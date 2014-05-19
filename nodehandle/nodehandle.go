package nodehandle

import (

	// go
	//"log"
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	// local
	"gnvm/config"
	"gnvm/util"
)

const (
	DIVIDE = "\\"
	NODE   = "node.exe"
)

var rootPath string

func init() {
	rootPath = util.GlobalNodePath + DIVIDE
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

	rootNodeExist := true

	// get true folder, e.g. folder is latest return x.xx.xx
	folder = GetTrueVersion(folder, true)

	if folder == config.UNKNOWN {
		fmt.Println("Waring: Unassigned Node.js latest version. See 'gnvm install latest'.")
		return false
	}

	// set rootNode
	rootNode := rootPath + NODE
	//log.Println("Root node path is: " + rootNode)

	// set usePath and useNode
	usePath := rootPath + folder + DIVIDE
	useNode := usePath + NODE
	//log.Println("Use node.exe path is: " + usePath)

	// get <root>/node.exe version
	rootVersion, err := util.GetNodeVersion(rootPath)
	if err != nil {
		fmt.Println("Waring: not found global node version, please use 'gnvm install x.xx.xx -g'. See 'gnvm help install'.")
		rootNodeExist = false
	}
	//log.Printf("Root node.exe verison is: %v \n", rootVersion)

	// <root>/folder is exist
	if isDirExist(usePath) != true {
		fmt.Printf("[%v] folder is not exist. Get local node.exe version. See 'gnvm ls'.\n", folder)
		return false
	}

	// check folder is rootVersion
	if folder == rootVersion {
		fmt.Printf("Current node.exe version is [%v], not re-use. See 'gnvm node-version'.\n", folder)
		return false
	}

	// set rootFolder
	rootFolder := rootPath + rootVersion

	// <root>/rootVersion is exist
	if isDirExist(rootFolder) != true {

		// create rootVersion folder
		if err := cmd("md", rootFolder); err != nil {
			fmt.Printf("Create %v folder Error: %v.\n", rootVersion, err.Error())
			return false
		}

	}

	if rootNodeExist {
		// copy rootNode to <root>/rootVersion
		if err := copy(rootNode, rootFolder); err != nil {
			fmt.Printf("copy %v to %v folder Error: %v.\n", rootNode, rootFolder, err.Error())
			return false
		}

		// delete <root>/node.exe
		if err := os.Remove(rootNode); err != nil {
			fmt.Printf("remove %v folder Error: %v.\n", rootNode, err.Error())
			return false
		}

	}

	// copy useNode to rootPath
	if err := copy(useNode, rootPath); err != nil {
		fmt.Printf("copy %v to %v folder Error: %v.\n", useNode, rootPath, err.Error())
		return false
	}

	fmt.Printf("Set success, Current Node.exe version is [%v].\n", folder)

	return true

}

func VerifyNodeVersion(version string) bool {
	result := true
	if version == config.UNKNOWN {
		return true
	}
	arr := strings.Split(version, ".")
	if len(arr) != 3 {
		return false
	}
	for _, v := range arr {
		_, err := strconv.ParseInt(v, 10, 0)
		if err != nil {
			result = false
			break
		}
	}
	return result
}

func GetTrueVersion(latest string, isPrint bool) string {
	if latest == config.LATEST {
		latest = config.GetConfig(config.LATEST_VERSION)
		if isPrint {
			fmt.Printf("Current latest version is [%v] \n", latest)

		}
	}
	return latest
}

