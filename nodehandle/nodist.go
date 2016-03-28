package nodehandle

import (
	// go
	"fmt"
	"regexp"
	"strconv"
	"strings"

	// local
	"gnvm/util"
)

type (
	Node struct {
		Version string
		Exec    string
	}

	NPM struct {
		Version string
	}

	NodeList struct {
		ID   int
		Date string
		Node
		NPM
	}

	NL map[string]NodeList
)

var sorts []string

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

func (this NL) Filter(idx int, value map[string]interface{}, regx *regexp.Regexp) (NodeList, bool) {
	isfilter := false
	ver, _ := value["version"].(string)
	if ok := regx.MatchString(ver[1:]); ok {
		date, _ := value["date"].(string)
		npm, _ := value["npm"].(string)
		if npm == "" {
			npm = "[x]"
		}
		exe := filter(ver[1:])
		this[ver] = NodeList{idx, date, Node{ver, exe}, NPM{npm}}
		isfilter = true
	}
	return this[ver], isfilter
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

func filter(version string) (exec string) {
	switch util.GetNodeVerLev(util.FormatNodeVer(version)) {
	case 0:
		exec = "[x]"
	case 1:
		exec = "x86"
	default:
		exec = "x86 x64"
	}
	return
}

func format(value string, max int) string {
	if len(value) > max {
		max = len(value)
	}
	newValue := strings.Repeat(" ", max-len(value))
	return value + newValue
}
