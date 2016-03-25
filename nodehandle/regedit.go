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

const NODE_HOME, PATH = "NODE_HOME2", "path"

var nodehome, noderoot string

func init() {
	noderoot = config.GetConfig(config.NODEROOT)
	nodehome = os.Getenv(NODE_HOME)
	if nodehome == "" && config.GetConfig(config.GLOBAL_VERSION) == config.UNKNOWN {
		P(NOTICE, "not found environment variable '%v', please use '%v'. See '%v'.\n", NODE_HOME, "gnvm reg noderoot", "gnvm help reg")
	}
}

func Reg(s string) {
	prompt := "n"

	if s != "noderoot" {
		P(ERROR, "parameter %v error, only support %v, please check your input. See '%v'.\n", s, "noderoot", "gnvm help reg")
		return
	}

	P(WARING, "tis is the %v, need %v permission, please note!\n", "experimental function", "Administrator")
	if nodehome != "" {
		P(NOTICE, "current environment variable %v is %v\n", NODE_HOME, nodehome)
	}
	P(NOTICE, "current config %v is %v\n", "noderoot", noderoot)
	P(NOTICE, "set environment variable %v new value is %v [Y/n]? ", NODE_HOME, noderoot)

	fmt.Scanf("%s\n", &prompt)
	prompt = strings.ToLower(prompt)

	if prompt == "y" {
		if _, err := regAdd(NODE_HOME, noderoot); err == nil {
			if arr, err := regQuery(PATH); err == nil {
				prompt = "n"
				P(NOTICE, "if add environment variable %v to %v [Y/n]? ", NODE_HOME, PATH)
				fmt.Scanf("%s\n", &prompt)

				prompt = strings.ToLower(prompt)
				if prompt == "y" {
					regval := ""
					if len(arr) > 0 {
						regval = ";" + arr[0].Value
					}
					if _, err := regAdd(PATH, noderoot+regval); err != nil {
						P(ERROR, "set environment variable %v failed. Error: %v", PATH, err.Error())
					}
				} else {
					P(NOTICE, "oeration has been cancelled.")
				}
			} else if err != nil {
				P(ERROR, "serch environment variable %v failed. Error: %v", PATH, err.Error())
			}
		} else {
			P(ERROR, "add environment variable %v failed. Error: %v", NODE_HOME, err.Error())
		}
	} else {
		P(NOTICE, "operation has been cancelled.")
	}
}

func regAdd(key, value string) ([]regedit.Reg, error) {
	reg := regedit.New(regedit.Add, regedit.HKCU, "\\Environment")
	regcmd := reg.Add(regedit.Reg{key, regedit.Types[regedit.SZ], value})
	return regcmd.Exec()
}

func regQuery(key string) ([]regedit.Reg, error) {
	reg := regedit.New(regedit.Query, regedit.HKCU, "\\Environment")
	regcmd := reg.Search(regedit.Reg{Key: key})
	return regcmd.Exec()
}