func NodeVersion(args []string, remote bool) {

	// try catch
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
	}()

	latest := config.GetConfig(config.LATEST_VERSION)
	global := config.GetConfig(config.GLOBAL_VERSION)

	if len(args) == 0 || len(args) > 1 {
		fmt.Printf("Node.exe latest verson is [%v].\n", latest)
		fmt.Printf("Node.exe global verson is [%v].\n", global)
	} else {
		switch {
		case args[0] == "global":
			fmt.Printf("Node.exe global verson is [%v].\n", global)
		case args[0] == "latest" && !remote:
			fmt.Printf("Node.exe latest verson is [%v].\n", latest)
		case args[0] == "latest" && remote:
			remoteVersion := getLatestVersionByRemote()
			if remoteVersion == "" {
				fmt.Printf("Error: get remote [%v] latest version error, please check. See 'gnvm config help'.\n", config.GetConfig("registry")+config.LATEST+"/"+config.NODELIST)
				fmt.Printf("Node.exe latest verson is [%v].\n", latest)
				return
			}
			fmt.Printf("Node.exe remote [%v] verson is [%v].\n", config.GetConfig("registry"), remoteVersion)
			fmt.Printf("Node.exe latest verson is [%v].\n", latest)
		}
	}

	switch {
	case len(args) == 0 && (global == config.UNKNOWN || latest == config.UNKNOWN):
		fmt.Printf("Waring: when version is [%v], please Use 'gnvm config INIT'. See 'gnvm help config'.\n", config.UNKNOWN)
	case len(args) > 0 && args[0] == "latest" && latest == config.UNKNOWN:
		fmt.Printf("Waring: when version is [%v], please Use 'gnvm config INIT'. See 'gnvm help config'.\n", config.UNKNOWN)
	case len(args) > 0 && args[0] == "global" && global == config.UNKNOWN:
		fmt.Printf("Waring: when version is [%v], please Use 'gnvm config INIT'. See 'gnvm help config'.\n", config.UNKNOWN)
	}
}

func Uninstall(folder string) {

	// set removePath
	removePath := rootPath + folder

	if folder == config.UNKNOWN {
		fmt.Println("Waring: Unassigned Node.js latest version. See 'gnvm install latest'.")
		return
	}

	// rootPath/version is exist
	if isDirExist(removePath) != true {
		fmt.Printf("Waring: [%v] folder is not exist. Get local node.exe version. See 'gnvm ls'.\n", folder)
		return
	}

	// remove rootPath/version folder
	if err := os.RemoveAll(removePath); err != nil {
		fmt.Printf("Uinstall [%v] fail, Error: %v.\n", folder, err.Error())
		return
	}

	fmt.Printf("Node.exe version [%v] uninstall success.\n", folder)
}

func LS(isPrint bool) ([]string, error) {
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
		if ok := VerifyNodeVersion(version); ok {

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
					fmt.Println("v" + version + desc)
				}
			}
		}

		// return
		return nil
	})

	// show error
	if err != nil {
		fmt.Printf("'gnvm ls' Error: %v.\n", err.Error())
		return lsArr, err
	}

	// version is exist
	if !existVersion {
		fmt.Println("Waring: don't have any available version, please check. See 'gnvm help install'.")
	}

	return lsArr, err
}

