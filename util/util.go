package util

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const (
	NODE   = "node.exe"
	DIVIDE = "\\"
)

var GlobalNodePath string

func init() {
	GlobalNodePath = getGlobalNodePath()
}

func Exec() {
	//fmt.Println("GlobalNodePath = " + GlobalNodePath)
}

func ConverFloat(str string) (float64, error) {
	args := strings.Split(str, ".")
	var newStr string
	for k, v := range args {
		if k == 0 {
			newStr = string(v) + "."
		} else {
			newStr = newStr + string(v)
		}
	}
	version, err := strconv.ParseFloat(newStr, 64)
	return version, err
}

func getGlobalNodePath() string {
	var path string
	file, err := exec.LookPath(NODE)
	if err != nil {
		path = getCurrentPath()
	} else {
		// relpace "\\node.exe"
		path = strings.Replace(file, DIVIDE+NODE, "", -1)
	}

	// gnvm.exe and node.exe the same path
	if path == "." {
		path = getCurrentPath()
	}

	return path
}

func getCurrentPath() string {
	path, err := os.Getwd()
	if err != nil {
		fmt.Println("Get current path Error: " + err.Error())
		return ""
	}
	return path
}
