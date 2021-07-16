package main

import (
	"fmt"
	"os"
)

type sCsBuilder struct {
	mPath    string //文件夹路径
	mTypeMap map[string]string
}

func (pOwn *sCsBuilder) getCommandDesc() string {
	return "-cs [path]. optional command, [path] is the output (.cs) files floder."
}

func (pOwn *sCsBuilder) init(aCmdParm []string) bool {
	if len(aCmdParm) != 1 {
		logErr("the command -cs needs 1 (only 1) argument.")
		return false
	}
	pOwn.mPath = aCmdParm[0]
	if pOwn.mPath[len(pOwn.mPath)-1] != '/' {
		pOwn.mPath += "/"
	}

	pOwn.mTypeMap = make(map[string]string)
	pOwn.mTypeMap["int32"] = "int"
	pOwn.mTypeMap["int64"] = "long"
	pOwn.mTypeMap["float32"] = "float"
	pOwn.mTypeMap["float64"] = "double"
	pOwn.mTypeMap["string"] = "string"
	pOwn.mTypeMap["object"] = "string"
	pOwn.mTypeMap["[]int32"] = "List<int>"
	pOwn.mTypeMap["[]int64"] = "List<long>"
	pOwn.mTypeMap["[]float32"] = "List<float>"
	pOwn.mTypeMap["[]float64"] = "List<double>"
	pOwn.mTypeMap["[]string"] = "List<string>"
	pOwn.mTypeMap["[]nnkv"] = "List<NNKV>"

	return true
}

func (pOwn *sCsBuilder) build() bool {
	//生成定义文件
	if pOwn.buildDefineFile() == false {
		return false
	}

	for _, v := range gTables {
		if pOwn.doBuild(v) == false {
			return false
		}
	}
	return true
}

