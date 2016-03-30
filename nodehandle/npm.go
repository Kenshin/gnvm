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
	LATNPMURL  = "https://raw.githubusercontent.com/npm/npm/master/package.json"
	NPMTAOBAO  = "http://npm.taobao.org/mirrors/npm/"
	NPMDEFAULT = "https://github.com/npm/npm/releases/"
	ZIP        = ".zip"
)

/*
- root:     config.GetConfig(config.NODEROOT)
- zipname:  v3.8.5.zip
- ziproot:  <v3.8.5.zip>/<root_folder>
- zippath:  /<root>/v3.8.5.zip
- modules:  /<root>/node_modules
- npmpath:  /<root>/node_modules/npm
- npmbin:   /<root>/node_modules/npm/bin
- command1: npm
- command2: npm.cmd
*/
type NPMange struct {
	root     string
	zipname  string
	ziproot  string
	zippath  string
	modules  string
	npmpath  string
	npmbin   string
	command1 string
	command2 string
}

var npm = new(NPMange)

/*
 Crete NPMange
*/
func (this *NPMange) New(zip string) {
	(*this).root = config.GetConfig(config.NODEROOT)
	(*this).modules = (*this).root + util.DIVIDE + "node_modules"
	(*this).zipname = zip
	(*this).zippath = (*this).root + util.DIVIDE + (*this).zipname
	(*this).npmpath = (*this).modules + util.DIVIDE + util.NPM
	(*this).npmbin = (*this).npmpath + util.DIVIDE + "bin"
	(*this).command1 = "npm"
	(*this).command2 = "npm.cmd"
}

/*
 Custom Print
*/
func (this *NPMange) String() string {
	s := fmt.Sprintf("root     = %v\n", this.root)
	s += fmt.Sprintf("zipname  = %v\n", this.zipname)
	s += fmt.Sprintf("ziproot  = %v\n", this.ziproot)
	s += fmt.Sprintf("zippath  = %v\n", this.zippath)
	s += fmt.Sprintf("modules  = %v\n", this.modules)
	s += fmt.Sprintf("npmpath  = %v\n", this.npmpath)
	s += fmt.Sprintf("npmbin   = %v\n", this.npmbin)
	s += fmt.Sprintf("command1 = %v\n", this.command1)
	s += fmt.Sprintf("command2 = %v", this.command2)
	return s
}

/*
 Create node_modules folder
*/
func (this *NPMange) CreateModules() {
	if !isDirExist(this.modules) {
		if err := os.Mkdir(this.modules, 0755); err != nil {
			P(ERROR, "create %v foler error, Error: %v\n", this.modules, err.Error())
		} else {
			P(NOTICE, "%v folder create success.\n", this.modules)
		}
	}
}

/*
 Download npm zip

 Param:
    - url: download url

 Return:
    - error
*/
func (this *NPMange) Download(url, name string) error {
	if _, errs := curl.New(url, name, name, (*this).root); len(errs) > 0 {
		err := errs[0]
		P(ERROR, "%v an error has occurred, url %v, Error is %v. See '%v'.\n", "gnvm npm", url, err, "gnvm help npm")
		return err
	}
	return nil
}

/*
  Unzip file

  Return:
    - error
    - code
        - -1: open  zip file error
        - -2: open  file error
        - -3: write file error
        - -4: copy  file error
*/
func (this *NPMange) Unzip() (int, error) {
	path, dest := (*this).zippath, (*this).modules
	unzip, err := zip.OpenReader(path)
	if err != nil {
		return -1, err
	}
	defer unzip.Close()

	extractAndWriteFile := func(file *zip.File, idx int) (int, error) {
		rc, err := file.Open()
		if err != nil {
			return -2, err
		}
		defer rc.Close()
		if idx == 0 {
			(*this).ziproot = strings.Replace(file.Name, "/", "", -1)
		}
		path = filepath.Join(dest, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
		} else {
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
			if err != nil {
				return -3, err
			}
			defer f.Close()
			if _, err := io.Copy(f, rc); err != nil {
				return -4, err
			}
		}
		return 0, nil
	}

	idx := 0
	for _, file := range unzip.File {
		if code, err := extractAndWriteFile(file, idx); err != nil {
			return code, err
		}
		idx++
	}
	return 0, nil
}

