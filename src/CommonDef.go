package main

import xlsx "libxlsx"

const (
	cUnexportableColor  = "FF808080"      //不导出标识颜色
	cMultiLanguageColor = "FF00B0F0"      //多语言标识颜色
	cRofDefaultSize     = 1024 * 1024 * 1 //Rof文件初始大小1M
)

//列头信息
type sColHeadInfo struct {
	Index int    //列索引
	Name  string //列名
	Type  string //数据类型
	IsLan bool   //是否是多语言字段
}

//表信息
type sTableInfo struct {
	Dir         string          //xlsx路径
	RelativeDir string          //xlsx相对路径
	FileName    string          //xlsx文件名
	RofName     string          //rof文件名
	TableName   string          //表名
	File        *xlsx.File      //excel表对象
	ColHeadList []*sColHeadInfo //列头信息
}

type sCommandItem struct {
	Cmd        string
	Parms      []string
	Builder    IBuilder
	CanExecute bool
}

type IBuilder interface {
	getCommandDesc() string
	init(aCmdParm []string) bool
	build() bool
}
