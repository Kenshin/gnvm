package util

import (
	"fmt"
	"os/exec"
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
	fmt.Println("GlobalNodePath = " + GlobalNodePath)
}

func getGlobalNodePath() string {
	var path string
	file, err := exec.LookPath(NODE)
	if err != nil {
		path = "root"
	} else {
		// relpace "\\node.exe"
		path = strings.Replace(file, DIVIDE+NODE, "", -1)
	}

	// gnvm.exe and node.exe the same path
	if path == "." {
		path = "root"
	}

	return path
}