package main

import (
	"encoding/binary"
	"math"
	"os"
	"strconv"
	"strings"
)

func (pOwn *sMaker) templateRof() bool {
	strRofPath := strings.Replace(pOwn.XlsxPath, gCommand.XlsxPath, gCommand.RofPath, 1)
	strRofName := strRofPath + pOwn.RofName + ".bytes"
	os.MkdirAll(strRofPath, os.ModeDir)
	pSheet := pOwn.File.Sheets[0]
	//填写头
	nRealRowNum := pSheet.MaxRow - 3
	nRealColNum := len(pOwn.ColHeadList)
	pBuffer := make([]byte, cRofDefaultSize)
	nOffset := 64

	binary.BigEndian.PutUint32(pBuffer[nOffset:], uint32(nRealRowNum))
	nOffset += 4
	binary.BigEndian.PutUint32(pBuffer[nOffset:], uint32(nRealColNum))
	nOffset += 4

	//填写列属性
	for i := 0; i < len(pOwn.ColHeadList); i++ {
		pInfo := pOwn.ColHeadList[i]
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
		for j := 0; j < len(pOwn.ColHeadList); j++ {
			pColHead := pOwn.ColHeadList[j]
			cell := pRow.Cells[pColHead.Index]

			switch pColHead.Type {
			case "int32":
				{
					nValue, _ := strconv.ParseInt(cell.String(), 10, 32)
					binary.BigEndian.PutUint32(byteTempBuffer, uint32(nValue))
					copybuffer(&pBuffer, nOffset, byteTempBuffer, 4)
					nOffset += 4
				}
				break
			case "int64":
				{
					nValue, _ := strconv.ParseInt(cell.String(), 10, 64)
					binary.BigEndian.PutUint64(byteTempBuffer, uint64(nValue))
					copybuffer(&pBuffer, nOffset, byteTempBuffer, 8)
					nOffset += 8
				}
				break
			case "float32":
				{
					fValue, _ := strconv.ParseFloat(cell.String(), 32)
					bits := math.Float32bits(float32(fValue))
					binary.BigEndian.PutUint32(byteTempBuffer, bits)
					copybuffer(&pBuffer, nOffset, byteTempBuffer, 4)
					nOffset += 4
				}
				break
			case "float64":
				{
					fValue, _ := strconv.ParseFloat(cell.String(), 64)
					bits := math.Float64bits(fValue)
					binary.BigEndian.PutUint64(byteTempBuffer, bits)
					copybuffer(&pBuffer, nOffset, byteTempBuffer, 8)
					nOffset += 8
				}
				break
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
				break
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
