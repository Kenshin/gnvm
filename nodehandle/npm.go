package nodehandle

import (
	/*
		"github.com/Kenshin/curl"
		"github.com/pierrre/archivefile/zip"

		// go
		//"log"
		"fmt"
		"io"
		"io/ioutil"
		"os"
		"os/exec"
		"runtime"
		"strconv"
		"time"
	*/

	// lib
	. "github.com/Kenshin/cprint"

	// go
	"fmt"
	"strings"

	// local
	"gnvm/util"
)

func InstallNPM(version string) {
	version = strings.ToLower(version)
	switch version {
	case util.LATEST:
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
