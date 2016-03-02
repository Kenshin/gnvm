package nodehandle

import (
	"fmt"
	. "github.com/Kenshin/cprint"
	//"gnvm/util"
	"os"
)

func init() {
	// verify GNVM_SESSION_NODE_HOME exist.
	// if exist, 'gnvm install -g', 'gnvm update -g' 'gnvm use x.xx.xx' can't be use.
}

func Run(param string) {

	// try catch
	defer func() {
		if err := recover(); err != nil {
			msg := fmt.Sprintf("'gnvm session' an error has occurred. please check. \nError: ")
			Error(ERROR, msg, err)
			os.Exit(0)
		}
	}()

	/*
		nodePath := util.GlobalNodePath + DIVIDE + version + DIVIDE + NODE
		// <root>/version/node.exe is exist
		if isDirExist(nodePath) != true {
			P(WARING, "%v folder is not exist from %v, use '%v' get local node.exe list. See '%v'.\n", version, rootPath, "gnvm ls", "gnvm help ls")
			return false
		}
	*/

	fmt.Println("asdfasdfasdf " + param)

	// session.cmd is exist? exist quit, not exist, create session.cmd
}
