package main

import (
	"os"
)

type sGoDefineBuilder struct {
	mPath string //文件夹路径
}

func (pOwn *sGoDefineBuilder) getCommandDesc() string {
	return "-godef [path]. optional command, [path] is the output (.go) files floder."
}

func (pOwn *sGoDefineBuilder) init(aCmdParm []string) bool {
	if len(aCmdParm) != 1 {
		logErr("the command -godef needs 1 (only 1) argument.")
		return false
	}
	pOwn.mPath = aCmdParm[0]
	if pOwn.mPath[len(pOwn.mPath)-1] != '/' {
		pOwn.mPath += "/"
	}

	return true
}

func (pOwn *sGoDefineBuilder) build() bool {
	os.MkdirAll(pOwn.mPath, os.ModeDir)
	strGoName := pOwn.mPath + "RofDefine.go"
	pFile, err := os.Create(strGoName)
	if err != nil {
		logErr("can not create RofDefine.go")
		return false
	}
	defer pFile.Close()
	pFile.WriteString("package rof\n")
	pFile.WriteString("type nnkv struct {\nk int32\nv float64\n}")
	return true
}
