package nodehandle

import (
	. "github.com/Kenshin/cprint"
	"regexp"
	"runtime"
	"strconv"
)

/*
type Node struct {
	version string
	exec    string
}

type NPM struct {
	version string
}

type NodeList struct {
	id   int
	date string
	Node
	NPM
}
*/

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
