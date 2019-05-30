package main

import (
	xlsx "libxlsx"
	"os"
)

func analysisFileList() bool {
	//判断资源路径是否存在
	_, statErr := os.Stat(gCommand.XlsxPath)
	if os.IsNotExist(statErr) == true {
		logErr("can not find xlsx path.")
		return false
	}
	return pushFileInfo(gCommand.XlsxPath)
}

func pushFileInfo(aDir string) bool {
	fileDir, _ := os.Open(aDir)
	defer fileDir.Close()
	fileList, _ := fileDir.Readdir(0)
	for _, v := range fileList {
		fileName := v.Name()
		if v.IsDir() == true {
			pushFileInfo(aDir + fileName + "/")
		}
		if len(fileName) < 6 || fileName[len(fileName)-5:] != ".xlsx" {
			continue
		}
		if fileName[:2] == "~$" {
			continue
		}

		strXlsxFileName := aDir + fileName
		pFile, err := xlsx.OpenFile(strXlsxFileName)
		if err != nil {
			logErr("open xlsx file error:%s, file: %s", err.Error(), strXlsxFileName)
			return false
		}
		pSheet := pFile.Sheets[0]

		pMaker := new(sMaker)
		pMaker.XlsxPath = aDir
		pMaker.XlsxName = fileName
		pMaker.RofName = "Rof" + pSheet.Name
		pMaker.File = pFile

		//判断表名是否重复
		pOldMaker, exist := gMakerMap[pMaker.RofName]
		if exist == true {
			logErr("repetitive table name: %s, %s is same as %s", pMaker.RofName, strXlsxFileName, (pOldMaker.XlsxPath + pOldMaker.XlsxName))
			return false
		}
		gMakerMap[pMaker.RofName] = pMaker
	}
	return true
}
