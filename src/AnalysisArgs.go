package main

//eg. xlsx2rof.exe -xlsx ./xlsx -rof ./rof -json ./json -go ./go -cs ./cs
func analysisArgs(args []string) bool {
	var curFunc dealCommand
	for i := 0; i < len(args); i++ {
		parm := args[i]
		if parm[0] == '-' {
			if curFunc != nil {
				if curFunc("") == false {
					return false
				}
			}
			curFunc = checkCommand(parm)
			if curFunc == nil {
				logErr("illegal command:" + parm)
				return false
			}
		} else {
			if curFunc == nil {
				logErr("illegal command:" + parm)
				return false
			}
			if curFunc(parm) == false {
				return false
			}
			curFunc = nil
		}
	}

	if curFunc != nil {
		if curFunc("") == false {
			return false
		}
	}

	return true
}

func checkCommand(aCommand string) dealCommand {
	switch aCommand {
	case "-xlsx":
		return dealCommandXlsx
	case "-rof":
		return dealCommandRof
	case "-json":
		return dealCommandJson
	case "-go":
		return dealCommandGO
	case "-cs":
		return dealCommandCS
	case "-h":
		return dealCommandHelp
	}
	return nil
}

type dealCommand func(arg string) bool

func dealCommandXlsx(arg string) bool {
	if len(arg) <= 0 {
		logErr("the command -xlsx lack of arg.")
		return false
	}
	if arg[len(arg)-1] != '/' {
		gCommand.XlsxPath = arg + "/"
	} else {
		gCommand.XlsxPath = arg
	}
	return true
}

func dealCommandRof(arg string) bool {
	if len(arg) <= 0 {
		logErr("the command -rof lack of arg.")
		return false
	}
	if arg[len(arg)-1] != '/' {
		gCommand.RofPath = arg + "/"
	} else {
		gCommand.RofPath = arg
	}
	return true
}

func dealCommandJson(arg string) bool {
	if len(arg) <= 0 {
		logErr("the command -json lack of arg.")
		return false
	}
	if arg[len(arg)-1] != '/' {
		gCommand.JsonPath = arg + "/"
	} else {
		gCommand.JsonPath = arg
	}
	return true
}

func dealCommandGO(arg string) bool {
	if len(arg) <= 0 {
		logErr("the command -go lack of arg.")
		return false
	}
	if arg[len(arg)-1] != '/' {
		gCommand.GoFilePath = arg + "/"
	} else {
		gCommand.GoFilePath = arg
	}
	return true
}

func dealCommandCS(arg string) bool {
	if len(arg) <= 0 {
		logErr("the command -cs lack of arg.")
		return false
	}
	if arg[len(arg)-1] != '/' {
		gCommand.CsFilePath = arg + "/"
	} else {
		gCommand.CsFilePath = arg
	}
	return true
}

func dealCommandHelp(arg string) bool {
	log("-xlsx : [path] input .xlsx files floder path")
	log("-rof : optional command.  [path] output .bytes files floder path")
	log("-json : optional command. [path] output .json files floder path")
	log("-go : optional command. [path] output .go files floder path")
	log("-cs : optional command. [path] output .cs files floder path")
	log("-h : show help")
	return false
}
