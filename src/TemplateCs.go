package main

import (
	"fmt"
	"os"
	"strings"
)

func (pOwn *sMaker) templateCs() bool {
	strCsPath := strings.Replace(pOwn.XlsxPath, gCommand.XlsxPath, gCommand.CsFilePath, 1)
	strCsName := strCsPath + pOwn.RofName + ".cs"
	os.MkdirAll(strCsPath, os.ModeDir)
	pFile, err := os.Create(strCsName)
	if err != nil {
		logErr("can not create cs file:%s", strCsName)
		return false
	}
	defer pFile.Close()

	pFile.WriteString("using System;\n")
	pFile.WriteString("using System.Text;\n")
	pFile.WriteString("using System.Collections.Generic;\n")
	pFile.WriteString("namespace Rof\n{\n")

	strRowClassName := fmt.Sprintf("%sRow", pOwn.RofName)
	strTableClassName := fmt.Sprintf("%sTable", pOwn.RofName)
	//row类
	strContent := fmt.Sprintf("public class %s\n{\n", strRowClassName)
	for i := 0; i < len(pOwn.ColHeadList); i++ {
		cell := pOwn.ColHeadList[i]
		strType := cCsTypeMap[cell.Type]
		strContent += fmt.Sprintf("public %s %s { get; private set; }\n", strType, cell.Name)
	}

	//ReadBody
	strContent += "public int ReadBody(byte[] rData, int nOffset)\n{\n"
	for i := 0; i < len(pOwn.ColHeadList); i++ {
		cell := pOwn.ColHeadList[i]
		switch cell.Type {
		case "int32":
			{
				strContent += "if (BitConverter.IsLittleEndian){Array.Reverse(rData, nOffset, 4);}\n"
				strContent += fmt.Sprintf("%s = (int)BitConverter.ToUInt32(rData, nOffset); nOffset += 4;\n", cell.Name)
			}
			break
		case "int64":
			{
				strContent += "if (BitConverter.IsLittleEndian){Array.Reverse(rData, nOffset, 8);}\n"
				strContent += fmt.Sprintf("%s = (long)BitConverter.ToUInt64(rData, nOffset); nOffset += 8;\n", cell.Name)
			}
			break
		case "float32":
			{
				strContent += "if (BitConverter.IsLittleEndian){Array.Reverse(rData, nOffset, 4);}\n"
				strContent += fmt.Sprintf("%s = BitConverter.ToSingle(rData, nOffset); nOffset += 4;\n", cell.Name)
			}
			break
		case "float64":
			{
				strContent += "if (BitConverter.IsLittleEndian){Array.Reverse(rData, nOffset, 8);}\n"
				strContent += fmt.Sprintf("%s = BitConverter.ToDouble(rData, nOffset); nOffset += 8;\n", cell.Name)
			}
			break
		case "string", "object":
			{
				strContent += "if (BitConverter.IsLittleEndian){Array.Reverse(rData, nOffset, 4);}\n"
				strContent += fmt.Sprintf("int n%sLen = (int)BitConverter.ToUInt32(rData, nOffset); nOffset += 4;\n", cell.Name)
				strContent += fmt.Sprintf("%s = Encoding.UTF8.GetString(rData, nOffset, n%sLen); nOffset += n%sLen;\n", cell.Name, cell.Name, cell.Name)
			}
			break
		}
	}
	strContent += "return nOffset;\n}\n"
	strContent += "}\n"

	//table类
	strContent += fmt.Sprintf("public class %s\n{\n", strTableClassName)
	strContent += "private int mColNum;\nprivate int mRowNum;\n"
	strContent += fmt.Sprintf("private Dictionary<int, %s> mIDMap;\n", strRowClassName)
	strContent += "private Dictionary<int, int> mRowMap;\n"
	strContent += "public int RowNum { get { return this.mRowNum; } }\n"
	strContent += "public int ColNum { get { return this.mColNum; } }\n"
	strContent += "public void Init(byte[] rTotalBuffer)\n{\n"
	strContent += fmt.Sprintf("mIDMap = new Dictionary<int, %s>();\n", strRowClassName)
	strContent += "this.mRowMap = new Dictionary<int, int>();\n"
	strContent += "int nOffset = 64;\n"
	strContent += "if (BitConverter.IsLittleEndian) { Array.Reverse(rTotalBuffer, nOffset, 4); }\n"
	strContent += "this.mRowNum = (int)BitConverter.ToUInt32(rTotalBuffer, nOffset); nOffset += 4;\n"
	strContent += "if (BitConverter.IsLittleEndian) { Array.Reverse(rTotalBuffer, nOffset, 4); }\n"
	strContent += "this.mColNum = (int)BitConverter.ToUInt32(rTotalBuffer, nOffset); nOffset += 4;\n"
	strContent += "for (int i = 0; i < this.mColNum; i++)\n"
	strContent += "{\n"
	strContent += "int nNameLen = (int)rTotalBuffer[nOffset];\n"
	strContent += "nOffset += 1 + nNameLen;\n"
	strContent += "int nTypeLen = (int)rTotalBuffer[nOffset];\n"
	strContent += "nOffset += 1 + nTypeLen;\n"
	strContent += "}\n"
	strContent += "for (int i = 0; i < this.mRowNum; i++)\n"
	strContent += "{\n"
	strContent += "if (BitConverter.IsLittleEndian) { Array.Reverse(rTotalBuffer, nOffset, 4); }\n"
	strContent += "int nID = (int)BitConverter.ToUInt32(rTotalBuffer, nOffset);\n"
	strContent += "if (BitConverter.IsLittleEndian) { Array.Reverse(rTotalBuffer, nOffset, 4); }\n"
	strContent += fmt.Sprintf("%s rModel = new %s();\n", strRowClassName, strRowClassName)
	strContent += "nOffset = rModel.ReadBody(rTotalBuffer, nOffset);\n"
	strContent += "this.mIDMap.Add(nID, rModel);\n"
	strContent += "this.mRowMap.Add(i, nID);\n"
	strContent += "}\n"
	strContent += "}\n"
	strContent += fmt.Sprintf("public %s GetDataByID(int nID)\n", strRowClassName)
	strContent += "{\n"
	strContent += "if (this.mIDMap.ContainsKey(nID) == false)\n"
	strContent += "{\n"
	strContent += "return null;\n"
	strContent += "}\n"
	strContent += "return this.mIDMap[nID];\n"
	strContent += "}\n"
	strContent += fmt.Sprintf("public %s GetDataByRow(int nIndex)\n", strRowClassName)
	strContent += "	{\n"
	strContent += "if (mRowMap.ContainsKey(nIndex) == false)\n"
	strContent += "{\n"
	strContent += "return null;\n"
	strContent += "}\n"
	strContent += "int nID = mRowMap[nIndex];\n"
	strContent += "return mIDMap[nID];\n"
	strContent += "}\n"

	strContent += "}\n"
	strContent += "}\n"
	pFile.WriteString(strContent)
	return true
}
