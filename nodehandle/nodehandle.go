package nodehandle

import (

	// lib
	"github.com/pierrre/archivefile/zip"

	// go
	//"log"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	// local
	"gnvm/config"
	"gnvm/util"
	"gnvm/util/curl"
	. "gnvm/util/p"
)

const (
	DIVIDE      = "\\"
	NODE        = "node.exe"
	TIMEFORMART = "02-Jan-2006 15:04"
	GNVMHOST    = "http://k-zone.cn/gnvm/version.txt"
)

var rootPath string

func init() {
	rootPath = util.GlobalNodePath + DIVIDE
}

func TransLatestVersion(latest string, isPrint bool) string {
	if latest == config.LATEST {
		latest = config.GetConfig(config.LATEST_VERSION)
		if isPrint {
			P(NOTICE, "current latest version is [%v]", latest)
		}
	}
	return latest
}

/**
 * rootNode is rootPath + "/node.exe", e.g. <root>/node.exe
 *
 * usePath  is use node version path,  e.g. <root>/x.xx.xx
 * useNode  is usePath + "/node.exe",  e.g. <root>/x.xx.xx/node.exe
 *
 * rootVersion is <root>/node.exe version
 * rootFolder  is <root>/rootVersion
 */
func Use(folder string) bool {

	// try catch
	defer func() {
		if err := recover(); err != nil {
			msg := fmt.Sprintf("'gnvm use %v' an error has occurred. please check. \nError: ", folder)
			Error(ERROR, msg, err)
			os.Exit(0)
		}
	}()

	rootNodeExist := true

	// get true folder, e.g. folder is latest return x.xx.xx
	folder = TransLatestVersion(folder, true)

	if folder == config.UNKNOWN {
		P(ERROR, "unassigned Node.js latest version. See 'gnvm install latest'.")
		return false
	}

	// set rootNode
	rootNode := rootPath + NODE

	// set usePath and useNode
	usePath := rootPath + folder + DIVIDE
	useNode := usePath + NODE

	// get <root>/node.exe version
	rootVersion, err := util.GetNodeVersion(rootPath)
	if err != nil {
		P(WARING, "not found global node version, please use 'gnvm install x.xx.xx -g'. See 'gnvm help install'.")
		rootNodeExist = false
	}

	// <root>/folder is exist
	if isDirExist(usePath) != true {
		P(WARING, "[%v] folder is not exist from [%v]. Get local node.exe version. See 'gnvm ls'.", folder, rootPath)
		return false
	}

	// check folder is rootVersion
	if folder == rootVersion {
		P(WARING, "current node.exe version is [%v], not re-use. See 'gnvm node-version'.", folder)
		return false
	}

	// set rootFolder
	rootFolder := rootPath + rootVersion

	// <root>/rootVersion is exist
	if isDirExist(rootFolder) != true {

		// create rootVersion folder
		if err := cmd("md", rootFolder); err != nil {
			P(ERROR, "create %v folder Error: %v.", rootVersion, err.Error())
			return false
		}

	}

	if rootNodeExist {
		// copy rootNode to <root>/rootVersion
		if err := copy(rootNode, rootFolder); err != nil {
			P(ERROR, "copy %v to %v folder Error: %v.\n", rootNode, rootFolder, err.Error())
			return false
		}

		// delete <root>/node.exe
		if err := os.Remove(rootNode); err != nil {
			P(ERROR, "remove %v folder Error: %v.\n", rootNode, err.Error())
			return false
		}

	}

	// copy useNode to rootPath
	if err := copy(useNode, rootPath); err != nil {
		P(ERROR, "copy %v to %v folder Error: %v.\n", useNode, rootPath, err.Error())
		return false
	}

	P(DEFAULT, "Set success, Current Node.exe version is [%v].\n", folder)

	return true
}

