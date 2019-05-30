package main

import (
	"fmt"
	"os"
)

func (pOwn *sMaker) templateGo() bool {
	strGoName := gCommand.GoFilePath + pOwn.RofName + ".go"
	os.MkdirAll(gCommand.GoFilePath, os.ModeDir)
	pFile, err := os.Create(strGoName)
	if err != nil {
		logErr("can not create go file:%s", strGoName)
		return false
	}
	defer pFile.Close()

	bIncludeMath := false
	strRowClassName := fmt.Sprintf("s%sRow", pOwn.RofName)
	strTableClassName := fmt.Sprintf("s%sTable", pOwn.RofName)
	//row 结构体
	strContent := fmt.Sprintf("type %s struct {\n", strRowClassName)
	for i := 0; i < len(pOwn.ColHeadList); i++ {
		cell := pOwn.ColHeadList[i]
		strType := cGoTypeMap[cell.Type]
		strContent += fmt.Sprintf("m%s %s\n", cell.Name, strType)
	}
	strContent += "}\n"

	//ReadBody
	strContent += fmt.Sprintf("func (pOwn *%s) readBody(aBuffer []byte) int32 {\nvar nOffset int32\n", strRowClassName)
	for i := 0; i < len(pOwn.ColHeadList); i++ {
		cell := pOwn.ColHeadList[i]
		switch cell.Type {
		case "int32":
			{
				strContent += fmt.Sprintf("pOwn.m%s = int32(binary.BigEndian.Uint32(aBuffer[nOffset:]))\nnOffset+=4\n", cell.Name)
			}
			break
		case "int64":
			{
				strContent += fmt.Sprintf("pOwn.m%s = int64(binary.BigEndian.Uint64(aBuffer[nOffset:]))\nnOffset+=8\n", cell.Name)
			}
			break
		case "float32":
			{
				strContent += fmt.Sprintf("pOwn.m%s = math.Float32frombits(binary.BigEndian.Uint32(aBuffer[nOffset:]))\nnOffset+=4\n", cell.Name)
				bIncludeMath = true
			}
			break
		case "float64":
			{
				strContent += fmt.Sprintf("pOwn.m%s = math.Float64frombits(binary.BigEndian.Uint64(aBuffer[nOffset:]))\nnOffset+=8\n", cell.Name)
				bIncludeMath = true
			}
			break
		case "string", "object":
			{
				strContent += fmt.Sprintf("n%sLen := int32(binary.BigEndian.Uint32(aBuffer[nOffset:]))\nnOffset+=4\n", cell.Name)
				strContent += fmt.Sprintf("pOwn.m%s = string(aBuffer[nOffset:nOffset+n%sLen])\nnOffset+=n%sLen\n", cell.Name, cell.Name, cell.Name)
			}
			break
		}
	}
	strContent += "return nOffset\n}\n"

	//函数
	for i := 0; i < len(pOwn.ColHeadList); i++ {
		cell := pOwn.ColHeadList[i]
		strType := cGoTypeMap[cell.Type]
		strContent += fmt.Sprintf("func (pOwn *%s) Get%s() %s { return pOwn.m%s } \n", strRowClassName, cell.Name, strType, cell.Name)
	}

	//table 结构体
	strContent += fmt.Sprintf("type %s struct { \nmRowNum int32\nmColNum int32\nmIDMap  map[int32]*%s\nmRowMap map[int32]int32\n}\n", strTableClassName, strRowClassName)
	strContent += fmt.Sprintf("func (pOwn *%s) GetDataByID(aID int32) *%s {return pOwn.mIDMap[aID]}\n", strTableClassName, strRowClassName)
	strContent += fmt.Sprintf("func (pOwn *%s) GetDataByRow(aIndex int32) *%s {\n", strTableClassName, strRowClassName)
	strContent += fmt.Sprintf("nID, ok := pOwn.mRowMap[aIndex]\nif ok == false {return nil}\nreturn pOwn.mIDMap[nID]\n}\n")
	strContent += fmt.Sprintf("func (pOwn *%s) GetRows() int32 {return pOwn.mRowNum}\n", strTableClassName)
	strContent += fmt.Sprintf("func (pOwn *%s) GetCols() int32 {return pOwn.mColNum}\n", strTableClassName)

	strContent += fmt.Sprintf("func (pOwn *%s) init(aPath string) bool {\n", strTableClassName)
	strContent += fmt.Sprintf("pOwn.mIDMap = make(map[int32]*%s)\n", strRowClassName)
	strContent += fmt.Sprintf("pOwn.mRowMap = make(map[int32]int32)\n")
	strContent += fmt.Sprintf("pFile, err := os.Open(aPath)\nif err != nil {\nreturn false\n}\ndefer pFile.Close()\n")
	strContent += fmt.Sprintf("pFileInfo, _ := pFile.Stat()\nnFileSize := pFileInfo.Size()\npBuffer := make([]byte, nFileSize)\n")
	strContent += fmt.Sprintf("_, err = pFile.Read(pBuffer)\nif err != nil {\nreturn false\n}\n")
	strContent += fmt.Sprintf("var nOffset int32 = 64\n")
	strContent += fmt.Sprintf("pOwn.mRowNum = int32(binary.BigEndian.Uint32(pBuffer[nOffset:]))\nnOffset += 4\n")
	strContent += fmt.Sprintf("pOwn.mColNum = int32(binary.BigEndian.Uint32(pBuffer[nOffset:]))\nnOffset += 4\n")

	strContent += fmt.Sprintf("for i := 0; i < int(pOwn.mColNum); i++ {\n")
	strContent += fmt.Sprintf("nNameLen := int8(pBuffer[nOffset])\nnOffset += 1 + int32(nNameLen)\n")
	strContent += fmt.Sprintf("nTypeLen := int8(pBuffer[nOffset])\nnOffset += 1 + int32(nTypeLen)\n")
	strContent += fmt.Sprintf("}\n")
	strContent += fmt.Sprintf("for i := int32(0); i < pOwn.mRowNum; i++ {\n")
	strContent += fmt.Sprintf("nID := int32(binary.BigEndian.Uint32(pBuffer[nOffset:]))\n")
	strContent += fmt.Sprintf("pData := new(%s)\n", strRowClassName)
	strContent += fmt.Sprintf("nOffset += pData.readBody(pBuffer[nOffset:])\n")
	strContent += fmt.Sprintf("pOwn.mIDMap[nID] = pData\n")
	strContent += fmt.Sprintf("pOwn.mRowMap[i] = nID\n")
	strContent += fmt.Sprintf("}\nreturn true\n}\n")

	pFile.WriteString("package rof\n")
	pFile.WriteString("import \"encoding/binary\"\n")
	pFile.WriteString("import \"os\"\n")
	if bIncludeMath == true {
		pFile.WriteString("import \"math\"\n")
	}
	pFile.WriteString(strContent)
	return true
}
