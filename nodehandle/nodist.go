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

var sorts []string

func ParseFloat(version string) float64 {
	reg, _ := regexp.Compile(`\.(\d){0,2}`)
	ver := ""
	arr := reg.FindAllString(version, -1)
	for _, v := range arr {
		v = v[1:]
		if len(v) == 1 {
			ver += "0" + v
		} else if len(v) == 2 {
			ver += v
		}
	}
	reg, _ = regexp.Compile(`^(\d){1,2}\.`)
	prefix := reg.FindString(version)
	ver = prefix + ver
	float64, _ := strconv.ParseFloat(ver, 64)
	return float64
}

func GetNodePath(version string) string {
	ver := ParseFloat(version)
	path := "/"
	switch {
	case ver <= 0.0500:
		P(ERROR, "downlaod node.exe version: %v, not %v. See '%v'.\n", version, "node.exe", "gnvm help install")
	case ver >= 0.0501 && ver <= 0.0612:
		P(WARING, "downlaod node.exe version: %v, not %v node.exe.\n", version, "x64")
	case ver > 0.0612 && ver < 4:
		if runtime.GOARCH == "amd64" {
			path = "/x64/"
		}
	case ver >= 4:
		if runtime.GOARCH == "amd64" {
			path = "/win-x64/"
		} else {
			path = "/win-x86/"
		}
	}
	return "v" + version + path
}

func filter(version string) string {
	ver := ParseFloat(version)
	exec := ""
	switch {
	case ver <= 0.0500:
		exec = "[x]"
	case ver >= 0.0501 && ver <= 0.0612:
		exec = "x86"
	case ver > 0.0612:
		exec = "x86 x64"
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
	exe := filter(ver[1:])
	nl[ver] = NodeList{idx, date, Node{ver, exe}, NPM{npm}}
	return nl[ver]
}

func (nl *NL) Print(nodeist NodeList) {
	msg := fmt.Sprintf("id: %v date: %v node version: %v os support: %v npm version: %v", nodeist.ID, nodeist.Date, nodeist.Node.Version, nodeist.Node.Exec, nodeist.NPM.Version)
	fmt.Println(msg)
}

func (nl NL) IndexBy(key string) {
	sorts = append(sorts, key)
}

func (nl NL) Detail(limit int) {
	table := `+--------------------------------------------------+
| No.   date         node ver    exec      npm ver |
+--------------------------------------------------+`
	if limit == 0 || limit > len(sorts) {
		limit = len(sorts)
	}
	for idx, v := range sorts {
		if idx == 0 {
			fmt.Println(table)
		}
		if idx >= limit {
			break
		}
		value := nl[v]
		id := format(strconv.Itoa(value.ID+1), 6)
		date := format(value.Date, 13)
		ver := format(value.Node.Version[1:], 12)
		exe := format(value.Node.Exec, 10)
		npm := format(value.NPM.Version, 9)
		fmt.Println("  " + id + date + ver + exe + npm)
		if idx == limit-1 {
			fmt.Println("+--------------------------------------------------+")
		}
	}
}
