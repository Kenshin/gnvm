package nodehandle

import (

	// lib
	. "github.com/Kenshin/cprint"

	// go
	"fmt"
	"os"
	"strings"

	// local
	"gnvm/config"
	"gnvm/util"
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

	if nodehome != "" {
		P(NOTICE, "current environment variable %v is %v\n", NODE_HOME, nodehome)
	}

	P(NOTICE, "current config %v is %v\n", "noderoot", noderoot)
	P(NOTICE, "set environment variable %v new value is %v [Y/n]? ", NODE_HOME, noderoot)

	fmt.Scanf("%s\n", &prompt)
	prompt = strings.ToLower(prompt)

	if prompt == "y" {
		if _, err := regAdd(NODE_HOME, noderoot); err == nil {
			if arr, err := regQuery(PATH); err == nil && len(arr) == 1 {
				regval := arr[0]

				prompt = "n"
				P(NOTICE, "if add environment variable %v to %v [Y/n]? ", NODE_HOME, PATH)
				fmt.Scanf("%s\n", &prompt)

				prompt = strings.ToLower(prompt)
				if prompt == "y" {
					if _, err := regAdd(PATH, noderoot+";"+regval.Value); err != nil {
						fmt.Println("adfasdfadfasfd")
					}
				}
			}
		}
	}
}

func regAdd(key, value string) ([]util.Reg, error) {
	reg := util.Regedit{util.Actions[util.Add], util.Fields[util.HKCU] + "\\Environment", key, util.Types[util.SZ], value}
	regcmd := reg.Add()
	return regcmd.Exec()
}

func regQuery(key string) ([]util.Reg, error) {
	reg := util.Regedit{Action: util.Actions[util.Query], Field: util.Fields[util.HKCU] + "\\Environment", Key: key}
	regcmd := reg.Search()
	return regcmd.Exec()
}
