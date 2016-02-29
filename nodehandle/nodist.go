package nodehandle

import (
	"fmt"
	. "github.com/Kenshin/cprint"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

type Node struct {
	Version string
	Exec    string
}

type NPM struct {
	Version string
}

type NodeList struct {
	ID   int
	Date string
	Node
	NPM
}

type NL map[string]NodeList

func GetNodePath(version string) string {
	reg1, _ := regexp.Compile(`^0\.`)
	// version format include: 0.xx.xx and ^0.xx.xx
	// if 0.xx.xx get 0.xx else get ^0.xx * 10.0
	regstr := `(^\d+\.\d+)`
	multiple := 10.0
	if format := reg1.MatchString(version); format {
		regstr = `(\d)+.(\d)+$`
		multiple = 1.0
	}
	reg2, _ := regexp.Compile(regstr)
	str := reg2.FindString(version)
	ver, _ := strconv.ParseFloat(str, 64)
	// get ture version
	ver = ver * multiple
	path := "/"
	switch {
	case ver <= 5:
		P(ERROR, "downlaod node.exe version: %v, not node.exe. See '%v'.\n", version, "gnvm help install")
	case ver > 5 && ver <= 6.12:
		P(WARING, "downlaod node.exe version: %v, not x64 node.exe.\n", version)
	case ver > 6.12 && ver <= 12.10:
		if runtime.GOARCH == "amd64" {
			path = "/x64/"
		}
	default:
		if runtime.GOARCH == "amd64" {
			path = "/win-x64/"
		} else {
			path = "/win-x86/"
		}
	}
	return "v" + version + path
}

func filter(files []interface{}) string {
	exec := ""
	reg, _ := regexp.Compile(`x(86|64)`)
	for _, file := range files {
		if ok, err := regexp.MatchString("^win-x(86|64)-exe", file.(string)); ok && err == nil {
			exec += reg.FindString(file.(string)) + " "
		}
	}
	if exec == "" {
		exec = "[x]"
	}
	return exec
}

func format(value string, max int) string {
	if len(value) > max {
		max = len(value)
	}
	newValue := strings.Repeat(" ", max-len(value))
	return value + newValue
}

func (nl NL) New(idx int, value map[string]interface{}) NodeList {
	ver, _ := value["version"].(string)
	date, _ := value["date"].(string)
	npm, _ := value["npm"].(string)
	if npm == "" {
		npm = "[x]"
	}
	exe := filter(value["files"].([]interface{}))
	nl[ver] = NodeList{idx, date, Node{ver, exe}, NPM{npm}}
	return nl[ver]
}

func (nl *NL) Print(nodeist NodeList) {
	msg := fmt.Sprintf("id: %v date: %v node version: %v os support: %v npm version: %v", nodeist.ID, nodeist.Date, nodeist.Node.Version, nodeist.Node.Exec, nodeist.NPM.Version)
	fmt.Println(msg)
}

func (nl NL) Detail() {
	fmt.Println("No.   date         node ver    exec      npm ver  ")
	fmt.Println("--------------------------------------------------")
	for _, value := range nl {
		id := format(strconv.Itoa(value.ID), 6)
		date := format(value.Date, 13)
		ver := format(value.Node.Version, 12)
		exe := format(value.Node.Exec, 10)
		npm := format(value.NPM.Version, 9)
		fmt.Println(id + date + ver + exe + npm)
	}
}
