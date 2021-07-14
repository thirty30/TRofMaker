package main

import (
	"encoding/binary"
	"math"
	"os"
	"strconv"
)

type sRofBuilder struct {
	mPath string //文件夹路径
}

func (pOwn *sRofBuilder) getCommandDesc() string {
	return "-rof [path]. optional command, [path] is the output (.bytes) files floder."
}

func (pOwn *sRofBuilder) init(aCmdParm []string) bool {
	if len(aCmdParm) != 1 {
		logErr("the command -rof needs 1 (only 1) argument.")
		return false
	}
	pOwn.mPath = aCmdParm[0]
	if pOwn.mPath[len(pOwn.mPath)-1] != '/' {
		pOwn.mPath += "/"
	}
	return true
}

func (pOwn *sRofBuilder) build() bool {
	for _, v := range gTables {
		if pOwn.doBuild(v) == false {
			return false
		}
	}
	return true
}

func (pOwn *sRofBuilder) doBuild(aInfo *sTableInfo) bool {
	strRofPath := pOwn.mPath + aInfo.RelativeDir
	strRofName := strRofPath + aInfo.RofName + ".bytes"
	os.MkdirAll(strRofPath, os.ModeDir)
	pSheet := aInfo.File.Sheets[0]
	//填写头
	nRealRowNum := pSheet.MaxRow - 3
	nRealColNum := len(aInfo.ColHeadList)
	pBuffer := make([]byte, cRofDefaultSize)
	nOffset := 64

	binary.BigEndian.PutUint32(pBuffer[nOffset:], uint32(nRealRowNum))
	nOffset += 4
	binary.BigEndian.PutUint32(pBuffer[nOffset:], uint32(nRealColNum))
	nOffset += 4

	//填写列属性
	for i := 0; i < len(aInfo.ColHeadList); i++ {
		pInfo := aInfo.ColHeadList[i]
		nNameLen := int8(len(pInfo.Name))
		pBuffer[nOffset] = byte(nNameLen)
		nOffset++

		copy(pBuffer[nOffset:], []byte(pInfo.Name))
		nOffset += len(pInfo.Name)

		nTypeLen := int8(len(pInfo.Type))
		pBuffer[nOffset] = byte(nTypeLen)
		nOffset++

		copy(pBuffer[nOffset:], []byte(pInfo.Type))
		nOffset += len(pInfo.Type)
	}

	byteTempBuffer := make([]byte, 8)
	//填写内容
	for i := 3; i < pSheet.MaxRow; i++ {
		pRow := pSheet.Rows[i]
		for j := 0; j < len(aInfo.ColHeadList); j++ {
			pColHead := aInfo.ColHeadList[j]
			cell := pRow.Cells[pColHead.Index]

			switch pColHead.Type {
			case "int32":
				{
					nValue, _ := strconv.ParseInt(cell.String(), 10, 32)
					binary.BigEndian.PutUint32(byteTempBuffer, uint32(nValue))
					copybuffer(&pBuffer, nOffset, byteTempBuffer, 4)
					nOffset += 4
				}
			case "int64":
				{
					nValue, _ := strconv.ParseInt(cell.String(), 10, 64)
					binary.BigEndian.PutUint64(byteTempBuffer, uint64(nValue))
					copybuffer(&pBuffer, nOffset, byteTempBuffer, 8)
					nOffset += 8
				}
			case "float32":
				{
					fValue, _ := strconv.ParseFloat(cell.String(), 32)
					bits := math.Float32bits(float32(fValue))
					binary.BigEndian.PutUint32(byteTempBuffer, bits)
					copybuffer(&pBuffer, nOffset, byteTempBuffer, 4)
					nOffset += 4
				}
			case "float64":
				{
					fValue, _ := strconv.ParseFloat(cell.String(), 64)
					bits := math.Float64bits(fValue)
					binary.BigEndian.PutUint64(byteTempBuffer, bits)
					copybuffer(&pBuffer, nOffset, byteTempBuffer, 8)
					nOffset += 8
				}
			case "string", "object":
				{
					strValue := cell.String()
					nLen := len(strValue)
					binary.BigEndian.PutUint32(byteTempBuffer, uint32(nLen))
					copybuffer(&pBuffer, nOffset, byteTempBuffer, 4)
					nOffset += 4

					copybuffer(&pBuffer, nOffset, []byte(strValue), nLen)
					nOffset += nLen
				}
			}
		}
	}

	pRofFile, err := os.Create(strRofName)
	if err != nil {
		logErr("can not create rof file: %s", strRofName)
		return false
	}
	_, err = pRofFile.Write(pBuffer[0:nOffset])
	if err != nil {
		logErr("write rof file error:%s", err.Error())
		return false
	}
	return true
}
