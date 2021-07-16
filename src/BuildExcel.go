package main

import (
	xlsx "libxlsx"
	"strconv"
	"strings"
)

type sExcelBuilder struct {
	mPath     string          //目标路径
	mIsFolder bool            //是否是文件夹路径
	mTypeMap  map[string]bool //预制类型
}

func (pOwn *sExcelBuilder) getCommandDesc() string {
	return "-xlsx [path] or [file]. necessary command, analysis excel files. need a floder or a file path"
}

func (pOwn *sExcelBuilder) init(aCmdParm []string) bool {
	if len(aCmdParm) != 1 {
		logErr("the command -xlsx needs 1 (only 1) argument, a floder or a file path.")
		return false
	}
	pOwn.mPath = aCmdParm[0]

	//单个文件
	if len(pOwn.mPath) >= 6 && pOwn.mPath[len(pOwn.mPath)-5:] == ".xlsx" {
		pOwn.mIsFolder = false
		tStr := strings.Split(pOwn.mPath, "/")
		fileName := tStr[len(tStr)-1]
		dir := strings.Replace(pOwn.mPath, fileName, "", 1)

		pInfo := new(sTableInfo)
		pInfo.Dir = dir
		pInfo.FileName = fileName
		pInfo.RelativeDir = ""
		gTables = append(gTables, pInfo)

	} else {
		pOwn.mIsFolder = true

		if pOwn.mPath[len(pOwn.mPath)-1] != '/' {
			pOwn.mPath += "/"
		}

		//解析待生成的文件的列表
		if analysisFileList(pOwn.mPath) == false {
			return false
		}
	}

	//填充表结构信息
	for _, v := range gTables {
		strXlsxFileName := v.Dir + v.FileName
		pFile, err := xlsx.OpenFile(strXlsxFileName)
		if err != nil {
			logErr("open xlsx file error:%s, file: %s", err.Error(), strXlsxFileName)
			return false
		}
		pSheet := pFile.Sheets[0]
		rofName := "Rof" + pSheet.Name

		//判断表名是否重复
		isRepeat, srcPath := isRofNameRepeated(rofName)
		if isRepeat == true {
			logErr("repetitive table name: %s, %s is same as %s", rofName, strXlsxFileName, srcPath)
			return false
		}

		v.RofName = rofName
		v.File = pFile
		v.TableName = pSheet.Name
		if pOwn.mIsFolder == true {
			v.RelativeDir = strings.Replace(v.Dir, pOwn.mPath, "", 1)
		}
	}

	pOwn.mTypeMap = make(map[string]bool)
	pOwn.mTypeMap["int32"] = true
	pOwn.mTypeMap["int64"] = true
	pOwn.mTypeMap["float32"] = true
	pOwn.mTypeMap["float64"] = true
	pOwn.mTypeMap["string"] = true
	pOwn.mTypeMap["object"] = true
	pOwn.mTypeMap["[]int32"] = true
	pOwn.mTypeMap["[]int64"] = true
	pOwn.mTypeMap["[]float32"] = true
	pOwn.mTypeMap["[]float64"] = true
	pOwn.mTypeMap["[]string"] = true
	pOwn.mTypeMap["[]nnkv"] = true

	return true
}

func (pOwn *sExcelBuilder) build() bool {
	for _, v := range gTables {
		if pOwn.preprocessTableHead(v) == false {
			return false
		}

		if pOwn.preprocessContent(v) == false {
			return false
		}
	}
	return true
}

