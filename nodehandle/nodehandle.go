package nodehandle

import (

	// lib
	. "github.com/Kenshin/cprint"
	"github.com/Kenshin/curl"

	// go
	//"log"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
	"strings"

	// local
	"gnvm/config"
	"gnvm/util"
)

const (
	TIMEFORMART   = "02-Jan-2006 15:04"
	GNVMHOST      = "http://k-zone.cn/gnvm/version.txt"
	PROCESSTAKEUP = "The process cannot access the file because it is being used by another process."
)

var rootPath, latURL string

func init() {
	rootPath = util.GlobalNodePath + util.DIVIDE
	latURL = config.GetConfig("registry") + config.LATEST + "/" + util.SHASUMS
}

/**
 * rootPath    : node.exe global path,         e.g. x:\xxx\xx\xx\
 *
 * global      : global node.exe version num,  e.g. x.xx.xx-x86 ( only rumtime.GOARCH == "amd64", suffix include: 'x86' and 'x64' )
 * globalPath  : global node.exe version path, e.g. x:\xxx\xx\xx\x.xx.xx-x86
 *
 * newer       : newer node.exe version num,   e.g. x.xx.xx
 * newerPath   : newer node.exe version path,  e.g. <rootPath>\x.xx.xx\
 *
 */
func Use(newer string) bool {

	// try catch
	defer func() {
		if err := recover(); err != nil {
			msg := fmt.Sprintf("'gnvm use %v' an error has occurred. please check. \nError: ", newer)
			Error(ERROR, msg, err)
			os.Exit(0)
		}
	}()

	// get true folder, e.g. folder is latest return x.xx.xx
	util.FormatLatVer(&newer, config.GetConfig(config.LATEST_VERSION), true)
	if newer == util.UNKNOWN {
		P(WARING, "current latest is %v, please usage '%v' first. See '%v'.\n", newer, "gnvm update latest", "gnvm help update")
		return false
	}

	// set newerPath and verify newerPath is exist?
	newerPath := rootPath + newer
	if _, err := util.GetNodeVer(newerPath); err != nil {
		P(WARING, "%v folder is not exist %v, use '%v' get local node.exe list. See '%v'.\n", newer, "node.exe", "gnvm ls", "gnvm help ls")
		return false
	}

	// get <root>/node.exe version, when exist, get full version, e.g. x.xx.xx-x86
	global, err := util.GetNodeVer(rootPath)
	if err != nil {
		P(WARING, "not found global node.exe version.\n")
	} else {
		if bit, err := util.Arch(rootPath); err == nil {
			if bit == "x86" && runtime.GOARCH == "amd64" {
				global += "-" + bit
			}
		}
	}

	// check newer is global
	if newer == global {
		P(WARING, "current node.exe version is %v, not re-use. See '%v'.\n", newer, "gnvm node-version")
		return false
	}

	// set globalPath
	globalPath := rootPath + global

	// <root>/global is exist? when not exist, create global folder
	if !util.IsDirExist(globalPath) {
		if err := os.Mkdir(globalPath, 0777); err != nil {
			P(ERROR, "create %v folder Error: %v.\n", global, err.Error())
			return false
		}
	}

	// backup copy <root>/node.exe to <root>/global/node.exe
	if global != "" {
		if err := util.Copy(rootPath, globalPath, util.NODE); err != nil {
			P(ERROR, "copy %v to %v folder Error: %v.\n", rootPath, globalPath, err.Error())
			return false
		}
	}

	// copy <root>/newer/node.exe to <root>/node.exe
	if err := util.Copy(newerPath, rootPath, util.NODE); err != nil {
		P(ERROR, "copy %v to %v folder Error: %v.\n", newerPath, rootPath, err.Error())
		return false
	}

	P(DEFAULT, "Set success, global node.exe version is %v.\n", newer)

	return true
}

