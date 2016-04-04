package config

import (
	// lib
	. "github.com/Kenshin/cprint"
	"github.com/tsuru/config"

	// go
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	// local
	"gnvm/util"
)

var configPath, globalversion, latsetversion string

const (
	VERSION  = "0.1.4 beta"
	CONFIG   = ".gnvmrc"
	NEWLINE  = "\n"
	UNKNOWN  = "unknown"
	LATEST   = "latest"
	NODELIST = "index.json"

	REGISTRY     = "registry"
	REGISTRY_KEY = "registry: "
	REGISTRY_VAL = "http://nodejs.org/dist/"
	TAOBAO       = "http://npm.taobao.org/mirrors/node/"

	NODEROOT     = "noderoot"
	NODEROOT_KEY = "noderoot: "
	NODEROOT_VAL = "root"

	GLOBAL_VERSION     = "globalversion"
	GLOBAL_VERSION_KEY = "globalversion: "
	GLOBAL_VERSION_VAL = UNKNOWN

	LATEST_VERSION     = "latestversion"
	LATEST_VERSION_KEY = "latestversion: "
	LATEST_VERSION_VAL = UNKNOWN

	//CURRENT_VERSION     = "currentversion"
	//CURRENT_VERSION_KEY = "currentversion: "
	//CURRENT_VERSION_VAL = UNKNOWN
)

func init() {

	// try catch
	defer func() {
		if err := recover(); err != nil {
			Error(ERROR, "gnvm.exe an error has occurred. please check. \nError: ", err)
			os.Exit(0)
		}
	}()

	// set config path
	configPath = util.GlobalNodePath + util.DIVIDE + CONFIG

	// config file is exist
	file, err := os.Open(configPath)
	defer file.Close()
	if err != nil && os.IsNotExist(err) {
		P(WARING, "config file %v is not exist.\n", configPath)
		createConfig()
	}

	// read config
	readConfig()

}

/*
 Create .gnvmrc file
*/
func createConfig() {

	// create file
	file, err := os.Create(configPath)
	defer file.Close()
	if err != nil {
		P(ERROR, "config file create Error: %v\n", err.Error())
		return
	}

	// get <root>/node.exe version
	version, err := util.GetNodeVer(util.GlobalNodePath)
	if err != nil {
		P(WARING, "not found global node.exe version, please use '%v'. See '%v'.\n", "gnvm install x.xx.xx -g", "gnvm help install")
		globalversion = GLOBAL_VERSION_VAL
	} else {
		globalversion = version
		// add suffix
		if runtime.GOARCH == "amd64" {
			if bit, err := util.Arch(util.GlobalNodePath); err == nil && bit == "x86" {
				globalversion += "-" + bit
			}
		}
	}

	//write init config
	_, fileErr := file.WriteString(REGISTRY_KEY + REGISTRY_VAL + NEWLINE + NODEROOT_KEY + util.GlobalNodePath + NEWLINE + GLOBAL_VERSION_KEY + globalversion + NEWLINE + LATEST_VERSION_KEY + LATEST_VERSION_VAL)
	if fileErr != nil {
		P(ERROR, "write config file Error: %v\n", fileErr.Error())
		return
	}

	P(DEFAULT, "Config file %v create success.\n", configPath)
}

/*
 Read .gnvmrc file
*/
func readConfig() {
	if err := config.ReadConfigFile(configPath); err != nil {
		P(ERROR, "read config file fail, please use '%v'. \nError: %v\n", "gnvm config INIT", err.Error())
		return
	}
}

/*
 Write config property value from .gnvmrc file

 Param:
 	- key:   config property, include: registry noderoot latestversion globalversion
 	- value: config property value

*/
func SetConfig(key string, value interface{}) string {
	if key == "registry" {
		if !strings.HasPrefix(value.(string), "http://") {
			P(WARING, "%v need %v", value.(string), "http://", "\n")
			value = "http://" + value.(string)
		}
		if !strings.HasSuffix(value.(string), "/") {
			value = value.(string) + "/"
		}
		reg, _ := regexp.Compile(`^https?:\/\/(w{3}\.)?([-a-zA-Z0-9.])+(\.[a-zA-Z]+)(:\d{1,4})?(\/)+`)
		if !reg.MatchString(value.(string)) {
			P(ERROR, "%v value %v must valid url.\n", "registry", value.(string))
			return ""
		}
	}

	// set new value
	config.Set(key, value)

	// delete old config
	if err := os.Remove(configPath); err != nil {
		P(ERROR, "remove config file Error: %v\n", err.Error())
	}

	// write new config
	if err := config.WriteConfigFile(configPath, 0777); err != nil {
		P(ERROR, "write config file Error: %v\n", err.Error())
	}

	return value.(string)
}