func NodeVersion(args []string, remote bool) {

	// try catch
	defer func() {
		if err := recover(); err != nil {
			msg := fmt.Sprintf("'gnvm node-version %v' an error has occurred. please check. \nError: ", strings.Join(args, " "))
			Error(ERROR, msg, err)
			os.Exit(0)
		}
	}()

	latest := config.GetConfig(config.LATEST_VERSION)
	global := config.GetConfig(config.GLOBAL_VERSION)

	if len(args) == 0 || len(args) > 1 {
		P(DEFAULT, "Node.exe latest verson is [%v].\n", latest)
		P(DEFAULT, "Node.exe global verson is [%v].\n", global)
	} else {
		switch {
		case args[0] == "global":
			P(DEFAULT, "Node.exe global verson is [%v].\n", global)
		case args[0] == "latest" && !remote:
			P(DEFAULT, "Node.exe latest verson is [%v].\n", latest)
		case args[0] == "latest" && remote:
			remoteVersion := getLatestVersionByRemote()
			if remoteVersion == "" {
				P(ERROR, "get remote [%v] latest version error, please check. See 'gnvm config help'.\n", config.GetConfig("registry")+config.LATEST+"/"+config.NODELIST)
				P(DEFAULT, "Node.exe latest verson is [%v].\n", latest)
				return
			}
			P(DEFAULT, "Node.exe remote [%v] verson is [%v].\n", config.GetConfig("registry"), remoteVersion)
			P(DEFAULT, "Node.exe latest verson is [%v].\n", latest)
		}
	}

	switch {
	case len(args) == 0 && (global == config.UNKNOWN || latest == config.UNKNOWN):
		P(WARING, "when version is [%v], please Use 'gnvm config INIT'. See 'gnvm help config'.\n", config.UNKNOWN)
	case len(args) > 0 && args[0] == "latest" && latest == config.UNKNOWN:
		P(WARING, "when version is [%v], please Use 'gnvm config INIT'. See 'gnvm help config'.\n", config.UNKNOWN)
	case len(args) > 0 && args[0] == "global" && global == config.UNKNOWN:
		P(WARING, "when version is [%v], please Use 'gnvm config INIT'. See 'gnvm help config'.\n", config.UNKNOWN)
	}
}

func Uninstall(folder string) {

	// try catch
	defer func() {
		if err := recover(); err != nil {
			msg := fmt.Sprintf("'gnvm uninstall %v' an error has occurred. please check. \nError: ", folder)
			Error(ERROR, msg, err)
			os.Exit(0)
		}
	}()

	// set removePath
	removePath := rootPath + folder

	if folder == config.UNKNOWN {
		P(ERROR, "unassigned Node.js latest version. See 'gnvm install latest'.")
		return
	}

	// rootPath/version is exist
	if isDirExist(removePath) != true {
		P(ERROR, "[%v] folder is not exist. Get local node.exe version. See 'gnvm ls'.\n", folder)
		return
	}

	// remove rootPath/version folder
	if err := os.RemoveAll(removePath); err != nil {
		P(ERROR, "uinstall [%v] fail, Error: %v.\n", folder, err.Error())
		return
	}

	P(DEFAULT, "Node.exe version [%v] uninstall success.\n", folder)
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

	if !isDirExist(rootPath + "npm.cmd") {
		P(ERROR, "[%v] not exist npm.\n", rootPath)
		return
	}

	// remove npm.cmd
	if err := os.RemoveAll(rootPath + "npm.cmd"); err != nil {
		removeFlag = false
		P(ERROR, "remove [%v] file fail from [%v], Error: %v.\n", "npm.cmd", rootPath, err.Error())
	}

	// remove node_modules/npm
	if err := os.RemoveAll(rootPath + "node_modules" + DIVIDE + "npm"); err != nil {
		removeFlag = false
		P(ERROR, "remove [%v] folder fail from [%v], Error: %v.\n", "npm", rootPath+"node_modules", err.Error())
	}

	if removeFlag {
		P(DEFAULT, "npm uninstall success from [%v].\n", rootPath)
	}
}

