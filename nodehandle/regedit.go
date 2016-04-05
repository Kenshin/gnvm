package nodehandle

import (

	// lib
	. "github.com/Kenshin/cprint"
	"github.com/Kenshin/regedit"

	// go
	"fmt"
	"os"
	"strings"

	// local
	"gnvm/config"
)

const NODE_HOME, PATH = "NODE_HOME", "Path"

var nodehome, noderoot string

func init() {
	noderoot = config.GetConfig(config.NODEROOT)
	nodehome = os.Getenv(NODE_HOME)
	if nodehome == "" && config.GetConfig(config.GLOBAL_VERSION) == config.UNKNOWN {
		P(NOTICE, "not found environment variable '%v', please use '%v'. See '%v'.\n", NODE_HOME, "gnvm reg noderoot", "gnvm help reg")
	}
}

/*
 Regedit

 Param:
 	- s: olny support 'noderoot'

*/
func Reg(s string) {
	prompt := "n"

	P(WARING, "this command is %v, need %v permission, please note!\n", "experimental function", "Administrator")
	if nodehome != "" {
		P(NOTICE, "current environment variable %v is %v\n", NODE_HOME, nodehome)
	}
	P(NOTICE, "current config %v is %v\n", "noderoot", noderoot)
	P(NOTICE, "set environment variable %v is %v [Y/n]? ", NODE_HOME, noderoot)

	fmt.Scanf("%s\n", &prompt)
	prompt = strings.ToLower(prompt)

	if prompt == "y" {
		if add(NODE_HOME, noderoot) == nil {
			if arr, err := query(PATH); err == nil {
				prompt = "n"
				P(NOTICE, "add environment variable %v to %v [Y/n]? ", NODE_HOME, PATH)
				fmt.Scanf("%s\n", &prompt)

				prompt = strings.ToLower(prompt)
				if prompt == "y" {
					regval := ""
					if len(arr) > 0 {
						regval = ";" + arr[0].Value
					}
					add(PATH, noderoot+regval)
				} else {
					P(NOTICE, "operation has been cancelled.")
				}
			}
		}
	} else {
		P(NOTICE, "operation has been cancelled.")
	}
}

func add(key, value string) (err error) {
	reg := regedit.New(regedit.Add, regedit.HKCU, "\\Environment")
	regcmd := reg.Add(regedit.Reg{key, regedit.Types[regedit.SZ], value})
	if _, err = regcmd.Exec(); err != nil {
		P(ERROR, "add environment variable %v failed. Error: %v", NODE_HOME, err.Error())
	}
	return err
}

func query(key string) (regs []regedit.Reg, err error) {
	reg := regedit.New(regedit.Query, regedit.HKCU, "\\Environment")
	regcmd := reg.Search(regedit.Reg{Key: key})
	if regs, err = regcmd.Exec(); err != nil {
		P(ERROR, "search environment variable %v failed. Error: %v", PATH, err.Error())
	}
	return regs, err
}
