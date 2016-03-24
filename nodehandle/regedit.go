package nodehandle

import (

	// lib
	. "github.com/Kenshin/cprint"

	// go
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	// local
	"gnvm/config"
)

const NODE_HOME = "NODE_HOME2"

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
	P(NOTICE, "if set environment variable %v new value is %v [Y/n]? ", NODE_HOME, noderoot)

	fmt.Scanf("%s", &prompt)
	prompt = strings.ToLower(prompt)

	if prompt == "y" {
		if err := regAdd(NODE_HOME, noderoot); err == nil {
			if path := regQuery("path"); path != "" {
				regAdd("path", noderoot+";"+path)
			}
		}
	}
}

func regAdd(key, value string) (err error) {
	regPath := "HKEY_CURRENT_USER\\Environment"
	cmd := exec.Command("cmd", "/c", "reg", "add", regPath, "/v", key, "/t", "REG_SZ", "/d", value)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		P(ERROR, "set failed. Error: %s\n", err.Error())
	}
	return err
}

func regQuery(value string) string {
	regPath := "HKEY_CURRENT_USER\\Environment"
	if out, err := exec.Command("cmd", "/c", "reg", "query", regPath, "/v", value).Output(); err != nil {
		P(ERROR, "get path failed. Error: %s\n", err.Error())
	} else {
		buff := bytes.NewBuffer(out)
		line := 1
		for {
			content, err := buff.ReadString('\n')
			content = strings.TrimSpace(content)
			if err != nil || err == io.EOF {
				break
			}
			if line == 3 {
				content = strings.Replace(content, "REG_SZ", "", -1)
				content = strings.Replace(content, value, "", -1)
				content = strings.TrimSpace(content)
				return content
			}
			line++
		}
	}
	return ""
}
