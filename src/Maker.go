package main

import (
	xlsx "libxlsx"
	"strconv"
)

type sMaker struct {
	XlsxPath    string          //xlsx路径
	XlsxName    string          //xlsx文件名
	RofName     string          //rof文件名
	ColHeadList []*sColHeadInfo //列头信息
	File        *xlsx.File
}

func (pOwn *sMaker) process() bool {
	if pOwn.preprocessTableHead() == false {
		return false
	}
	if pOwn.preprocessContent() == false {
		return false
	}
	if pOwn.doProcess() == false {
		return false
	}
	return true
}

//预处理表头
func (pOwn *sMaker) preprocessTableHead() bool {
	strXlsxName := pOwn.XlsxPath + pOwn.XlsxName
	pSheet := pOwn.File.Sheets[0]
	pFlagRow := pSheet.Rows[0]
	pNameRow := pSheet.Rows[1]
	pTypeRow := pSheet.Rows[2]

	//判断列数是否匹配
	if pSheet.MaxCol != len(pFlagRow.Cells) || pSheet.MaxCol != len(pNameRow.Cells) || pSheet.MaxCol != len(pTypeRow.Cells) {
		logErr("column num is not match max column num, file: %s", strXlsxName)
		return false
	}

	//得到真实的列
	pOwn.ColHeadList = make([]*sColHeadInfo, 0)
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
		for i := 0; i < len(pOwn.ColHeadList); i++ {
			if pOwn.ColHeadList[i].Name == pNameCell.String() {
				logErr("repetitive column name, column:%s, file:%s", pNameCell.String(), strXlsxName)
				return false
			}
		}

		//判断数据类型是否合法
		_, bExist := cDataType[pTypeCell.String()]
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

		pOwn.ColHeadList = append(pOwn.ColHeadList, pColInfo)
	}
	return true
}

//预处理内容
func (pOwn *sMaker) preprocessContent() bool {
	mapRepetitionID := make(map[int32]bool)
	pSheet := pOwn.File.Sheets[0]
	strXlsxName := pOwn.XlsxPath + pOwn.XlsxName

	//验证内容
	for i := 3; i < pSheet.MaxRow; i++ {
		pRow := pSheet.Rows[i]
		nShowRowNum := i + 1 //用于显示的行号
		for j := 0; j < len(pOwn.ColHeadList); j++ {
			pColHead := pOwn.ColHeadList[j]
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
				break
			case "int64":
				{
					_, err := strconv.ParseInt(cell.String(), 10, 64)
					if err != nil {
						logErr("the value's type is not int64 in row :%d, column: %s, file: %s", nShowRowNum, pColHead.Name, strXlsxName)
						return false
					}
				}
				break
			case "float32":
				{
					_, err := strconv.ParseFloat(cell.String(), 32)
					if err != nil {
						logErr("the value's type is not float32 in row :%d, column: %s, file: %s", nShowRowNum, pColHead.Name, strXlsxName)
						return false
					}
				}
				break
			case "float64":
				{
					_, err := strconv.ParseFloat(cell.String(), 64)
					if err != nil {
						logErr("the value's type is not float64 in row :%d, column: %s, file: %s", nShowRowNum, pColHead.Name, strXlsxName)
						return false
					}
				}
				break
			}
		}
	}
	return true
}

//生成目标文件
func (pOwn *sMaker) doProcess() bool {
	if len(gCommand.RofPath) > 0 {
		if pOwn.templateRof() == false {
			return false
		}
	}
	if len(gCommand.JsonPath) > 0 {
		if pOwn.templateJson() == false {
			return false
		}
	}
	if len(gCommand.GoFilePath) > 0 {
		if pOwn.templateGo() == false {
			return false
		}
	}
	if len(gCommand.CsFilePath) > 0 {
		if pOwn.templateCs() == false {
			return false
		}
	}
	return true
}