func LS(isPrint bool) ([]string, error) {

	// try catch
	defer func() {
		if err := recover(); err != nil {
			Error(ERROR, "'gnvm ls' an error has occurred. please check. \nError: ", err)
			os.Exit(0)
		}
	}()

	var lsArr []string
	existVersion := false
	err := filepath.Walk(rootPath, func(dir string, f os.FileInfo, err error) error {

		// check nil
		if f == nil {
			return err
		}

		// check dir
		if f.IsDir() == false {
			return nil
		}

		// set version
		version := f.Name()

		// check node version
		if ok := util.VerifyNodeVersion(version); ok {

			// <root>/x.xx.xx/node.exe is exist
			if isDirExist(rootPath + version + DIVIDE + NODE) {
				desc := ""
				switch {
				case version == config.GetConfig(config.GLOBAL_VERSION) && version == config.GetConfig(config.LATEST_VERSION):
					desc = " -- global, latest"
				case version == config.GetConfig(config.LATEST_VERSION):
					desc = " -- latest"
				case version == config.GetConfig(config.GLOBAL_VERSION):
					desc = " -- global"
				}

				// set true
				existVersion = true

				// set lsArr
				lsArr = append(lsArr, version)

				if isPrint {
					if desc == "" {
						P(DEFAULT, "v"+version+desc)
					} else {
						P(DEFAULT, "%v", "v"+version+desc)
					}

				}
			}
		}

		// return
		return nil
	})

	// show error
	if err != nil {
		P(ERROR, "'gnvm ls' Error: %v.\n", err.Error())
		return lsArr, err
	}

	// version is exist
	if !existVersion {
		P(WARING, "don't have any available version, please check. See 'gnvm help install'.")
	}

	return lsArr, err
}

func LsRemote() {

	// set url
	url := config.GetConfig(config.REGISTRY) + config.NODELIST

	// try catch
	defer func() {
		if err := recover(); err != nil {
			msg := fmt.Sprintf("'gnvm ls --remote' an error has occurred. please check registry: [%v]. \nError: ", url)
			Error(ERROR, msg, err)
			os.Exit(0)
		}
	}()

	// set exist version
	isExistVersion := false

	// print
	P(DEFAULT, "Read all Node.exe version list from [%v], please wait.", url)

	// get
	code, res, _ := curl.Get(url)
	if code != 0 {
		return
	}
	// close
	defer res.Body.Close()

	writeVersion := func(content string, line int) {
		// replace '\n'
		content = strings.Replace(content, "\n", "", -1)

		// split 'vx.xx.xx  1.1.0-alpha-2'
		args := strings.Split(content, " ")

		if ok := util.VerifyNodeVersion(args[0][1:]); ok {
			isExistVersion = true
			P(DEFAULT, args[0])
		}
	}

	if err := curl.ReadLine(res.Body, writeVersion); err != nil && err != io.EOF {
		P(ERROR, "gnvm ls --remote Error: %v", err)
	}

}

/*
 * return code same as download return code
 */
func Install(args []string, global bool) int {

	var currentLatest string
	var code int

	// try catch
	defer func() {
		if err := recover(); err != nil {
			msg := fmt.Sprintf("'gnvm install %v' an error has occurred. \nError: ", strings.Join(args, " "))
			Error(ERROR, msg, err)
			os.Exit(0)
		}
	}()

	for _, v := range args {

		if v == config.LATEST {

			version := getLatestVersionByRemote()
			if version == "" {
				P(ERROR, "get latest version error, please check. See 'gnvm config help'.")
				break
			}

			// set v
			v = version
			currentLatest = version
			P(NOTICE, "current latest version is [%v]", version)
		}

		// downlaod
		code = download(v)
		if code == 0 || code == 2 {

			if v == currentLatest {
				config.SetConfig(config.LATEST_VERSION, v)
			}

			if global && len(args) == 1 {
				if ok := Use(v); ok {
					config.SetConfig(config.GLOBAL_VERSION, v)
				}
			}
		}
	}

	return code

}

