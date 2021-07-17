package main

import (
	"os"
)

type sCsDefineBuilder struct {
	mPath string //文件夹路径
}

func (pOwn *sCsDefineBuilder) getCommandDesc() string {
	return "-csdef [path]. optional command, [path] is the output (.cs) files floder."
}

func (pOwn *sCsDefineBuilder) init(aCmdParm []string) bool {
	if len(aCmdParm) != 1 {
		logErr("the command -csdef needs 1 (only 1) argument.")
		return false
	}
	pOwn.mPath = aCmdParm[0]
	if pOwn.mPath[len(pOwn.mPath)-1] != '/' {
		pOwn.mPath += "/"
	}
	return true
}

func (pOwn *sCsDefineBuilder) build() bool {
	os.MkdirAll(pOwn.mPath, os.ModeDir)
	strGoName := pOwn.mPath + "RofDefine.cs"
	pFile, err := os.Create(strGoName)
	if err != nil {
		logErr("can not create RofDefine.cs")
		return false
	}
	defer pFile.Close()

	strContent := "namespace Rof\n"
	strContent += "{\n"
	strContent += "public class NNKV\n"
	strContent += "{\n"
	strContent += "public int Key { get; private set; }\n"
	strContent += "public double Value { get; private set; }\n"
	strContent += "public NNKV(int aKey, double aValue)\n"
	strContent += "{\n"
	strContent += "this.Key = aKey;\n"
	strContent += "this.Value = aValue;\n"
	strContent += "}\n"
	strContent += "}\n"
	strContent += "}\n"

	pFile.WriteString(strContent)
	return true
}
