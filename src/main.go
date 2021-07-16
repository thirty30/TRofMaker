package main

var gCommandItems []*sCommandItem
var gTables []*sTableInfo //表信息

func main() {
	gCommandItems = make([]*sCommandItem, 0, 16)
	gCommandItems = append(gCommandItems, &sCommandItem{Cmd: "-xlsx", Builder: new(sExcelBuilder), Parms: make([]string, 0, 2), CanExecute: false})
	gCommandItems = append(gCommandItems, &sCommandItem{Cmd: "-rof", Builder: new(sRofBuilder), Parms: make([]string, 0, 2), CanExecute: false})
	gCommandItems = append(gCommandItems, &sCommandItem{Cmd: "-go", Builder: new(sGoBuilder), Parms: make([]string, 0, 2), CanExecute: false})
	gCommandItems = append(gCommandItems, &sCommandItem{Cmd: "-cs", Builder: new(sCsBuilder), Parms: make([]string, 0, 2), CanExecute: false})
	gCommandItems = append(gCommandItems, &sCommandItem{Cmd: "-json", Builder: new(sJsonBuilder), Parms: make([]string, 0, 2), CanExecute: false})

	gTables = make([]*sTableInfo, 0, 1024)

	//初始化控制台输出颜色
	initConsoleColor()

	//解析命令
	if analysisArgs() == false {
		return
	}

	//初始化builder参数
	for _, v := range gCommandItems {
		if v.Cmd == "-xlsx" && v.CanExecute == false {
			logErr("lack necessary option -xlsx.")
			return
		}
		if v.CanExecute == true && v.Builder.init(v.Parms) == false {
			return
		}
	}

	//处理文件
	for _, v := range gCommandItems {
		if v.CanExecute == false {
			continue
		}
		if v.Builder.build() == false {
			return
		}
	}

	log("[SUCCESS] Generate completely!")
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

func isRofNameRepeated(aRofName string) (bool, string) {
	for _, v := range gTables {
		if v.RofName == aRofName {
			return true, (v.Dir + v.FileName)
		}
	}
	return false, ""
}
