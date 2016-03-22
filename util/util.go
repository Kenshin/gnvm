package util

import (

	// lib
	. "github.com/Kenshin/cprint"
	"github.com/Kenshin/curl"

	// go
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

const (
	NODE    = "node.exe"
	GNVM    = "gnvm.exe"
	DIVIDE  = "\\"
	SHASUMS = "SHASUMS256.txt"
	UNKNOWN = "unknown"
	LATEST  = "latest"
)

/*
  Golbal node.exe path
*/
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
  Get node.exe version, usage exec.Command()
*/
func GetNodeVer(path string) (string, error) {
	var newout string
	out, err := exec.Command(path+"node", "--version").Output()
	if err == nil {
		newout = strings.Replace(string(string(out[:])[1:]), "\r\n", "", -1)
	}
	return newout, err
}

/*
  Verify node version format.
  Node version format must be http://semver.org/
*/
func VerifyNodeVer(version string) bool {
	result := true
	version = strings.Split(version, "-")[0]
	version = strings.TrimSpace(version)
	reg, _ := regexp.Compile(`^([0]|[1-9]\d?)(\.([0]|[1-9]\d?)){2}$`)
	if version == UNKNOWN || version == LATEST {
		return true
	} else if format := reg.MatchString(version); !format {
		result = false
	}
	return result
}

/*
  Format node version
  x.xx.xx conver to float64
*/
func FormatNodeVer(version string) float64 {
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

/*
  Get node version level( 0 ~ 4 )
*/
func NodeVerLev(ver float64) (level int) {
	switch {
	case ver <= 0.0500:
		level = 0
	case ver >= 0.0501 && ver <= 0.0612:
		level = 1
	case ver > 0.0612 && ver < 1:
		level = 2
	case ver >= 1 && ver <= 0.331:
		level = 3
	case ver > 0.331:
		level = 4
	}
	return
}

/*
 parse arguments return version, io, suffix and arch
 s support format: <version>-<io>-<arch>, e.g.

 	- x.xx.xx
 	- x.xx.xx-io
 	- x.xx.xx-x86|x64
 	- x.xx.xx-io-x86|x64

   Return:
	- ver    : x.xx.xx
	- io     : true  and false
	- arch   : "386" and "amd64"
	- suffix : "x86" and "x64"  and ""

*/
func ParseNodeVer(s string) (ver string, io bool, arch, suffix string) {
	s = strings.ToLower(s)
	arr := strings.Split(s, "-")
	ver = arr[0]

	if ver == LATEST && len(arr) > 1 {
		P(WARING, "%v parameter not support suffix.\n", s)
		io = false
		arch = runtime.GOARCH
		suffix = ""
		return
	}

	switch len(arr) {
	case 1:
		io = false
	case 2:
		if arr[1] == "io" {
			io = true
		} else if ok, _ := regexp.MatchString(`^x?(86|64)$`, arr[1]); ok {
			io = false
			arch = arr[1]
		}
	case 3:
		if arr[1] != "io" {
			s := fmt.Sprintf("%v format error, second parameter must be '%v'.\n", arr[1], "io")
			panic(s)
		} else {
			io = true
		}
		if ok, _ := regexp.MatchString(`^x?(86|64)$`, arr[2]); !ok {
			s := fmt.Sprintf("%v format error, third parameter must be '%v' or '%v'.\n", arr[1], "x86", "x64")
			panic(s)
		} else {
			arch = arr[2]
		}
	}

	// correction arch
	switch arch {
	case "x86":
		arch = "386"
	case "x64":
		arch = "amd64"
	default:
		//P(WARING, "%v format error, only support %v and %v parameter.\n", ver, "x86", "x64")
		arch = runtime.GOARCH
	}

	// correction suffix
	if arch == runtime.GOARCH {
		suffix = ""
	} else {
		if arch == "386" {
			suffix = "x86"
		} else {
			suffix = "x64"
		}
	}

	return
}

/*
 Get remote latest version from url
*/
func GetLatVer(url string) string {

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
		}
		return false
	}

	if err := curl.ReadLine(res.Body, latestVersion); err != nil && err != io.EOF {
		P(ERROR, "%v Error: %v\n", "gnvm update latest", err)
	}

	return version
}

/*
  Ignore key case and return lowercase value
*/
func EqualAbs(key, value string) string {
	if strings.EqualFold(value, key) && value != key {
		P(WARING, "current value is %v, please use %v.\n", value, key)
		value = key
	}
	return value
}

/*
  Return session environment variable
*/
func IsSessionEnv() (string, bool) {
	env := os.Getenv("GNVM_SESSION_NODE_HOME")
	if env != "" {
		return env, true
	} else {
		return env, false
	}
}

/*
 Return binary arch
*/
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