func (pOwn *sCsBuilder) buildDefineFile() bool {
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

func (pOwn *sCsBuilder) doBuild(aInfo *sTableInfo) bool {
	strCsPath := pOwn.mPath + aInfo.RelativeDir
	strCsName := strCsPath + aInfo.RofName + ".cs"
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

	strRowClassName := fmt.Sprintf("%sRow", aInfo.RofName)
	strTableClassName := fmt.Sprintf("%sTable", aInfo.RofName)
	//row类
	strContent := fmt.Sprintf("public class %s\n{\n", strRowClassName)
	for i := 0; i < len(aInfo.ColHeadList); i++ {
		cell := aInfo.ColHeadList[i]
		strType := pOwn.mTypeMap[cell.Type]
		strContent += fmt.Sprintf("public %s %s { get; private set; }\n", strType, cell.Name)
	}

	//ReadBody
	strContent += "public int ReadBody(byte[] rData, int nOffset)\n{\n"
	for i := 0; i < len(aInfo.ColHeadList); i++ {
		cell := aInfo.ColHeadList[i]
		switch cell.Type {
		case "int32":
			{
				strContent += "if (BitConverter.IsLittleEndian){Array.Reverse(rData, nOffset, 4);}\n"
				strContent += fmt.Sprintf("this.%s = (int)BitConverter.ToUInt32(rData, nOffset); nOffset += 4;\n", cell.Name)
			}
		case "int64":
			{
				strContent += "if (BitConverter.IsLittleEndian){Array.Reverse(rData, nOffset, 8);}\n"
				strContent += fmt.Sprintf("this.%s = (long)BitConverter.ToUInt64(rData, nOffset); nOffset += 8;\n", cell.Name)
			}
		case "float32":
			{
				strContent += "if (BitConverter.IsLittleEndian){Array.Reverse(rData, nOffset, 4);}\n"
				strContent += fmt.Sprintf("this.%s = BitConverter.ToSingle(rData, nOffset); nOffset += 4;\n", cell.Name)
			}
		case "float64":
			{
				strContent += "if (BitConverter.IsLittleEndian){Array.Reverse(rData, nOffset, 8);}\n"
				strContent += fmt.Sprintf("this.%s = BitConverter.ToDouble(rData, nOffset); nOffset += 8;\n", cell.Name)
			}
		case "string", "object":
			{
				strContent += "if (BitConverter.IsLittleEndian){Array.Reverse(rData, nOffset, 4);}\n"
				strContent += fmt.Sprintf("int n%sLen = (int)BitConverter.ToUInt32(rData, nOffset); nOffset += 4;\n", cell.Name)
				strContent += fmt.Sprintf("this.%s = Encoding.UTF8.GetString(rData, nOffset, n%sLen); nOffset += n%sLen;\n", cell.Name, cell.Name, cell.Name)
			}
		case "[]int32":
			{
				strContent += fmt.Sprintf("this.%s = new List<int>();\n", cell.Name)
				strContent += "if (BitConverter.IsLittleEndian){Array.Reverse(rData, nOffset, 4);}\n"
				strContent += fmt.Sprintf("int n%sLen = (int)BitConverter.ToUInt32(rData, nOffset); nOffset += 4;\n", cell.Name)
				strContent += fmt.Sprintf("for (int i = 0; i < n%sLen; i++)\n", cell.Name)
				strContent += "{\n"
				strContent += "if (BitConverter.IsLittleEndian){Array.Reverse(rData, nOffset, 4);}\n"
				strContent += fmt.Sprintf("this.%s.Add((int)BitConverter.ToUInt32(rData, nOffset)); nOffset += 4;\n", cell.Name)
				strContent += "}\n"
			}
		case "[]int64":
			{
				strContent += fmt.Sprintf("this.%s = new List<long>();\n", cell.Name)
				strContent += "if (BitConverter.IsLittleEndian){Array.Reverse(rData, nOffset, 4);}\n"
				strContent += fmt.Sprintf("int n%sLen = (int)BitConverter.ToUInt32(rData, nOffset); nOffset += 4;\n", cell.Name)
				strContent += fmt.Sprintf("for (int i = 0; i < n%sLen; i++)\n", cell.Name)
				strContent += "{\n"
				strContent += "if (BitConverter.IsLittleEndian){Array.Reverse(rData, nOffset, 8);}\n"
				strContent += fmt.Sprintf("this.%s.Add((long)BitConverter.ToUInt64(rData, nOffset)); nOffset += 8;\n", cell.Name)
				strContent += "}\n"
			}
		case "[]float32":
			{
				strContent += fmt.Sprintf("this.%s = new List<float>();\n", cell.Name)
				strContent += "if (BitConverter.IsLittleEndian){Array.Reverse(rData, nOffset, 4);}\n"
				strContent += fmt.Sprintf("int n%sLen = (int)BitConverter.ToUInt32(rData, nOffset); nOffset += 4;\n", cell.Name)
				strContent += fmt.Sprintf("for (int i = 0; i < n%sLen; i++)\n", cell.Name)
				strContent += "{\n"
				strContent += "if (BitConverter.IsLittleEndian){Array.Reverse(rData, nOffset, 4);}\n"
				strContent += fmt.Sprintf("this.%s.Add(BitConverter.ToSingle(rData, nOffset)); nOffset += 4;\n", cell.Name)
				strContent += "}\n"
			}
		case "[]float64":
			{
				strContent += fmt.Sprintf("this.%s = new List<double>();\n", cell.Name)
				strContent += "if (BitConverter.IsLittleEndian){Array.Reverse(rData, nOffset, 4);}\n"
				strContent += fmt.Sprintf("int n%sLen = (int)BitConverter.ToUInt32(rData, nOffset); nOffset += 4;\n", cell.Name)
				strContent += fmt.Sprintf("for (int i = 0; i < n%sLen; i++)\n", cell.Name)
				strContent += "{\n"
				strContent += "if (BitConverter.IsLittleEndian){Array.Reverse(rData, nOffset, 8);}\n"
				strContent += fmt.Sprintf("this.%s.Add(BitConverter.ToDouble(rData, nOffset)); nOffset += 8;\n", cell.Name)
				strContent += "}\n"
			}
		case "[]string":
			{
				strContent += fmt.Sprintf("this.%s = new List<string>();\n", cell.Name)
				strContent += "if (BitConverter.IsLittleEndian){Array.Reverse(rData, nOffset, 4);}\n"
				strContent += fmt.Sprintf("int n%sLen = (int)BitConverter.ToUInt32(rData, nOffset); nOffset += 4;\n", cell.Name)
				strContent += fmt.Sprintf("string %sTempBuf = Encoding.UTF8.GetString(rData, nOffset, n%sLen); nOffset += n%sLen;\n", cell.Name, cell.Name, cell.Name)
				strContent += fmt.Sprintf("string[] %sElements = %sTempBuf.Split(',');\n", cell.Name, cell.Name)
				strContent += fmt.Sprintf("foreach(string item in %sElements)\n", cell.Name)
				strContent += "{\n"
				strContent += fmt.Sprintf("this.%s.Add(item);\n", cell.Name)
				strContent += "}\n"
			}
		case "[]nnkv":
			{
				strContent += fmt.Sprintf("this.%s = new List<NNKV>();\n", cell.Name)
				strContent += "if (BitConverter.IsLittleEndian){Array.Reverse(rData, nOffset, 4);}\n"
				strContent += fmt.Sprintf("int n%sLen = (int)BitConverter.ToUInt32(rData, nOffset); nOffset += 4;\n", cell.Name)
				strContent += fmt.Sprintf("for (int i = 0; i < n%sLen; i++)\n", cell.Name)
				strContent += "{\n"
				strContent += "if (BitConverter.IsLittleEndian){Array.Reverse(rData, nOffset, 4);}\n"
				strContent += "int nKey = (int)BitConverter.ToUInt32(rData, nOffset); nOffset += 4;\n"
				strContent += "if (BitConverter.IsLittleEndian){Array.Reverse(rData, nOffset, 8);}\n"
				strContent += "double nValue = BitConverter.ToDouble(rData, nOffset); nOffset += 8;\n"
				strContent += fmt.Sprintf("this.%s.Add(new NNKV(nKey, nValue));\n", cell.Name)
				strContent += "}\n"
			}
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