func InstallNpm() {

	// try catch
	defer func() {
		if err := recover(); err != nil {
			Error(ERROR, "'gnvm install npm' an error has occurred. \nError: ", err)
			os.Exit(0)
		}
	}()

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

	getNpmVersion := func(content string, line int) {
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
	}

	if err := curl.ReadLine(res.Body, getNpmVersion); err != nil && err != io.EOF {
		P(ERROR, "parse npm version Error: %v, from %v", err, url)
		return
	}

	if maxVersion == "" {
		P(ERROR, "get npm version fail from [%v], please check. See 'gnvm help config'.\n", url)
		return
	}

	P(NOTICE, "the latest version is [%v] from [%v].\n", maxVersion, config.GetConfig(config.REGISTRY))

	// download zip
	zipPath := os.TempDir() + DIVIDE + maxVersion
	if code := downloadNpm(maxVersion); code == 0 {

		P(DEFAULT, "Start unarchive file [%v].\n", maxVersion)

		//unzip(maxVersion)

		if err := zip.UnarchiveFile(zipPath, config.GetConfig(config.NODEROOT), nil); err != nil {
			panic(err)
		}

		P(DEFAULT, "End unarchive.")
	}

	// remove temp zip file
	if err := os.RemoveAll(zipPath); err != nil {
		P(ERROR, "remove zip file fail from [%v], Error: %v.\n", zipPath, err.Error())
	}

}

func Update(global bool) {

	// try catch
	defer func() {
		if err := recover(); err != nil {
			Error(ERROR, "'gnvm updte latest' an error has occurred. \nError: ", err)
			os.Exit(0)
		}
	}()

	localVersion := config.GetConfig(config.LATEST_VERSION)
	P(NOTICE, "local latest version is [%v].\n", localVersion)

	remoteVersion := getLatestVersionByRemote()
	if remoteVersion == "" {
		P(ERROR, "get latest version error, please check. See 'gnvm config help'.")
		return
	}
	P(NOTICE, "remote [%v] latest version is [%v].\n", config.GetConfig("registry"), remoteVersion)

	local, _ := util.ConverFloat(localVersion)
	remote, _ := util.ConverFloat(remoteVersion)

	var args []string
	args = append(args, remoteVersion)

	switch {
	case localVersion == config.UNKNOWN:
		P(WARING, "local latest version undefined.")
		if code := Install(args, global); code == 0 || code == 2 {
			config.SetConfig(config.LATEST_VERSION, remoteVersion)
			P(DEFAULT, "Update latest success, current latest version is [%v].\n", remoteVersion)
		}
	case local == remote:

		if isDirExist(rootPath + localVersion) {
			cc := CC{Red, false, None, false, "="}
			P(DEFAULT, "Remote latest version [%v] %v latest version [%v], don't need to upgrade.\n", remoteVersion, cc, localVersion)
			if global {
				if ok := Use(localVersion); ok {
					config.SetConfig(config.GLOBAL_VERSION, localVersion)
				}
			}
		} else if !isDirExist(rootPath + localVersion) {
			P(WARING, "local not exist %v", localVersion)
			if code := Install(args, global); code == 0 || code == 2 {
				P(DEFAULT, "Download latest version [%v] success.\n", localVersion)
			}
		}

	case local > remote:
		cc := CC{Red, false, None, false, ">"}
		P(WARING, "local latest version [%v] %v remote latest version [%v].\nPlease check your registry. See 'gnvm help config'.\n", localVersion, cc, remoteVersion)
	case local < remote:
		cc := CC{Red, false, None, false, ">"}
		P(WARING, "remote latest version [%v] %v local latest version [%v].\n", remoteVersion, cc, localVersion)
		if code := Install(args, global); code == 0 || code == 2 {
			config.SetConfig(config.LATEST_VERSION, remoteVersion)
			P(DEFAULT, "Update latest success, current latest version is [%v].\n", remoteVersion)
		}
	}
}

