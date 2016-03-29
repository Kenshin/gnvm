package nodehandle

import (

	// lib
	. "github.com/Kenshin/cprint"
	"github.com/Kenshin/curl"
	"github.com/bitly/go-simplejson"

	// go
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	// local
	"gnvm/config"
	"gnvm/util"
)

const (
	LATNPMURL   = "https://raw.githubusercontent.com/npm/npm/master/package.json"
	NPMTAOBAO   = "http://npm.taobao.org/mirrors/npm/"
	NPMDEFAULT  = "https://github.com/npm/npm/releases/"
	ZIP         = ".zip"
	NODEMODULES = "node_modules"
	NPMBIN      = "bin"
	NPMCOMMAND1 = "npm"
	NPMCOMMAND2 = "npm.cmd"
)

func InstallNPM(version string) {
	// try catch
	defer func() {
		if err := recover(); err != nil {
			P(ERROR, "%v an error has occurred, Error is %v \n", "gnvm npm", err)
			os.Exit(0)
		}
	}()

	version = strings.ToLower(version)
	prompt := "n"
	switch version {
	case util.LATEST:
		remote, local := getLatNPMVer(), getGlobalNPMVer()
		v1, v2 := util.FormatNodeVer(remote), util.FormatNodeVer(local)
		cp := CP{Red, false, None, false, "="}
		switch {
		case v1 > v2:
			cp.Value = ">"
			P(WARING, "npm remote latest version %v %v local latest version %v.\n", remote, cp, local)
			P(NOTICE, "is update local npm version [Y/n]? ")
			fmt.Scanf("%s\n", &prompt)
			prompt = strings.ToLower(prompt)
			if prompt == "y" {
				downloadNpm(remote)
			} else {
				P(NOTICE, "you need use '%v' update local version. \n", "npm install -g npm")
			}
		case v1 < v2:
			cp.Value = "<"
			P(WARING, "npm remote latest version %v %v local latest version %v.\n", remote, cp, local)
		case v1 == v2:
			P(WARING, "npm remote latest version %v %v local latest version %v.\n", remote, cp, local)
		}
	case util.GLOBAL:
	default:
		P(ERROR, "'%v' param only support [%v] [%v] [%v], please check your input. See '%v'.\n", "gnvm npm", "latest", "global", "x.xx.xx", "gnvm help npm")
		return
	}
}

func UninstallNPM() {
	fmt.Println("UnInstall")
}

/*
 Get Latest NPM version
*/
func getLatNPMVer() string {
	_, res, err := curl.Get(LATNPMURL)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	json, err := simplejson.NewJson(body)
	if err != nil {
		panic(err)
	}
	ver, _ := json.Get("version").String()
	return ver
}

/*
 Get global( local ) NPM version
*/
func getGlobalNPMVer() string {
	out, err := exec.Command(rootPath+util.NPM, "-v").Output()
	if err != nil {
		P(WARING, "current path %v not exist npm.\n", rootPath)
		return util.UNKNOWN
	}
	return strings.TrimSpace(string(out[:]))
}

/*
 Download and unzip npm.zip
*/
func downloadNpm(version string) {
	version = "v" + version + ZIP
	url := NPMTAOBAO + version
	if config.GetConfig(config.REGISTRY) != config.TAOBAO {
		url = NPMDEFAULT + version
	}
	if dl, errs := curl.New(url); len(errs) > 0 {
		err := errs[0]
		P(ERROR, "%v an error has occurred, Error is %v \n", "gnvm npm latest", err)
		return
	} else {
		fmt.Println(dl)
	}
}

/*
 Create node_modules folder
*/
func MkNPM(path, zip string) {
	//dest := config.GetConfig(config.NODEROOT) + util.DIVIDE + NODEMODULES
	dest := path + util.DIVIDE + NODEMODULES
	zip = path + util.DIVIDE + zip
	npm := dest + util.DIVIDE + util.NPM
	npmbin := npm + util.DIVIDE + NPMBIN

	// verify node_modules exist
	if !isDirExist(dest) {
		if err := os.Mkdir(dest, 0755); err != nil {
			P(ERROR, "create %v foler error, Error: %v\n", dest, err.Error())
			return
		} else {
			P(NOTICE, "%v folder create success.\n", dest)
		}
	}

	// verify node_modules/npm exist
	if isDirExist(npm) {
		if err := os.RemoveAll(npm); err != nil {
			P(ERROR, "remove %v folder Error: %v.\n", npm, err.Error())
			return
		}
	}

	// unzip
	if code, err := unzip(zip, dest); err != nil {
		fmt.Println(code)
		fmt.Println(err)
	} else {
		if err := os.Rename(dest+util.DIVIDE+"npm-3.8.5", dest+util.DIVIDE+util.NPM); err != nil {
			P(ERROR, "unzip fail, Error: %v", err.Error())
			return
		} else {
			// copy <root>\node_modules\npm\bin npm and npm.cmd to <root>\
			if err := copyFile(npmbin, path, NPMCOMMAND1); err != nil {
				P(ERROR, "copy %v to %v faild, Error: %v \n", npmbin, path)
				return
			}
			if err := copyFile(npmbin, path, NPMCOMMAND2); err != nil {
				P(ERROR, "copy %v to %v faild, Error: %v \n", npmbin, path)
				return
			}
			// remove download zip file
			if err := os.RemoveAll(zip); err != nil {
				P(ERROR, "remove %v folder Error: %v.\n", npm, err.Error())
				return
			}
			P(NOTICE, "unzip complete.\n")
		}
	}
}

/*
  Unzip file

  Param:
    - path: zip file path
    - dest: unzip dest folder

  Return:
    - error
    - code
        - -1: open  zip file error
        - -2: open  file error
        - -3: write file error
        - -4: copy  file error
*/
func unzip(path, dest string) (string, error) {
	unzip, err := zip.OpenReader(path)
	if err != nil {
		return "-1", err
	}
	defer unzip.Close()
	idx, root := 0, ""
	for _, file := range unzip.File {
		rc, err := file.Open()
		if err != nil {
			return "-2", err
		}
		defer rc.Close()
		if idx == 0 {
			root = file.Name
		}
		path = filepath.Join(dest, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
		} else {
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
			if err != nil {
				return "-3", err
			}
			defer f.Close()
			if _, err := io.Copy(f, rc); err != nil {
				return "-4", err
			}
		}
		idx++
	}
	return root, nil
}

/*
 Copy file from src to dest
*/
func copyFile(src, dst, name string) (err error) {
	src = src + util.DIVIDE + name
	dst = dst + util.DIVIDE + name
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