/*
 Read config property value from .gnvmrc file

 Param:
 	- key:   config property, include: registry noderoot latestversion globalversion

 Return:
 	- value: config property value

*/
func GetConfig(key string) string {
	value, err := config.GetString(key)
	if err != nil {
		value = UNKNOWN
	}
	return value
}

/*
 Init config property value from .gnvmrc file
*/
func ReSetConfig() {
	if newValue := SetConfig(REGISTRY, REGISTRY_VAL); newValue != "" {
		P(NOTICE, "%v      init success, new value is %v\n", REGISTRY, newValue)
	}
	if newValue := SetConfig(NODEROOT, util.GlobalNodePath); newValue != "" {
		P(NOTICE, "%v      init success, new value is %v\n", NODEROOT, newValue)
	}
	version, err := util.GetNodeVer(util.GlobalNodePath)
	if err != nil {
		P(WARING, "not found global node.exe version, please use '%v'. See '%v'.\n", "gnvm install x.xx.xx -g", "gnvm help install")
		globalversion = GLOBAL_VERSION_VAL
	} else {
		globalversion = version
		// add suffix
		if runtime.GOARCH == "amd64" {
			if bit, err := util.Arch(util.GlobalNodePath); err == nil && bit == "x86" {
				globalversion += "-" + bit
			}
		}
	}
	if newValue := SetConfig(GLOBAL_VERSION, globalversion); newValue != "" {
		P(NOTICE, "%v init success, new value is %v\n", GLOBAL_VERSION, newValue)
	}
}

/*
 Print all config property value from .gnvmrc file
*/
func List() {
	P(NOTICE, "config file path %v \n", configPath)
	f, err := os.Open(configPath)
	if err != nil {
		P(ERROR, "read config file fail, please use '%v'. \nError: %v\n", "gnvm config INIT", err.Error())
		return
	}
	buf := bufio.NewReader(f)
	for {
		line, _, err := buf.ReadLine()
		if err == io.EOF {
			break
		}
		arr := strings.SplitN(string(line), ":", 2)
		if len(arr) == 2 {
			P(DEFAULT, "gnvm config %v is %v\n", strings.TrimSpace(arr[0]), strings.TrimSpace(arr[1]))
		}
	}
}

/*
 Get io.js url

 Param:
 	- url:   config property, include: registry noderoot latestversion globalversion

 Return:
 	- value: config property value

*/
func GetIOURL(url string) string {
	if url == TAOBAO {
		url = strings.Replace(url, "/node", "/iojs", -1)
	} else if url == REGISTRY_VAL {
		url = strings.Replace(url, "nodejs.org", "iojs.org", -1)
	}
	return url
}

/*
 Verify config registry url structural correctness, include:
 	- url:  <url>
 	- json: <url>/index.json
*/
func Verify() {
	code := make(chan int)
	fail := make(chan interface{})
	finish := false
	registry := GetConfig(REGISTRY)
	wait := func() {
		wait := ""
		for {
			time.Sleep(time.Millisecond * 500)
			if finish {
				break
			}
			wait += "."
			fmt.Printf(wait)
		}
	}

	go verifyURL("url", registry, code, fail)
	go wait()

	for {
		select {
		case <-time.After(time.Second * 10):
			cp1 := CP{Red, false, None, false, "fail"}
			cp2 := CP{Red, false, None, false, "time out"}
			P(DEFAULT, "%v. \n", cp1)
			P(ERROR, "gnvm config registry %v vaild %v, Error: %v.", registry, cp1, cp2)
			return
		case value, ok := <-code:
			if ok && value == 200 {
				finish = true
				P(DEFAULT, "%v.\n", " ok")
				go verifyURL("json", registry+NODELIST, code, fail)
				finish = false
				go wait()
			} else if !ok {
				P(DEFAULT, "%v.\n", " ok")
				return
			}
		case value, _ := <-fail:
			cp := CP{Red, false, None, false, " fail"}
			if v, ok := value.(int); ok {
				P(DEFAULT, "%v, response code: %v.\n", cp, strconv.Itoa(v))
			} else {
				P(DEFAULT, "%v.\n", cp)
				Error(ERROR, "", value)
			}
			close(fail)
			finish = true
			return
		}
	}
}

func verifyURL(status string, url string, code chan int, fail chan interface{}) {
	P(NOTICE, "gnvm config registry %v valid ", url)
	time.Sleep(time.Second * 2)
	if resp, err := http.Get(url); err == nil {
		if resp.StatusCode == 200 {
			if status == "url" {
				code <- resp.StatusCode
			} else {
				close(code)
			}
		} else {
			fail <- resp.StatusCode
		}
	} else {
		fail <- err
	}
}