func LsRemote() {

	// set exist version
	isExistVersion := false

	// set url
	registry := config.GetConfig("registry")
	url := registry + config.NODELIST

	// print
	fmt.Println("Read all Node.exe version list from " + url + ", please wait.")

	// try catch
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("'gnvm ls --remote' an error has occurred. please check registry: [%v], Error: %v.\n", url, err)
			fmt.Println(err)
			os.Exit(0)
		}
	}()

	// get res
	res, err := http.Get(url)

	// close
	defer res.Body.Close()

	// err
	if err != nil {
		panic(err)
	}

	// check state code
	if res.StatusCode != 200 {
		fmt.Printf("registry [%v] an [%v] error occurred, please check. See 'gnvm config help'.\n", url, res.StatusCode)
		return
	}

	// set buff
	buff := bufio.NewReader(res.Body)

	for {
		// set line
		line, err := buff.ReadString('\n')

		// when EOF or err break
		if err != nil || err == io.EOF {
			break
		}

		// replace '\n'
		line = strings.Replace(line, "\n", "", -1)

		// splite 'vx.xx.xx  1.1.0-alpha-2'
		args := strings.Split(line, " ")

		if ok := VerifyNodeVersion(args[0][1:]); ok {
			isExistVersion = true
			// print all node.exe version
			fmt.Println(args[0])
		}
	}

	if !isExistVersion {
		fmt.Printf("Not found any Node.exe version list from %v, please check it.\n", url)
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
			fmt.Println(err)
			os.Exit(0)
		}
	}()

	for _, v := range args {

		if v == config.LATEST {

			version := getLatestVersionByRemote()
			if version == "" {
				fmt.Println("Get latest version error, please check. See 'gnvm config help'.")
				break
			}

			// set v
			v = version
			currentLatest = version
			fmt.Printf("Current latest version is [%v]\n", version)
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

func Update(global bool) {

	// try catch
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
	}()

	localVersion := config.GetConfig(config.LATEST_VERSION)
	fmt.Printf("local latest version is [%v].\n", localVersion)

	remoteVersion := getLatestVersionByRemote()
	if remoteVersion == "" {
		fmt.Println("Get latest version error, please check. See 'gnvm config help'.")
		return
	}
	fmt.Printf("remote [%v] latest version is [%v].\n", config.GetConfig("registry"), remoteVersion)

	local, _ := util.ConverFloat(localVersion)
	remote, _ := util.ConverFloat(remoteVersion)

	var args []string
	args = append(args, remoteVersion)

	switch {
	case localVersion == config.UNKNOWN:
		fmt.Println("Waring: local latest version undefined.")
		if code := Install(args, global); code == 0 || code == 2 {
			config.SetConfig(config.LATEST_VERSION, remoteVersion)
			fmt.Printf("Update latest success, current latest version is [%v].\n", remoteVersion)
		}
	case local == remote:
		fmt.Printf("Remote latest version [%v] = latest version [%v].\n", remoteVersion, localVersion)
	case local > remote:
		fmt.Printf("Waring: local latest version [%v] > remote latest version [%v], please check your registry. See 'gnvm help config'.\n", localVersion, remoteVersion)
	case local < remote:
		fmt.Printf("Remote latest version [%v] > local latest version [%v].\n", remoteVersion, localVersion)
		if code := Install(args, global); code == 0 || code == 2 {
			config.SetConfig(config.LATEST_VERSION, remoteVersion)
			fmt.Printf("Update latest success, current latest version is [%v].\n", remoteVersion)
		}
	}
}

/*
 * return code
 * 0: success
 * 1: status code != 200
 * 2: folder exist
 * 3: remove folder error
 * 4: create folder error
 * 5: download node.exe error
 *
 */
func download(version string) int {

	// get current os arch
	amd64 := "/"
	if runtime.GOARCH == "amd64" {
		amd64 = "/x64/"
	}

	// set url
	registry := config.GetConfig("registry")
	url := registry + "v" + version + amd64 + NODE

	// get res
	res, err := http.Get(url)

	// close
	defer res.Body.Close()

	// err
	if err != nil {
		panic(err)
	}

	// check state code
	if res.StatusCode != 200 {
		fmt.Printf("Downlaod url [%v] an [%v] error occurred, please check. See 'gnvm config help'.\n", url, res.StatusCode)
		return 1
	}

	// rootPath/version/node.exe is exist
	if _, err := util.GetNodeVersion(rootPath + version + DIVIDE); err == nil {
		fmt.Printf("Waring: [%v] folder exist.\n", version)
		return 2
	} else {
		if err := os.RemoveAll(rootPath + version); err != nil {
			fmt.Printf("Remove [%v] fail, Error: %v\n", version, err.Error())
			return 3
		}
		fmt.Printf("Remove empty [%v] folder success.\n", version)
	}

	// rootPath/version is exist
	if isDirExist(rootPath+version) != true {
		if err := os.Mkdir(rootPath+version, 0777); err != nil {
			panic(err)
		}
	}

	// create file
	file, createErr := os.Create(rootPath + version + DIVIDE + NODE)
	if createErr != nil {
		fmt.Println("Create file error, Error: " + createErr.Error())
		return 4
	}
	defer file.Close()

	fmt.Printf("Start download node.exe version [%v] from %v.\n", version, url)

	// loop buff to file
	buf := make([]byte, res.ContentLength)
	var m float32
	isShow, oldCurrent := false, 0
	for {
		n, err := res.Body.Read(buf)

		// write complete
		if n == 0 {
			fmt.Println("100% \nEnd download.")
			break
		}

		//error
		if err != nil {
			panic(err)
		}

		/* show console e.g.
		 * Start download node.exe version [x.xx.xx] from http://nodejs.org/dist/.
		 * 10% 20% 30% 40% 50% 60% 70% 80% 90% 100%
		 * End download.
		 */
		m = m + float32(n)
		current := int(m / float32(res.ContentLength) * 100)

		if current > oldCurrent {
			switch current {
			case 10, 20, 30, 40, 50, 60, 70, 80, 90:
				isShow = true
			}

			if isShow {
				fmt.Printf("%d%v", current, "% ")
			}

			isShow = false
		}

		oldCurrent = current

		file.WriteString(string(buf[:n]))
	}

	// valid download exe
	fi, err := file.Stat()
	if err == nil {
		if fi.Size() != res.ContentLength {
			fmt.Printf("Error: Downlaod node.exe version [%v] size error, please check your network and run 'gnvm uninstall %v'.\n", version, version)
			return 5
		}
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
