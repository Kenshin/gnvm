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

type NPMDownload struct {
	root    string
	zipname string
	zippath string
	modules string
	npmpath string
	npmbin  string
}

var npm = new(NPMDownload)

/*
 Crete NPMDownload
*/
func (this *NPMDownload) New(zip string) {
	(*this).root = config.GetConfig(config.NODEROOT)
	(*this).modules = (*this).root + util.DIVIDE + NODEMODULES
	(*this).zipname = zip
	(*this).zippath = (*this).root + util.DIVIDE + (*this).zipname
	(*this).npmpath = (*this).modules + util.DIVIDE + util.NPM
	(*this).npmbin = (*this).npmpath + util.DIVIDE + NPMBIN
}

/*
 Custom Print
*/
func (this *NPMDownload) String() string {
	return fmt.Sprintf("root   = %v \nzipname= %v\nzippath= %v\nmodules= %v\nnpmpath= %v\nnpmbin = %v\n", this.root, this.zipname, this.zippath, this.modules, this.npmpath, this.npmbin)
}

/*
 Remove file
*/
func (this *NPMDownload) Clean(path string) error {
	if isDirExist(path) {
		if err := os.RemoveAll(path); err != nil {
			P(ERROR, "remove %v folder Error: %v.\n", path, err.Error())
			return err
		}
	}
	return nil
}

/*
 Remove <root>/node_modules/npm <root>/npm <root>/npm.cmd
*/
func (this *NPMDownload) CleanAll() {
	paths := [3]string{this.npmpath, this.root + util.DIVIDE + NPMCOMMAND1, this.root + util.DIVIDE + NPMCOMMAND2}
	for _, v := range paths {
		this.Clean(v)
	}
}

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
		ts := dl[0]
		MkNPM(ts.Name)
	}
}

/*
 Create npm folder

 Param:
    - path: npm root path
    - zip: download zip file name
*/
func MkNPM(zip string) {
	npm.New(zip)
	fmt.Println(npm)

	// verify node_modules exist
	if !isDirExist(npm.modules) {
		if err := os.Mkdir(npm.modules, 0755); err != nil {
			P(ERROR, "create %v foler error, Error: %v\n", npm.modules, err.Error())
			return
		} else {
			P(NOTICE, "%v folder create success.\n", npm.modules)
		}
	}

	// clean all npm files
	npm.CleanAll()

	// unzip
	if code, err := unzip(npm.zippath, npm.modules); err != nil {
		fmt.Println(code)
		fmt.Println(err)
	} else {
		if err := os.Rename(npm.modules+util.DIVIDE+code, npm.modules+util.DIVIDE+util.NPM); err != nil {
			P(ERROR, "unzip fail, Error: %v", err.Error())
			return
		} else {
			// copy <root>\node_modules\npm\bin npm and npm.cmd to <root>\
			if err := copyFile(npm.npmbin, npm.root, NPMCOMMAND1); err != nil {
				P(ERROR, "copy %v to %v faild, Error: %v \n", npm.npmbin, npm.root)
				return
			}
			if err := copyFile(npm.npmbin, npm.root, NPMCOMMAND2); err != nil {
				P(ERROR, "copy %v to %v faild, Error: %v \n", npm.npmbin, npm.root)
				return
			}
			// remove download zip file
			if err := os.RemoveAll(npm.zipname); err != nil {
				P(ERROR, "remove %v folder Error: %v.\n", npm.zipname, err.Error())
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

	extractAndWriteFile := func(file *zip.File, idx int, root string) (string, error) {
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
		return root, nil
	}

	idx, root, err := 0, "", nil
	for _, file := range unzip.File {
		root, err = extractAndWriteFile(file, idx, root)
		idx++
	}
	return root, err
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
