package main

import (
	"os"
)

func analysisFileList(aDir string) bool {
	//判断资源路径是否存在
	_, statErr := os.Stat(aDir)
	if os.IsNotExist(statErr) == true {
		logErr("can not find xlsx path.")
		return false
	}
	return pushFileInfo(aDir)
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

		pInfo := new(sTableInfo)
		pInfo.Dir = aDir
		pInfo.FileName = fileName
		gTables = append(gTables, pInfo)
	}
	return true
}