//预处理表头
func (pOwn *sExcelBuilder) preprocessTableHead(aInfo *sTableInfo) bool {
	strXlsxName := aInfo.Dir + aInfo.FileName
	pSheet := aInfo.File.Sheets[0]
	pFlagRow := pSheet.Rows[0]
	pNameRow := pSheet.Rows[1]
	pTypeRow := pSheet.Rows[2]

	//判断列数是否匹配
	if pSheet.MaxCol != len(pFlagRow.Cells) || pSheet.MaxCol != len(pNameRow.Cells) || pSheet.MaxCol != len(pTypeRow.Cells) {
		logErr("column num is not match max column num, file: %s", strXlsxName)
		return false
	}

	//得到真实的列
	aInfo.ColHeadList = make([]*sColHeadInfo, 0)
	for i := 0; i < pSheet.MaxCol; i++ {
		pFlagCell := pFlagRow.Cells[i]
		pNameCell := pNameRow.Cells[i]
		pTypeCell := pTypeRow.Cells[i]
		strFlagColor := pFlagCell.GetStyle().Fill.FgColor
		//判断第一列规范
		if i == 0 {
			//必须可导出
			if strFlagColor == cUnexportableColor {
				logErr("the first column can not be marked as unexportable, file: %s", strXlsxName)
				return false
			}
			//不能是多语言列
			if strFlagColor == cMultiLanguageColor {
				logErr("the first column can not be marked as multi-language, file: %s", strXlsxName)
				return false
			}
			//列名必须是ID
			if pNameCell.String() != "ID" {
				logErr("the first column name must be ID, file: %s", strXlsxName)
				return false
			}
			//列类型必须是int32
			if pTypeCell.String() != "int32" {
				logErr("the first column type must be int32, file: %s", strXlsxName)
				return false
			}
		}

		//判断是不导出列
		if strFlagColor == cUnexportableColor {
			continue
		}

		//判断列名存在
		if len(pNameCell.String()) <= 0 {
			logErr("the column name is empty, file: %s", strXlsxName)
			return false
		}
		if len(pNameCell.String()) > 255 {
			logErr("the length of column name is more than 255, file: %s", strXlsxName)
			return false
		}

		//判断列名重复
		for i := 0; i < len(aInfo.ColHeadList); i++ {
			if aInfo.ColHeadList[i].Name == pNameCell.String() {
				logErr("repetitive column name, column:%s, file:%s", pNameCell.String(), strXlsxName)
				return false
			}
		}

		//判断数据类型是否合法
		_, bExist := pOwn.mTypeMap[pTypeCell.String()]
		if bExist == false {
			logErr("illegal data type:%s, column:%s, file:%s", pTypeCell.String(), pNameCell.String(), strXlsxName)
		}

		//判断多语言列的类型必须是int32
		bIsLan := false
		if strFlagColor == cMultiLanguageColor {
			bIsLan = true
			if pTypeCell.String() != "int32" {
				logErr("the language column type is not int32, file: %s", strXlsxName)
				return false
			}
		}

		pColInfo := new(sColHeadInfo)
		pColInfo.Index = i
		pColInfo.Name = pNameCell.String()
		pColInfo.Type = pTypeCell.String()
		pColInfo.IsLan = bIsLan

		aInfo.ColHeadList = append(aInfo.ColHeadList, pColInfo)
	}
	return true
}