/*
 Install node

 Param:
 	- args  : install node.exe versions, include: x.xx.xx latest x.xx.xx-io-x86 x.xx.xx-x86
 	- global: when global == true, call Use func.

 Return:
 	- code: dl[0].Code, usage 'gnvm update latest'

*/
func InstallNode(args []string, global bool) int {

	localVersion, isLatest, code, dl, ts := "", false, 0, new(curl.Download), new(curl.Task)

	// try catch
	defer func() {
		if err := recover(); err != nil {
			if strings.HasPrefix(err.(string), "CURL Error:") {
				fmt.Printf("\n")
			}
			msg := fmt.Sprintf("'gnvm install %v' an error has occurred. \nError: ", strings.Join(args, " "))
			Error(ERROR, msg, err)
			os.Exit(0)
		}
	}()

	for _, v := range args {
		ver, io, arch, suffix, err := util.ParseNodeVer(v)
		if err != nil {
			switch err.Error() {
			case "1":
				P(ERROR, "%v not node.exe download.\n", v)
			case "2":
				P(ERROR, "%v format error, must be '%v' or '%v'.\n", v, "x86", "x64")
			case "3":
				P(ERROR, "%v format error, parameter must be '%v' or '%v'.\n", v, "x.xx.xx", "x.xx.xx-x86|x64")
			case "4":
				P(ERROR, "%v format error, the correct format is %v or %v. \n", v, "0.xx.xx", "^0.xx.xx")
			case "5":
				P(WARING, "'%v' command is no longer supported. See '%v'.\n", "gnvm install npm", "gnvm help npm")
			}
			continue
		}

		// check latest and get remote latest
		v = util.EqualAbs("latest", v)
		if ver == config.LATEST {
			localVersion = config.GetConfig(config.LATEST_VERSION)
			P(NOTICE, "local  latest version is %v.\n", localVersion)

			version := util.GetLatVer(latURL)
			if version == "" {
				P(ERROR, "get latest version error, please check. See '%v'.\n", "gnvm config help")
				break
			}

			isLatest = true
			ver = version
			P(NOTICE, "remote latest version is %v.\n", version)
		} else {
			isLatest = false
		}

		// get folder
		folder := rootPath + ver
		if suffix != "" {
			folder += "-" + suffix
		}
		if _, err := util.GetNodeVer(folder); err == nil {
			P(WARING, "%v folder exist.\n", ver)
			continue
		}

		// get and set url( include iojs)
		url := config.GetConfig(config.REGISTRY)
		if io {
			url = config.GetIOURL(url)
		}

		// add task
		if url, err := util.GetRemoteNodePath(url, ver, arch); err == nil {
			dl.AddTask(ts.New(url, ver, util.NODE, folder))
		}
	}

	// downlaod
	if len(*dl) > 0 {
		curl.Options.Header = false
		curl.Options.Footer = false
		arr := (*dl).GetValues("Title")
		P(DEFAULT, "Start download [%v].\n", strings.Join(arr, ", "))
		newDL, errs := curl.New(*dl)
		for _, task := range newDL {
			v := strings.Replace(task.Dst, rootPath, "", -1)
			if v != localVersion && isLatest {
				config.SetConfig(config.LATEST_VERSION, v)
				P(DEFAULT, "Set success, %v new value is %v\n", config.LATEST_VERSION, v)
			}
			if global && len(args) == 1 {
				if ok := Use(v); ok {
					config.SetConfig(config.GLOBAL_VERSION, v)
				}
			}
		}
		if len(errs) > 0 {
			code = (*dl)[0].Code
			s := ""
			for _, v := range errs {
				s += v.Error()
			}
			P(WARING, s)
		}
		P(DEFAULT, "End download.")
	}

	return code
}

/*
 Uninstall node and npm

 Param:
 	- folder: version

*/
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
		P(ERROR, "unassigned node.exe latest version. See '%v'.\n", "gnvm config INIT")
		return
	}

	// rootPath/version is exist
	if util.IsDirExist(removePath) != true {
		P(ERROR, "%v folder is not exist. See '%v'.\n", folder, "gnvm ls")
		return
	}

	// remove rootPath/version folder
	if err := os.RemoveAll(removePath); err != nil {
		P(ERROR, "uninstall %v fail, Error: %v.\n", folder, err.Error())
		return
	}

	P(DEFAULT, "Node.exe version %v uninstall success.\n", folder)
}

