package main

const (
	cUnexportableColor  = "FF808080"      //不导出标识颜色
	cMultiLanguageColor = "FF00B0F0"      //多语言标识颜色
	cRofDefaultSize     = 1024 * 1024 * 1 //Rof文件初始大小1M
)

type sCommand struct {
	XlsxPath   string //xlsx文件夹路径
	RofPath    string //rof输出路径
	JsonPath   string //json输出路径
	GoFilePath string
	CsFilePath string
}

//列头信息
type sColHeadInfo struct {
	Index int    //列索引
	Name  string //列名
	Type  string //数据类型
	IsLan bool   //是否是多语言字段
}

var gCommand sCommand            //指令
var gMakerMap map[string]*sMaker //文件Maker的map
var cDataType map[string]bool    //规定的数据类型
var cGoTypeMap map[string]string //go的数据类型对应关系
var cCsTypeMap map[string]string //cs的数据类型对应关系

func init() {
	gCommand.XlsxPath = ""
	gCommand.RofPath = ""
	gCommand.JsonPath = ""
	gCommand.GoFilePath = ""
	gCommand.CsFilePath = ""

	gMakerMap = make(map[string]*sMaker)
	cDataType = make(map[string]bool)
	cDataType["int32"] = true
	cDataType["int64"] = true
	cDataType["float32"] = true
	cDataType["float64"] = true
	cDataType["string"] = true
	cDataType["object"] = true

	cGoTypeMap = make(map[string]string)
	cGoTypeMap["int32"] = "int32"
	cGoTypeMap["int64"] = "int64"
	cGoTypeMap["float32"] = "float32"
	cGoTypeMap["float64"] = "float64"
	cGoTypeMap["string"] = "string"
	cGoTypeMap["object"] = "string"

	cCsTypeMap = make(map[string]string)
	cCsTypeMap["int32"] = "int"
	cCsTypeMap["int64"] = "long"
	cCsTypeMap["float32"] = "float"
	cCsTypeMap["float64"] = "double"
	cCsTypeMap["string"] = "string"
	cCsTypeMap["object"] = "string"
}

func copybuffer(aSrc *[]byte, aOffset int, aDest []byte, aLen int) {
	nTotalLen := cap(*aSrc)
	if nTotalLen-aOffset < aLen {
		aNewBuffer := make([]byte, nTotalLen*2)
		copy(aNewBuffer, *aSrc)
		*aSrc = aNewBuffer
	}
	copy((*aSrc)[aOffset:], aDest[0:aLen])
}