/*
 Rename <root>\node_modules\folder to <root>\node_modules\npm
 Copy <root>\node_modules\npm\bin\ npm and npm.cmd to <root>\
*/
func (this *NPMange) Install() error {
	if err := os.Rename(this.modules+util.DIVIDE+this.ziproot, this.npmpath); err != nil {
		P(ERROR, "rename fail, Error: %v\n", err.Error())
		return err
	} else {
		files := [2]string{this.command1, this.command2}
		for _, v := range files {
			if err := util.Copy(this.npmbin, this.root, v); err != nil {
				P(ERROR, "copy %v to %v faild, Error: %v \n", this.npmbin, this.root)
				return err
			}
		}
	}
	return nil
}

/*
 Remove file, inlcude:
    - <root>/node_modules/npm
    - <root>/npm
    - <root>/npm.cmd
    - <root>/<npm.zip>
*/
func (this *NPMange) Clean(path string) error {
	if isDirExist(path) {
		if err := os.RemoveAll(path); err != nil {
			P(ERROR, "remove %v folder Error: %v.\n", path, err.Error())
			return err
		}
	}
	return nil
}

/*
 Remove <root>/node_modules/npm, <root>/npm, <root>/npm.cmd
*/
func (this *NPMange) CleanAll() error {
	paths := [3]string{this.npmpath, this.root + util.DIVIDE + this.command1, this.root + util.DIVIDE + this.command2}
	for _, v := range paths {
		if err := this.Clean(v); err != nil {
			return err
		}
	}
	return nil
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
	if version != util.LATEST && version != util.GLOBAL && !util.VerifyNodeVer(version) {
		P(ERROR, "'%v' param only support [%v] [%v] [%v], please check your input. See '%v'.\n", "gnvm npm", "latest", "global", "valid version", "gnvm help npm")
		return
	}

	prompt, local, newver := "n", getGlobalNPMVer(), version

	if version == util.LATEST {
		newver = getLatNPMVer()
	}
	cp := CP{Red, false, None, false, newver}

	P(NOTICE, "local    npm version is %v\n", local)
	P(NOTICE, "download npm version is %v\n", cp)
	P(NOTICE, "download %v version [Y/n]? ", cp)
	fmt.Scanf("%s\n", &prompt)
	prompt = strings.ToLower(prompt)
	if prompt == "y" {
		downloadNpm(newver)
	} else {
		P(NOTICE, "operation has been cancelled.")
	}

	//v1, v2 := util.FormatNodeVer(newver), util.FormatNodeVer(local)

	/*
		cp := CP{Red, false, None, false, "="}
		switch {
		case v1 > v2:
			cp.Value = ">"
			P(WARING, "npm newver latest version %v %v local latest version %v.\n", newver, cp, local)
			P(NOTICE, "is update local npm version [Y/n]? ")
			fmt.Scanf("%s\n", &prompt)
			prompt = strings.ToLower(prompt)
			if prompt == "y" {
				downloadNpm(newver)
			} else {
				P(NOTICE, "you need use '%v' update local version. \n", "npm install -g npm")
			}
		case v1 < v2:
			cp.Value = "<"
			P(WARING, "npm newver latest version %v %v local latest version %v.\n", newver, cp, local)
		case v1 == v2:
			P(WARING, "npm newver latest version %v %v local latest version %v.\n", newver, cp, local)
		}
	*/
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

	// create npm
	npm.New(version)

	// download
	if err := npm.Download(url, version); err != nil {
		return
	}

	// create node_modules
	npm.CreateModules()

	// clean all npm files
	npm.CleanAll()

	// unzip
	if _, err := npm.Unzip(); err != nil {
		P(ERROR, "unzip %v an error has occurred. \nError: ", npm.zipname, err.Error())
		return
	}

	// install
	if err := npm.Install(); err == nil {
		npm.Clean(npm.zippath)
		P(NOTICE, "unzip complete.\n")
	}

	fmt.Println(npm)
}