func Version(remote bool) {

	// try catch
	defer func() {
		if err := recover(); err != nil {
			Error(ERROR, "'gnvm version --remote' an error has occurred. \nError: ", err)
			os.Exit(0)
		}
	}()

	localVersion := config.VERSION

	P(DEFAULT, "Current version %v", localVersion)

	if !remote {
		return
	}

	code, res, _ := curl.Get(GNVMHOST)
	if code != 0 {
		return
	}
	defer res.Body.Close()

	versionFunc := func(content string, line int) {
		if content != "" && line == 1 {
			arr := strings.Fields(content)
			if len(arr) == 2 {

				P(DEFAULT, "Latest version %v, publish data %v", arr[0][1:], arr[1])

				latestVersion, msg := arr[0][1:], ""
				localArr, latestArr := strings.Split(localVersion, "."), strings.Split(latestVersion, ".")

				switch {
				case latestArr[0] > localArr[0]:
					msg = "must be upgraded"
				case latestArr[1] > localArr[1]:
					msg = "suggest to upgrade"
				case latestArr[2] > localArr[2]:
					msg = "optional upgrade"
				}

				if msg != "" {
					P(NOTICE, msg)
				}

			} else {
				return
			}
		}
		if line > 1 {
			P(DEFAULT, content)
		}
	}

	if err := curl.ReadLine(res.Body, versionFunc); err != nil && err != io.EOF {
		P(ERROR, "gnvm version --remote Error: %v", err)
	}

}

func isDirExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	} else {
		// return file.IsDir()
		return true
	}
}

func cmd(name, arg string) error {
	_, err := exec.Command("cmd", "/C", name, arg).Output()
	return err
}

func copy(src, dest string) error {
	_, err := exec.Command("cmd", "/C", "copy", "/y", src, dest).Output()
	return err
}

/*
 * return code
 * 0: success
 * 1: remove folder error
 * 2: folder exist
 * 3: create folder error
 *
 */
func download(version string) int {

	// get current os arch
	amd64 := "/"
	if runtime.GOARCH == "amd64" {
		amd64 = "/x64/"
	}

	// rootPath/version/node.exe is exist
	if _, err := util.GetNodeVersion(rootPath + version + DIVIDE); err == nil {
		P(WARING, "[%v] folder exist.\n", version)
		return 2
	} else {
		if err := os.RemoveAll(rootPath + version); err != nil {
			P(ERROR, "remove [%v] fail, Error: %v\n", version, err.Error())
			return 1
		}
		//P(DEFAULT, "Remove empty [%v] folder success.\n", version)
	}

	// rootPath/version is exist
	if isDirExist(rootPath+version) != true {
		if err := os.Mkdir(rootPath+version, 0777); err != nil {
			P(ERROR, "create [%v] fail, Error: %v\n", version, err.Error())
			return 3
		}
	}

	// set url
	url := config.GetConfig(config.REGISTRY) + "v" + version + amd64 + NODE

	// download
	if code := curl.New(url, version, rootPath+version+DIVIDE+NODE); code != 0 {
		if code == -1 {
			if err := os.RemoveAll(rootPath + version); err != nil {
				P(ERROR, "remove [%v] fail, Error: %v\n", version, err.Error())
				return 1
			}
		}
		return code
	}

	return 0
}

/*
 * return code
 * 0: success
 *
 */
func downloadNpm(version string) int {

	// set url
	url := config.GetConfig(config.REGISTRY) + "npm/" + version

	// download
	if code := curl.New(url, version, os.TempDir()+DIVIDE+version); code != 0 {
		return code
	}

	return 0
}

func getLatestVersionByRemote() string {

	var version string

	// set url
	url := config.GetConfig("registry") + "latest/" + util.SHASUMS

	version = util.GetLatestVersion(url)

	return version

}