//预处理内容
func (pOwn *sExcelBuilder) preprocessContent(aInfo *sTableInfo) bool {
	mapRepetitionID := make(map[int32]bool)
	pSheet := aInfo.File.Sheets[0]
	strXlsxName := aInfo.Dir + aInfo.FileName

	//验证内容
	for i := 3; i < pSheet.MaxRow; i++ {
		pRow := pSheet.Rows[i]
		nShowRowNum := i + 1 //用于显示的行号
		for j := 0; j < len(aInfo.ColHeadList); j++ {
			pColHead := aInfo.ColHeadList[j]
			cell := pRow.Cells[pColHead.Index]
			//判断不能为空
			if len(cell.String()) <= 0 {
				logErr("cell's value is empty in row: %d, column: %s, file: %s", nShowRowNum, pColHead.Name, strXlsxName)
				return false
			}

			//检查ID不能重复
			if pColHead.Index == 0 {
				nValue, err := strconv.ParseInt(cell.String(), 10, 32)
				if err != nil {
					logErr("the ID is not int32 in row num:%d, file: %s,", nShowRowNum, strXlsxName)
					return false
				}
				_, ok := mapRepetitionID[int32(nValue)]
				if ok == true {
					logErr("repetitive ID in row: %d, file: %s", nShowRowNum, strXlsxName)
					return false
				}
				mapRepetitionID[int32(nValue)] = true
			}

			//检查类型正确
			switch pColHead.Type {
			case "int32":
				{
					_, err := strconv.ParseInt(cell.String(), 10, 32)
					if err != nil {
						logErr("the value's type is not int32 in row :%d, column: %s, file: %s", nShowRowNum, pColHead.Name, strXlsxName)
						return false
					}
				}
			case "int64":
				{
					_, err := strconv.ParseInt(cell.String(), 10, 64)
					if err != nil {
						logErr("the value's type is not int64 in row :%d, column: %s, file: %s", nShowRowNum, pColHead.Name, strXlsxName)
						return false
					}
				}
			case "float32":
				{
					_, err := strconv.ParseFloat(cell.String(), 32)
					if err != nil {
						logErr("the value's type is not float32 in row :%d, column: %s, file: %s", nShowRowNum, pColHead.Name, strXlsxName)
						return false
					}
				}
			case "float64":
				{
					_, err := strconv.ParseFloat(cell.String(), 64)
					if err != nil {
						logErr("the value's type is not float64 in row :%d, column: %s, file: %s", nShowRowNum, pColHead.Name, strXlsxName)
						return false
					}
				}
			case "[]int32":
				{
					content := cell.String()
					elements := strings.Split(content, ",")
					for _, item := range elements {
						_, err := strconv.ParseInt(item, 10, 32)
						if err != nil {
							logErr("can not analysis the value with []int32 in row :%d, column: %s, file: %s", nShowRowNum, pColHead.Name, strXlsxName)
							return false
						}
					}
				}
			case "[]int64":
				{
					content := cell.String()
					elements := strings.Split(content, ",")
					for _, item := range elements {
						_, err := strconv.ParseInt(item, 10, 64)
						if err != nil {
							logErr("can not analysis the value with []int64 in row :%d, column: %s, file: %s", nShowRowNum, pColHead.Name, strXlsxName)
							return false
						}
					}
				}
			case "[]float32":
				{
					content := cell.String()
					elements := strings.Split(content, ",")
					for _, item := range elements {
						_, err := strconv.ParseFloat(item, 32)
						if err != nil {
							logErr("can not analysis the value with []float32 in row :%d, column: %s, file: %s", nShowRowNum, pColHead.Name, strXlsxName)
							return false
						}
					}
				}
			case "[]float64":
				{
					content := cell.String()
					elements := strings.Split(content, ",")
					for _, item := range elements {
						_, err := strconv.ParseFloat(item, 64)
						if err != nil {
							logErr("can not analysis the value with []float64 in row :%d, column: %s, file: %s", nShowRowNum, pColHead.Name, strXlsxName)
							return false
						}
					}
				}
			case "[]nnkv":
				{
					content := cell.String()
					elements := strings.Split(content, ",")
					for _, item := range elements {
						pair := strings.Split(item, ":")
						if len(pair) != 2 {
							logErr("wrong kv pair format in row :%d, column: %s, file: %s", nShowRowNum, pColHead.Name, strXlsxName)
							return false
						}

						k := pair[0]
						v := pair[1]

						_, err := strconv.ParseInt(k, 10, 32)
						if err != nil {
							logErr("can not analysis the pair's key with int32 in row :%d, column: %s, file: %s", nShowRowNum, pColHead.Name, strXlsxName)
							return false
						}
						_, err = strconv.ParseFloat(v, 64)
						if err != nil {
							logErr("can not analysis the value with int32 nor float32 in row :%d, column: %s, file: %s", nShowRowNum, pColHead.Name, strXlsxName)
							return false
						}
					}
				}
			}
		}
	}
	return true
}
