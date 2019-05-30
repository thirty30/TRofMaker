package main

import (
	"os"
)

func main() {
	//初始化控制台输出颜色
	initConsoleColor()

	//解析命令
	if analysisArgs(os.Args[1:]) == false {
		return
	}

	//组织待生成的文件的列表
	if analysisFileList() == false {
		return
	}

	//生成文件
	if process() == false {
		return
	}

	log("[SUCCESS] Generate completely!")
}

func process() bool {
	for _, pMaker := range gMakerMap {
		if pMaker.process() == false {
			return false
		}
	}
	return true
}
