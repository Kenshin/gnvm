package util

import (

	// go
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"

	// local
	"gnvm/util/curl"
	"gnvm/util/p"
)

const (
	NODE    = "node.exe"
	DIVIDE  = "\\"
	SHASUMS = "SHASUMS.txt"
)

var GlobalNodePath string

func init() {
	GlobalNodePath = getGlobalNodePath()
	p.Exec()
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

func GetNodeVersion(path string) (string, error) {
	var newout string
	out, err := exec.Command(path+"node", "--version").Output()
	//string(out[:]) bytes to string
	if err == nil {
		// replace \r\n
		newout = strings.Replace(string(string(out[:])[1:]), "\r\n", "", -1)
	}
	return newout, err
}

func GetLatestVersion(url string) string {

	var version string

	// curl
	code, res, _ := curl.Get(url)
	if code != 0 {
		return ""
	}
	// close
	defer res.Body.Close()

	// set buff
	buff := bufio.NewReader(res.Body)

	for {
		// set line
		line, err := buff.ReadString('\n')

		if line != "" {

			args1 := strings.Split(line, "  ")
			if len(args1) < 2 {
				fmt.Printf("Error: URL [%v] format error, please change registry. See 'gnvm help config'.\n", url)
				break
			}

			args2 := strings.Split(args1[1], "-")
			if len(args2) < 2 {
				fmt.Printf("Error: URL [%v] format error, please change registry. See 'gnvm help config'.\n", url)
				break
			}

			if len(args2[1]) < 2 {
				fmt.Printf("Error: URL [%v] format error, please change registry. See 'gnvm help config'.\n", url)
				break
			}

			// set version
			version = args2[1][1:]
			break
		}

		// when EOF or err break
		if err != nil || err == io.EOF {
			break
		}

	}

	return version

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
