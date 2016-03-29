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

/*
func InstallNpm() {

	// try catch
	defer func() {
		if err := recover(); err != nil {
			if strings.HasPrefix(err.(string), "CURL Error:") {
				fmt.Printf("\n")
			}
			Error(ERROR, "'gnvm install npm' an error has occurred. \nError: ", err)
			os.Exit(0)
		}
	}()

	out, err := exec.Command(rootPath+"npm", "--version").Output()
	if err == nil {
		P(WARING, "current path %v exist npm, version is %v", rootPath, string(out[:]), "\n")
		return
	}

	url := config.GetConfig(config.REGISTRY) + "npm"

	// get
	code, res, _ := curl.Get(url)
	if code != 0 {
		return
	}
	// close
	defer res.Body.Close()

	maxTime, _ := time.Parse(TIMEFORMART, TIMEFORMART)
	var maxVersion string

	getNpmVersion := func(content string, line int) bool {
		if strings.Index(content, `<a href="`) == 0 && strings.Contains(content, ".zip") {

			// parse
			newLine := strings.Replace(content, `<a href="`, "", -1)
			newLine = strings.Replace(newLine, `</a`, "", -1)
			newLine = strings.Replace(newLine, `">`, " ", -1)

			// e.g. npm-1.3.9.zip npm-1.3.9.zip> 23-Aug-2013 21:14 1535885
			orgArr := strings.Fields(newLine)

			// e.g. npm-1.3.9.zip
			version := orgArr[0:1][0]

			// e.g. 23-Aug-2013 21:14
			sTime := strings.Join(orgArr[2:len(orgArr)-1], " ")

			// bubble sort
			if t, err := time.Parse(TIMEFORMART, sTime); err == nil {
				if t.Sub(maxTime).Seconds() > 0 {
					maxTime = t
					maxVersion = version
				}
			}
		}
		return false
	}

	if err := curl.ReadLine(res.Body, getNpmVersion); err != nil && err != io.EOF {
		P(ERROR, "parse npm version Error: %v, from %v\n", err, url)
		return
	}

	if maxVersion == "" {
		P(ERROR, "get npm version fail from %v, please check. See '%v'.\n", url, "gnvm help config")
		return
	}

	P(NOTICE, "the latest version is %v from %v.\n", maxVersion, config.GetConfig(config.REGISTRY))

	// download zip
	zipPath := os.TempDir() + util.DIVIDE + maxVersion
	if code := downloadNpm(maxVersion); code == 0 {

		P(DEFAULT, "Start unarchive file %v.\n", maxVersion)

		//unzip(maxVersion)

		if err := zip.UnarchiveFile(zipPath, config.GetConfig(config.NODEROOT), nil); err != nil {
			panic(err)
		}

		P(DEFAULT, "End unarchive.\n")
	}

	// remove temp zip file
	if err := os.RemoveAll(zipPath); err != nil {
		P(ERROR, "remove zip file fail from %v, Error: %v.\n", zipPath, err.Error())
	}

}

func UninstallNpm() {

	// try catch
	defer func() {
		if err := recover(); err != nil {
			Error(ERROR, "'gnvm uninstall npm' an error has occurred. please check. \nError: ", err)
			os.Exit(0)
		}
	}()

	removeFlag := true

	if !isDirExist(rootPath+"npm.cmd") && !isDirExist(rootPath+"node_modules"+util.DIVIDE+"npm") {
		P(WARING, "%v not exist %v.\n", rootPath, "npm.cmd")
		return
	}

	// remove npm.cmd
	if err := os.RemoveAll(rootPath + "npm.cmd"); err != nil {
		removeFlag = false
		P(ERROR, "remove %v file fail from %v, Error: %v.\n", "npm.cmd", rootPath, err.Error())
	}

	// remove node_modules/npm
	if err := os.RemoveAll(rootPath + "node_modules" + util.DIVIDE + "npm"); err != nil {
		removeFlag = false
		P(ERROR, "remove %v folder fail from %v, Error: %v.\n", "npm", rootPath+"node_modules", err.Error())
	}

	if removeFlag {
		P(DEFAULT, "npm uninstall success from %v.\n", rootPath)
	}
}

func downloadNpm(version string) int {
   // set url
   url := config.GetConfig(config.REGISTRY) + "npm/" + version
   // download
   if code := curl.New(url, version, os.TempDir()+DIVIDE+version); code != 0 {
       return code
   }
	return 0
}
*/
