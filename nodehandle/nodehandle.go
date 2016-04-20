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

var rootPath, latURL string

func init() {
	rootPath = util.GlobalNodePath + util.DIVIDE
	latURL = config.GetConfig(config.REGISTRY) + util.LATEST + "/" + util.SHASUMS
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
		P(WARING, "%v folder is not exist %v, use '%v' get local Node.js version list. See '%v'.\n", newer, "node.exe", "gnvm ls", "gnvm help ls")
		return false
	}

	// get <root>/node.exe version, when exist, get full version, e.g. x.xx.xx-x86
	global, err := util.GetNodeVer(rootPath)
	if err != nil {
		P(WARING, "not found %v Node.js version.\n", "global")
	} else {
		if bit, err := util.Arch(rootPath); err == nil {
			if bit == "x86" && runtime.GOARCH == "amd64" {
				global += "-" + bit
			}
		}
	}

	// check newer is global
	if newer == global {
		P(WARING, "current Node.js version is %v, not re-use. See '%v'.\n", newer, "gnvm node-version")
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

	P(DEFAULT, "Set success, global Node.js version is %v.\n", newer)

	return true
}

/*
 Install node

 Param:
 	- args  : install Node.js versions, include: x.xx.xx latest x.xx.xx-io-x86 x.xx.xx-x86
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
				P(ERROR, "%v format error, suffix only must be '%v' or '%v'.\n", v, "x86", "x64")
			case "3":
				P(ERROR, "%v format error, parameter must be '%v' or '%v'.\n", v, "x.xx.xx", "x.xx.xx-x86|x64")
			case "4":
				P(ERROR, "%v not an %v Node.js version.\n", v, "valid")
			case "5":
				P(WARING, "'%v' command is no longer supported. See '%v'.\n", "gnvm install npm", "gnvm help npm")
			}
			continue
		}

		// when os is 386, not download 64 bit node.exe
		if runtime.GOARCH == "386" && suffix == "x64" {
			P(WARING, "current operating system is %v, not support %v suffix.\n", "32-bit", "x64")
			continue
		}

		// check local latest and get remote latest
		v = util.EqualAbs("latest", v)
		if ver == util.LATEST {
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

		// ture folder name
		if suffix != "" {
			ver += "-" + suffix
		}

		// verify <root>/folder is exist
		folder := rootPath + ver
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
		arr := (*dl).GetValues("Title")
		P(DEFAULT, "Start download Node.js versions [%v].\n", strings.Join(arr, ", "))
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
			msg := fmt.Sprintf("gnvm uninstall %v an error has occurred. please check your input. \nError: ", folder)
			Error(ERROR, msg, err)
			os.Exit(0)
		}
	}()

	// set removePath
	removePath := rootPath + folder

	if folder == util.UNKNOWN {
		P(ERROR, "current latest version is %v, please usage '%v' first. See '%v'.\n", folder, "gnvm update latest", "gnvm help update")
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

	P(DEFAULT, "Node.js version %v uninstall success.\n", folder)
}

/*
 Update local Node.js latest verion

 - localVersion, remoteVersion: string  Node.js version
 - local, remote:               float64 Node.js version

 Param:
 	- global: when global == true, call Use func.

*/
func Update(global bool) {

	// try catch
	defer func() {
		if err := recover(); err != nil {
			msg := fmt.Sprintf("'%v' an error has occurred. \nError: ", "gnvm updte latest")
			Error(ERROR, msg, err)
			os.Exit(0)
		}
	}()

	localVersion, remoteVersion := config.GetConfig(config.LATEST_VERSION), util.GetLatVer(latURL)

	P(NOTICE, "local  Node.js latest version is %v.\n", localVersion)
	if remoteVersion == "" {
		P(ERROR, "get latest version error, please check. See '%v'.\n", "gnvm help config")
		return
	}
	P(NOTICE, "remote Node.js latest version is %v from %v.\n", remoteVersion, config.GetConfig("registry"))

	local, remote, args := util.FormatNodeVer(localVersion), util.FormatNodeVer(remoteVersion), []string{remoteVersion}

	switch {
	case localVersion == util.UNKNOWN:
		if code := InstallNode(args, global); code == 0 {
			config.SetConfig(config.LATEST_VERSION, remoteVersion)
			P(DEFAULT, "Update Node.js latest success, current latest version is %v.\n", remoteVersion)
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
			P(WARING, "%v folder is not exist. See '%v'.\n", localVersion, "gnvm ls")
			if code := InstallNode(args, global); code == 0 {
				P(DEFAULT, "Local Node.js latest version is %v.\n", localVersion)
			}
		}
	case local > remote:
		cp := CP{Red, false, None, false, ">"}
		P(WARING, "local latest version %v %v remote latest version %v.\nPlease check your config %v. See '%v'.\n", localVersion, cp, remoteVersion, "registry", "gnvm help config")
	case local < remote:
		cp := CP{Red, false, None, false, ">"}
		P(WARING, "remote latest version %v %v local latest version %v.\n", remoteVersion, cp, localVersion)
		if code := InstallNode(args, global); code == 0 {
			config.SetConfig(config.LATEST_VERSION, remoteVersion)
			P(DEFAULT, "Update success, Node.js latest version is %v.\n", remoteVersion)
		}
	}
}

