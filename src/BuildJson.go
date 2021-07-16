package main

import (
	"fmt"
	"os"
)

type sJsonBuilder struct {
	mPath string //文件夹路径
}

func (pOwn *sJsonBuilder) getCommandDesc() string {
	return "-json [path]. optional command, [path] is the output (.json) files floder."
}

func (pOwn *sJsonBuilder) init(aCmdParm []string) bool {
	if len(aCmdParm) != 1 {
		logErr("the command -go needs 1 (only 1) argument.")
		return false
	}
	pOwn.mPath = aCmdParm[0]
	if pOwn.mPath[len(pOwn.mPath)-1] != '/' {
		pOwn.mPath += "/"
	}
	return true
}

func (pOwn *sJsonBuilder) build() bool {
	for _, v := range gTables {
		if pOwn.doBuild(v) == false {
			return false
		}
	}
	return true
}

func (pOwn *sJsonBuilder) doBuild(aInfo *sTableInfo) bool {
	strJsonPath := pOwn.mPath + aInfo.RelativeDir
	strJsonName := strJsonPath + aInfo.RofName + ".json"
	os.MkdirAll(strJsonPath, os.ModeDir)
	pSheet := aInfo.File.Sheets[0]

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
		for j := 0; j < len(aInfo.ColHeadList); j++ {
			pColHead := aInfo.ColHeadList[j]
			cell := pRow.Cells[pColHead.Index]

			if pColHead.Index == 0 {
				strRowContent = fmt.Sprintf("\"%s\":{", cell.String())
			}

			switch pColHead.Type {
			case "int32", "int64", "float32", "float64", "object":
				{
					strRowContent += fmt.Sprintf("\"%s\":%s", pColHead.Name, cell.String())
				}
			case "string":
				{
					strRowContent += fmt.Sprintf("\"%s\":\"%s\"", pColHead.Name, cell.String())
				}
			}

			if j == len(aInfo.ColHeadList)-1 {
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