/*
 Update local latest node.exe verion

 - localVersion, remoteVersion: string  node.exe version
 - local, remote:               float64 node.exe version

 Param:
 	- global: when global == true, call Use func.

*/
func Update(global bool) {

	// try catch
	defer func() {
		if err := recover(); err != nil {
			Error(ERROR, "'gnvm updte latest' an error has occurred. \nError: ", err)
			os.Exit(0)
		}
	}()

	localVersion, remoteVersion := config.GetConfig(config.LATEST_VERSION), util.GetLatVer(latURL)

	P(NOTICE, "local latest version is %v.\n", localVersion)
	if remoteVersion == "" {
		P(ERROR, "get latest version error, please check. See '%v'.\n", "gnvm help config")
		return
	}
	P(NOTICE, "remote %v latest version is %v.\n", config.GetConfig("registry"), remoteVersion)

	local, remote, args := util.FormatNodeVer(localVersion), util.FormatNodeVer(remoteVersion), []string{remoteVersion}

	switch {
	case localVersion == config.UNKNOWN:
		if code := InstallNode(args, global); code == 0 {
			config.SetConfig(config.LATEST_VERSION, remoteVersion)
			P(DEFAULT, "Update latest success, current latest version is %v.\n", remoteVersion)
		}
	case local == remote:
		if util.IsDirExist(rootPath + localVersion) {
			cp := CP{Red, false, None, false, "="}
			P(DEFAULT, "Remote latest version %v %v latest version %v, don't need to upgrade.\n", remoteVersion, cp, localVersion)
			if global {
				if ok := Use(localVersion); ok {
					config.SetConfig(config.GLOBAL_VERSION, localVersion)
				}
			}
		} else {
			P(WARING, "local not exist %v\n", localVersion)
			if code := InstallNode(args, global); code == 0 {
				P(DEFAULT, "Download latest version %v success.\n", localVersion)
			}
		}
	case local > remote:
		cp := CP{Red, false, None, false, ">"}
		P(WARING, "local latest version %v %v remote latest version %v.\nPlease check your registry. See 'gnvm help config'.\n", localVersion, cp, remoteVersion)
	case local < remote:
		cp := CP{Red, false, None, false, ">"}
		P(WARING, "remote latest version %v %v local latest version %v.\n", remoteVersion, cp, localVersion)
		if code := InstallNode(args, global); code == 0 {
			config.SetConfig(config.LATEST_VERSION, remoteVersion)
			P(DEFAULT, "Update latest success, current latest version is %v.\n", remoteVersion)
		}
	}
}

/*
 Search node.exe version and Print

 Param:
 	- s: node.exe version, inlcude: *.*.* 0.*.* 0.10.* \<regexp>\ latest 0.10.10

*/
func Search(s string) {
	regex, err := util.FormatWildcard(s, latURL)
	if err != nil {
		P(ERROR, "[%v] %v\n", s, err.Error())
		return
	}

	// set url
	url := config.GetConfig(config.REGISTRY)
	if arr := strings.Split(s, "."); len(arr) == 3 {
		if ver, _ := strconv.Atoi(arr[0]); ver >= 1 && ver <= 3 {
			url = config.GetIOURL(url)
		}
	}
	url += config.NODELIST

	// try catch
	defer func() {
		if err := recover(); err != nil {
			msg := fmt.Sprintf("'gnvm search' an error has occurred. please check %v. \nError: ", url)
			Error(ERROR, msg, err)
			os.Exit(0)
		}
	}()

	// print
	P(DEFAULT, "Search node.exe version rules [%v] from %v, please wait.\n", s, url)

	// generate nodist
	nodist, err, code := New(url, regex)
	if err != nil {
		if code == -1 {
			P(ERROR, "'%v' get url %v error, Error: %v\n", "gnvm search", url, err)
		} else {
			P(ERROR, "%v an error has occurred. please check. Error: %v\n", "gnvm search", err)
		}
		return
	}

	if len(nodist.nl) > 0 {
		nodist.Detail(0)
	} else {
		P(WARING, "not search any node.exe version details, use rules [%v] from %v.\n", s, url)
	}
}