/*
 Search Node.js version and Print

 Param:
 	- s: Node.js version, inlcude: *.*.* 0.*.* 0.10.* /<regexp>/ latest 0.10.10

*/
func Search(s string) {
	regex, err := util.FormatWildcard(s, latURL)
	if err != nil {
		P(ERROR, "%v not an %v Node.js version.\n", s, "valid")
		return
	}

	// set url
	url := config.GetConfig(config.REGISTRY)
	if arr := strings.Split(s, "."); len(arr) == 3 {
		if ver, _ := strconv.Atoi(arr[0]); ver >= 1 && ver <= 3 {
			url = config.GetIOURL(url)
		}
	}
	url += util.NODELIST

	// try catch
	defer func() {
		if err := recover(); err != nil {
			msg := fmt.Sprintf("'%v' an error has occurred. please check your input.\nError: ", "gnvm search")
			Error(ERROR, msg, err)
			os.Exit(0)
		}
	}()

	// print
	P(DEFAULT, "Search Node.js version rules [%v] from %v, please wait.\n", s, url)

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
		P(WARING, "not search any Node.js version details, use rules [%v] from %v.\n", s, url)
	}
}

/*
 Print current local Node.js version list

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
				lsArr = append(lsArr, version)

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
 Print remote Node.js version list

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
	url += util.NODELIST

	// try catch
	defer func() {
		if err := recover(); err != nil {
			msg := fmt.Sprintf("'gnvm ls --remote' an error has occurred. please check your input %v. \nError: ", url)
			Error(ERROR, msg, err)
			os.Exit(0)
		}
	}()

	// print
	P(DEFAULT, "Read all Node.js version list from %v, please wait.\n", url)

	// generate nodist
	nodist, err, code := New(url, nil)
	if err != nil {
		if code == -1 {
			P(ERROR, "'%v' get url %v error, Error: %v\n", "gnvm ls -r -d", url, err)
		} else {
			P(ERROR, "%v an error has occurred. please check your input. Error: %v\n", "gnvm ls -r -d", err)
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
 Show local / global Node.js version

 Param:
 	- args:   include: latest global

*/
func NodeVersion(args []string) {

	// try catch
	defer func() {
		if err := recover(); err != nil {
			msg := fmt.Sprintf("'gnvm node-version %v' an error has occurred. please check. \nError: ", strings.Join(args, " "))
			Error(ERROR, msg, err)
			os.Exit(0)
		}
	}()

	/*
		if len(args) == 0 {
			P(DEFAULT, "Node.js %v version is %v.\n", "latest", latest)
			P(DEFAULT, "Node.js %v version is %v.\n", "global", global)
			if latest == util.UNKNOWN {
				P(WARING, "latest version is %v, please use '%v'. See '%v'.\n", util.UNKNOWN, "gnvm node-version latest -r", "gnvm help node-version")
			}
			if global == util.UNKNOWN {
				P(WARING, "global version is %v, please use '%v' or '%v'. See '%v'.\n", util.UNKNOWN, "gnvm install latest -g", "gnvm install x.xx.xx -g", "gnvm help install")
			}
		} else {
			switch {
			case args[0] == "global":
				P(DEFAULT, "Node.js global version is %v.\n", global)
				if global == util.UNKNOWN {
					P(WARING, "global version is %v, please use '%v' or '%v'. See '%v'.\n", util.UNKNOWN, "gnvm install latest -g", "gnvm install x.xx.xx -g", "gnvm help install")
				}
			case args[0] == "latest" && !remote:
				P(DEFAULT, "Node.js latest version is %v.\n", latest)
				if latest == util.UNKNOWN {
					P(WARING, "latest version is %v, please use '%v'. See '%v'.\n", util.UNKNOWN, "gnvm node-version latest -r", "gnvm help node-version")
				}
			case args[0] == "latest" && remote:
				remoteVersion := util.GetLatVer(latURL)
				if remoteVersion == "" {
					P(ERROR, "get remote %v Node.js %v error, please check your input. See '%v'.\n", config.GetConfig(config.REGISTRY), "latest version", "gnvm help config")
					return
				}
				P(DEFAULT, "Local  Node.js latest version is %v.\n", latest)
				P(DEFAULT, "Remote Node.js latest version is %v from %v.\n", remoteVersion, config.GetConfig(config.REGISTRY))
				if latest == util.UNKNOWN {
					config.SetConfig(config.LATEST_VERSION, remoteVersion)
					P(DEFAULT, "Set success, local Node.js %v version is %v.\n", util.LATEST, remoteVersion)
					return
				}
				v1 := util.FormatNodeVer(latest)
				v2 := util.FormatNodeVer(remoteVersion)
				if v1 < v2 {
					cp := CP{Red, false, None, false, ">"}
					P(WARING, "remote Node.js latest version %v %v local Node.js latest version %v, suggest to upgrade, usage '%v'.\n", remoteVersion, cp, latest, "gnvm update latest")
				}
			}
		}
	*/

	isLatest, isGlobal := false, false
	latest, global := config.GetConfig(config.LATEST_VERSION), config.GetConfig(config.GLOBAL_VERSION)
	if len(args) == 0 {
		isLatest = true
		isGlobal = true
	} else {
		if args[0] == util.LATEST {
			isLatest = true
		} else {
			isGlobal = true
		}
	}

	if isGlobal {
		if global == util.UNKNOWN {
			P(WARING, "global Node.js version is %v.\n", util.UNKNOWN)
			if global, err := util.GetNodeVer(rootPath); err == nil {
				config.SetConfig(config.GLOBAL_VERSION, global)
				P(DEFAULT, "Set success, %v new value is %v.\n", config.GLOBAL_VERSION, global)
			} else {
				P(WARING, "global Node.js version is %v, please use '%v' or '%v'. See '%v'.\n", util.UNKNOWN, "gnvm install latest -g", "gnvm install x.xx.xx -g", "gnvm help install")
			}
		} else {
			P(DEFAULT, "Node.js %v version is %v.\n", "global", global)
		}
	}

	if isLatest {
		if latest == util.UNKNOWN {
			P(WARING, "local  Node.js latest version is %v.\n", util.UNKNOWN)
		} else {
			P(DEFAULT, "Node.js %v version is %v.\n", "latest", latest)
		}
		remoteVersion := util.GetLatVer(latURL)
		if remoteVersion == "" {
			P(ERROR, "get remote %v Node.js %v error, please check your input. See '%v'.\n", config.GetConfig(config.REGISTRY), "latest version", "gnvm help config")
			return
		}
		if latest == util.UNKNOWN {
			P(NOTICE, "remote Node.js latest version is %v from %v.\n", remoteVersion, config.GetConfig(config.REGISTRY))
			config.SetConfig(config.LATEST_VERSION, remoteVersion)
			P(DEFAULT, "Set success, local Node.js %v version is %v.\n", util.LATEST, remoteVersion)
			return
		}
		v1 := util.FormatNodeVer(latest)
		v2 := util.FormatNodeVer(remoteVersion)
		if v1 < v2 {
			cp := CP{Red, false, None, false, ">"}
			P(WARING, "remote Node.js latest version %v %v local Node.js latest version %v, suggest to upgrade, usage '%v'.\n", remoteVersion, cp, latest, "gnvm update latest")
		}
	}
}

/*
 Print gnvm.exe version

 Param:
 	- remote: when remote == true, print CHANGELOG

*/
func Version(remote, detail bool) {

	defer func() {
		if err := recover(); err != nil {
			msg := fmt.Sprintf("'%v' an error has occurred. please check. \nError: ", "gnvm version -r")
			Error(ERROR, msg, err)
			os.Exit(0)
		}
	}()

	localVersion, arch := config.VERSION, "32 bit"
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

	code, res, err := curl.Get("http://ksria.com/gnvm/CHANGELOG.md")
	if code != 0 {
		panic(err)
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
		if line > 2 && detail {
			P(DEFAULT, content)
		}

		return false
	}

	if err := curl.ReadLine(res.Body, versionFunc); err != nil && err != io.EOF {
		panic(err)
	}
}
