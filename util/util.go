package util

import (

	// lib
	. "github.com/Kenshin/cprint"
	"github.com/Kenshin/curl"

	// go
	"encoding/hex"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

const (
	NODE    = "node.exe"
	GNVM    = "gnvm.exe"
	DIVIDE  = "\\"
	SHASUMS = "SHASUMS256.txt"
	UNKNOWN = "unknown"
)

var GlobalNodePath string

func init() {

	// try catch
	defer func() {
		if err := recover(); err != nil {
			Error(ERROR, "initialize gnvm.exe an error has occurred. please check. \nError: ", err)
			os.Exit(0)
		}
	}()

	GlobalNodePath = getGlobalNodePath()
}

/*
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
*/

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

	latestVersion := func(content string, line int) bool {
		if content != "" && line == 1 {
			reg, _ := regexp.Compile(`\d(\.\d){2}`)
			version = reg.FindString(content)
			/*
				args1 := strings.Split(content, "  ")
				if len(args1) < 2 {
					P(ERROR, "URL %v format error, please change registry. See '%v'.\n", url, "gnvm help config")
					return true
				}

				args2 := strings.Split(args1[1], "-")
				if len(args2) < 2 {
					P(ERROR, "URL %v format error, please change registry. See '%v'.\n", url, "gnvm help config")
					return true
				}

				if len(args2[1]) < 2 {
					P(ERROR, "URL %v format error, please change registry. See '%v'.\n", url, "gnvm help config")
					return true
				}

				version = args2[1][1:]
			*/
		}

		return false

	}

	if err := curl.ReadLine(res.Body, latestVersion); err != nil && err != io.EOF {
		P(ERROR, "%v Error: %v\n", "gnvm update latest", err)
	}

	return version
}

func VerifyNodeVersion(version string) bool {
	result := true
	version = strings.Split(version, "-")[0]
	version = strings.TrimSpace(version)
	reg, _ := regexp.Compile(`^([0]|[1-9]\d?)(\.([0]|[1-9]\d?)){2}$`)
	if version == UNKNOWN {
		return true
	} else if format := reg.MatchString(version); !format {
		result = false
	}
	return result
}

func EqualAbs(key, value string) string {
	if strings.EqualFold(value, key) && value != key {
		P(WARING, "current value is %v, please use %v.\n", value, key)
		value = key
	}
	return value
}

func IsSessionEnv() (string, bool) {
	env := os.Getenv("GNVM_SESSION_NODE_HOME")
	if env != "" {
		return env, true
	} else {
		return env, false
	}
}

func getGlobalNodePath() string {
	var path string

	if env, ok := IsSessionEnv(); ok {
		if reg, err := regexp.Compile(`\\([0]|[1-9]\d?)(\.([0]|[1-9]\d?)){2}\\$`); err == nil {
			ver := reg.FindString(env)
			path = strings.Replace(env, ver, "", -1)
		}
		return path
	}

	file, err := exec.LookPath(NODE)
	if err != nil {
		if file, err := exec.LookPath(GNVM); err != nil {
			path = getCurrentPath()
		} else {
			path = strings.Replace(file, DIVIDE+GNVM, "", -1)
		}
	} else {
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
		panic("get current path Error: " + err.Error())
	}
	return path
}

func Arch(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	bit32, _ := hex.DecodeString("504500004C")
	bit64, _ := hex.DecodeString("504500006486")
	byte := make([]byte, 5)
	j := 0
	for {
		j++
		byte = byte[:cap(byte)]
		n, err := f.Read(byte)
		if err == io.EOF {
			return "x64", nil
		}
		byte = byte[:n]
		if string(byte[:]) == string(bit32[:]) {
			return "x86", nil
		}
		if string(byte[:]) == string(bit64[1:]) || string(byte[:]) == string(bit64[:len(bit32)]) {
			return "x64", nil
		}
		if j == 60 {
			return "x64", nil
		}
	}
	return "x64", nil
}