/*
 Print current local node.exe list

 Param:
 	- isPrint: when isPrint == true, print console

*/
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
	files, err := ioutil.ReadDir(rootPath)

	// show error
	if err != nil {
		P(ERROR, "'%v' Error: %v.\n", "gnvm ls", err.Error())
		return lsArr, err
	}

	P(NOTICE, "gnvm.exe root is %v \n", rootPath)
	for _, file := range files {
		// set version
		version := file.Name()

		// check node version
		if util.VerifyNodeVer(version) {

			// <root>/x.xx.xx/node.exe is exist
			if util.IsDirExist(rootPath + version + util.DIVIDE + util.NODE) {
				desc := ""
				switch {
				case version == config.GetConfig(config.GLOBAL_VERSION) && version == config.GetConfig(config.LATEST_VERSION):
					desc = " -- global, latest"
				case version == config.GetConfig(config.LATEST_VERSION):
					desc = " -- latest"
				case version == config.GetConfig(config.GLOBAL_VERSION):
					desc = " -- global"
				}

				ver, _, _, suffix, _ := util.ParseNodeVer(version)
				if suffix == "x86" {
					desc = " -- x86"
				} else if suffix == "x64" {
					desc = " -- x64"
				}

				// set true
				existVersion = true

				// set lsArr
				lsArr = append(lsArr, ver)

				if isPrint {
					if desc == "" {
						P(DEFAULT, "v"+ver+desc, "\n")
					} else {
						P(DEFAULT, "%v", "v"+ver+desc, "\n")
					}

				}
			}
		}
	}

	// version is exist
	if !existVersion {
		P(WARING, "don't have any available Node.js version, please check your input. See '%v'.\n", "gnvm help install")
	}

	return lsArr, err
}

/*
 Print remote node.exe list

 Param:
 	- limit: print max line
 	- io:    when io == true, print iojs

*/
func LsRemote(limit int, io bool) {
	// set url
	url := config.GetConfig(config.REGISTRY)
	if io {
		url = config.GetIOURL(url)
	}
	url += config.NODELIST

	// try catch
	defer func() {
		if err := recover(); err != nil {
			msg := fmt.Sprintf("'gnvm ls --remote' an error has occurred. please check registry %v. \nError: ", url)
			Error(ERROR, msg, err)
			os.Exit(0)
		}
	}()

	// print
	P(DEFAULT, "Read all node.exe version list from %v, please wait.\n", url)

	// generate nodist
	nodist, err, code := New(url, nil)
	if err != nil {
		if code == -1 {
			P(ERROR, "'%v' get url %v error, Error: %v\n", "gnvm search", url, err)
		} else {
			P(ERROR, "%v an error has occurred. please check. Error: %v\n", "gnvm search", err)
		}
		return
	}

	if limit != -1 {
		nodist.Detail(limit)
	} else {
		for _, v := range nodist.Sorts {
			fmt.Println(v)
		}
	}
}

