package main

import (
	"fmt"
	"os"
	"strings"
)

func (pOwn *sMaker) templateJson() bool {
	strJsonPath := strings.Replace(pOwn.XlsxPath, gCommand.XlsxPath, gCommand.JsonPath, 1)
	strJsonName := strJsonPath + pOwn.RofName + ".json"
	os.MkdirAll(strJsonPath, os.ModeDir)
	pSheet := pOwn.File.Sheets[0]

	pOutFile, err := os.Create(strJsonName)
	if err != nil {
		logErr("can not create json file:%s", strJsonName)
		return false
	}
	pOutFile.WriteString("{\n")

	//填写内容
	for i := 3; i < pSheet.MaxRow; i++ {
		pRow := pSheet.Rows[i]
		strRowContent := ""
		for j := 0; j < len(pOwn.ColHeadList); j++ {
			pColHead := pOwn.ColHeadList[j]
			cell := pRow.Cells[pColHead.Index]

			if pColHead.Index == 0 {
				strRowContent = fmt.Sprintf("\"%s\":{", cell.String())
			}

			switch pColHead.Type {
			case "int32", "int64", "float32", "float64", "object":
				{
					strRowContent += fmt.Sprintf("\"%s\":%s", pColHead.Name, cell.String())
				}
				break
			case "string":
				{
					strRowContent += fmt.Sprintf("\"%s\":\"%s\"", pColHead.Name, cell.String())
				}
				break
			}

			if j == len(pOwn.ColHeadList)-1 {
				strRowContent += "}"
			} else {
				strRowContent += ","
			}
		}
		if i < pSheet.MaxRow-1 {
			strRowContent += ","
		}
		pOutFile.WriteString(strRowContent + "\n")
	}

	pOutFile.WriteString("}\n")
	pOutFile.Close()
	return true
}
