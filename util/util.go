package util

import (

	// lib
	. "github.com/Kenshin/cprint"
	"github.com/Kenshin/curl"

	// go
	"encoding/hex"
	"errors"
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
	IOJS    = "iojs.exe"
	SHASUMS = "SHASUMS256.txt"
	UNKNOWN = "unknown"
	LATEST  = "latest"
	GLOBAL  = "global"
	NPM     = "npm"
)

var DIVIDE = string(os.PathSeparator)

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
	VaildPath(&path)
	out, err := exec.Command(path+NODE, "--version").Output()
	if err == nil {
		return strings.TrimSpace(string(out[1:])), nil
	}
	return "", err
}

/*
  Verify node version format.
  Node version format must be http://semver.org/
*/
func VerifyNodeVer(version string) bool {
	version = strings.Split(version, "-")[0]
	version = strings.TrimSpace(version)
	reg, _ := regexp.Compile(`^([0]|[1-9]\d?)(\.([0]|[1-9]\d?)){2}$`)
	if version == UNKNOWN || version == LATEST {
		return true
	}
	return reg.MatchString(version)
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
  Format wildcard node version
  Do not allow the *.1.* model
	- `*.*.*`      - wildcard( include x|X )
	- `1.*.*`      - wildcard
	- `0.10.*`     - wildcard
	- `5.9.0`      - {num}.{num}.{num}
	- `\<regexp>\` - regexp
	- latest       - trans to true version

  Return:
	- regexp
	- error
*/
func FormatWildcard(version, url string) (*regexp.Regexp, error) {
	version = strings.ToLower(version)
	version = strings.Replace(version, "x", "*", -1)

	// *.*.* x.x.x X.x.x *.X.x
	reg1 := `^(\*)(\.(\*)){2}$`
	// {num}.*.*
	reg2 := `^(0{1}|[1-9]\d?)(\.\*{1}){2}$`
	// {num}.{num}.*
	reg3 := `^(0{1}\.|[1-9]\d?\.){2}\*$`

	if version == LATEST {
		version = GetLatVer(url)
		return regexp.Compile(version)
	} else if strings.HasPrefix(version, "/") && strings.HasSuffix(version, "/") {
		return regexp.Compile(version[1 : len(version)-1])
	} else if ok := VerifyNodeVer(version); ok {
		return regexp.Compile(version)
	} else if ok, _ := regexp.MatchString(reg1, version); ok {
		return regexp.Compile(`^([0]|[1-9]\d?)(\.([0]|[1-9]\d?)){2}$`)
	} else if ok, _ := regexp.MatchString(reg2, version); ok {
		return regexp.Compile(`^` + strings.Replace(version, ".*", "", -1) + `(\.([0]|[1-9]\d?)){2}$`)
	} else if ok, _ := regexp.MatchString(reg3, version); ok {
		return regexp.Compile(`^` + strings.Replace(version, "*", "", -1) + `([0]|[1-9]\d?)$`)
	} else {
		return nil, errors.New("parameter format error.")
	}
}

/*
 Conver latest to x.xx.xx( include unknown)
*/
func FormatLatVer(latest *string, value string, print bool) {
	if *latest == LATEST {
		*latest = value
		if print {
			P(NOTICE, "current latest version is %v.\n", *latest)
		}
	}
}

/*
  Get node version level( 0 ~ 4 )
  Return
	- 0: no exec
	- 1: only x86 exec
	- 2: x86 and x64 exec, folder is "x64/" and <root>
	- 3: io.js exec, folder is "win-x64/" and "win-x86/"
	- 4: x86 and x64 exec, folder is "win-x64/" and "win-x86/"
*/
func GetNodeVerLev(ver float64) (level int) {
	switch {
	case ver <= 0.0500:
		level = 0
	case ver > 0.0500 && ver <= 0.0612:
		level = 1
	case ver > 0.0612 && ver < 1:
		level = 2
	case ver >= 1 && ver <= 3.0301:
		level = 3
	case ver > 3.0301:
		level = 4
	}
	return
}

/*
 parse arguments return version, io, suffix and arch
 s support format: <version>-<io>-<arch>, e.g.

 	- x.xx.xx
 	- x.xx.xx-x86|x64

   Return:
	- ver    : x.xx.xx
	- iojs   : true  and false
	- arch   : "386" and "amd64"
	- suffix : "x86" and "x64"  and ""
	- err    : includ, "1" "2", "3", "4", "5"

*/
func ParseNodeVer(s string) (ver string, iojs bool, arch, suffix string, err error) {
	arr := strings.Split(strings.ToLower(s), "-")

	// get ver
	ver = arr[0]

	// verify npm
	if ver == NPM {
		err = errors.New("5")
		return
	}

	// verify latest
	if ver == LATEST {
		if len(arr) > 1 {
			P(WARING, "%v parameter not support suffix.\n", s)
		}
		iojs = false
		arch = runtime.GOARCH
		suffix = ""
		return
	}

	// verify ver
	if !VerifyNodeVer(ver) {
		err = errors.New("4")
		return
	}

	switch GetNodeVerLev(FormatNodeVer(ver)) {
	case 0:
		// no exec
		err = errors.New("1")
		return
	case 3:
		// get iojs
		iojs = true
	}

	// get arch
	if len(arr) == 2 {
		if ok, _ := regexp.MatchString(`^x?(86|64)$`, arr[1]); ok {
			arch = arr[1]
		} else {
			err = errors.New("2")
			return
		}
	} else if len(arr) > 2 {
		err = errors.New("3")
		return
	}

	// get arch
	switch arch {
	case "x86":
		arch = "386"
	case "x64":
		arch = "amd64"
	default:
		arch = runtime.GOARCH
	}

	// get suffix
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
			reg, _ := regexp.Compile(`([0]|[1-9]\d?)(\.([0]|[1-9]\d?)){2}`)
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
 Return node.exe real url, e.g.
 	- http://npm.taobao.org/mirrors/node/v5.9.0/win-x64/node.exe
 	- http://npm.taobao.org/mirrors/iojs/v1.0.0/win-x86/iojs.exe
*/
func GetRemoteNodePath(url, version, arch string) (string, error) {
	folder := "/"
	exec := NODE
	level := GetNodeVerLev(FormatNodeVer(version))

	switch level {
	case 0:
		P(ERROR, "downlaod node.exe version %v, not %v. See '%v'.\n", version, "node.exe", "gnvm help install")
		return "", errors.New("Not support version " + version + "download.")
	case 1:
		P(WARING, "downlaod node.exe version %v, not %v node.exe.\n", version, "x64")
	case 2:
		if arch == "amd64" {
			folder = "/x64/"
		}
	default:
		if arch == "amd64" {
			folder = "/win-x64/"
		} else {
			folder = "/win-x86/"
		}
	}

	// when level == 3, exec is "iojs.exe"
	if level == 3 {
		exec = IOJS
	}

	return url + "v" + version + folder + exec, nil
}

/*
 Get node.exe binary arch

   Param:
	- path:    node.exe path

   Return:
	- string: arch, inlcude: 'x86' 'x64'
	- error

*/
func Arch(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	bit32, _ := hex.DecodeString("504500004C")
	//bit64, _ := hex.DecodeString("504500006486")
	bytes, empty := make([]byte, 1), [5]byte{}
	i, j := 0, 0
	for {
		j++
		bytes = bytes[:cap(bytes)]
		n, err := f.Read(bytes)
		if err == io.EOF {
			return "x64", nil
		}
		bytes = bytes[:n]
		if i <= 4 {
			if string(bytes[:]) == string(bit32[i]) {
				empty[i] = bytes[0]
				if string(empty[:]) == string(bit32[:]) {
					return "x86", nil
				}
				i++
			}
		}
		if j == 500 {
			return "x64", nil
		}
	}
	return "x64", nil
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
 Vaild Path, e.g x:\aa\bb\cc to x:\aa\bb\cc\
*/
func VaildPath(path *string) {
	if !IsDirExist(*path) {
		P(WARING, "%v not a vaild directory.\n", *path)
	} else if !strings.HasSuffix(*path, DIVIDE) {
		*path += DIVIDE
	}
}

/*
 Copy file from src to dest

 Param:
 	- src:  copy file path
	- dst:  target file path
	- name: copy file name

 Return:
 	- error
*/
func Copy(src, dst, name string) (err error) {
	src = src + DIVIDE + name
	dst = dst + DIVIDE + name
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}

/*
 Judge path( folder ) exist

 Param:
    - path: valid path e.g. /gnvm/node_modules

 Return:
    - true : exist
    - false: no exit
*/
func IsDirExist(path string) bool {
	file, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	} else {
		return file.IsDir()
	}
	return false
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