/*
 Show local / global node.exe version

 Param:
 	- args:   include: latest global
 	- remote: when remote == true, print remote latest version

*/
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
		P(DEFAULT, "Node.exe %v version is %v.\n", "latest", latest)
		P(DEFAULT, "Node.exe %v version is %v.\n", "global", global)

		if latest == config.UNKNOWN {
			P(WARING, "latest version is %v, please use '%v'. See '%v'.\n", config.UNKNOWN, "gnvm node-version latest -r", "gnvm help node-version")
		}

		if global == config.UNKNOWN {
			P(WARING, "global version is %v, please use '%v' or '%v'. See '%v'.\n", config.UNKNOWN, "gnvm install latest -g", "gnvm install x.xx.xx -g", "gnvm help install")
		}
	} else {
		switch {
		case args[0] == "global":
			P(DEFAULT, "Node.exe global version is %v.\n", global)
		case args[0] == "latest" && !remote:
			P(DEFAULT, "Node.exe latest version is %v.\n", latest)
		case args[0] == "latest" && remote:
			remoteVersion := util.GetLatVer(latURL)
			if remoteVersion == "" {
				P(ERROR, "get remote %v latest version error, please check. See '%v'.\n", config.GetConfig("registry")+config.LATEST+"/"+config.NODELIST, "gnvm help config")
				P(DEFAULT, "Node.exe latest version is %v.\n", latest)
				return
			}
			P(DEFAULT, "Node.exe remote %v %v is %v.\n", config.GetConfig("registry"), "latest version", remoteVersion)
			P(DEFAULT, "Node.exe local  latest version is %v.\n", latest)
			if latest == config.UNKNOWN {
				config.SetConfig(config.LATEST_VERSION, remoteVersion)
				P(DEFAULT, "Set success, %v new value is %v\n", config.LATEST_VERSION, remoteVersion)
				return
			}
			v1 := util.FormatNodeVer(latest)
			v2 := util.FormatNodeVer(remoteVersion)
			if v1 < v2 {
				cp := CP{Red, false, None, false, ">"}
				P(WARING, "remote latest version %v %v local latest version %v, suggest to upgrade, usage 'gnvm update latest' or 'gnvm update latest -g'.\n", remoteVersion, cp, latest)
			}
		}
	}
}

/*
 Print gnvm.exe version

 Param:
 	- remote: when remote == true, print CHANGELOG

*/
func Version(remote bool) {

	// try catch
	defer func() {
		if err := recover(); err != nil {
			Error(ERROR, "'gnvm version --remote' an error has occurred. \nError: ", err)
			os.Exit(0)
		}
	}()

	localVersion := config.VERSION
	arch := "32 bit"
	if runtime.GOARCH == "amd64" {
		arch = "64 bit"
	}

	cp := CP{Red, true, None, true, "Kenshin Wang"}
	P(DEFAULT, "Current version %v %v.", localVersion, arch, "\n")
	P(DEFAULT, "Copyright (C) 2014-2016 %v <kenshin@ksria.com>", cp, "\n")
	cp.FgColor, cp.Value = Blue, "https://github.com/kenshin/gnvm"
	P(DEFAULT, "See %v for more information.", cp, "\n")

	if !remote {
		return
	}

	code, res, _ := curl.Get(GNVMHOST)
	if code != 0 {
		return
	}
	defer res.Body.Close()

	versionFunc := func(content string, line int) bool {
		if content != "" && line == 1 {
			arr := strings.Fields(content)
			if len(arr) == 2 {

				cp := CP{Red, true, None, true, arr[0][1:]}
				P(DEFAULT, "Latest version %v, publish data %v", cp, arr[1], "\n")

				latestVersion, msg := arr[0][1:], ""
				localArr, latestArr := strings.Split(localVersion, "."), strings.Split(latestVersion, ".")

				switch {
				case latestArr[0] > localArr[0]:
					msg = "must be upgraded."
				case latestArr[1] > localArr[1]:
					msg = "suggest to upgrade."
				case latestArr[2] > localArr[2]:
					msg = "optional upgrade."
				}

				if msg != "" {
					P(NOTICE, msg+" Please download latest %v from %v", "gnvm.exe", "https://github.com/kenshin/gnvm", "\n")
				}
			}

		}
		if line > 1 {
			P(DEFAULT, content)
		}

		return false
	}

	if err := curl.ReadLine(res.Body, versionFunc); err != nil && err != io.EOF {
		P(ERROR, "gnvm version --remote Error: %v\n", err)
	}
}
