package config

import (
	// lib
	"github.com/tsuru/config"

	// go
	"os"
	"regexp"
	"strings"

	// local
	"gnvm/util"
	. "gnvm/util/p"
)

var configPath, globalversion, latsetversion string

const (
	VERSION  = "0.1.0"
	CONFIG   = ".gnvmrc"
	NEWLINE  = "\n"
	UNKNOWN  = "unknown"
	LATEST   = "latest"
	NODELIST = "npm-versions.txt"

	REGISTRY     = "registry"
	REGISTRY_KEY = "registry: "
	REGISTRY_VAL = "http://nodejs.org/dist/"

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

	// set config path
	configPath = util.GlobalNodePath + "\\" + CONFIG

	// config file is exist
	file, err := os.Open(configPath)
	defer file.Close()
	if err != nil && os.IsNotExist(err) {
		P(WARING, "config file [%v] is not exist.\n", configPath)
		createConfig()
	}

	// read config
	readConfig()

}

func createConfig() {

	// create file
	file, err := os.Create(configPath)
	defer file.Close()
	if err != nil {
		P(ERROR, "config file create Error: %v", err.Error())
		return
	}

	// get <root>/node.exe version
	version, err := util.GetNodeVersion(util.GlobalNodePath + "\\")
	if err != nil {
		P(WARING, "not found global node version, please use 'gnvm install x.xx.xx -g'. See 'gnvm help install'.")
		globalversion = GLOBAL_VERSION_VAL
	} else {
		globalversion = version
	}

	//write init config
	_, fileErr := file.WriteString(REGISTRY_KEY + REGISTRY_VAL + NEWLINE + NODEROOT_KEY + util.GlobalNodePath + NEWLINE + GLOBAL_VERSION_KEY + globalversion + NEWLINE + LATEST_VERSION_KEY + LATEST_VERSION_VAL)
	if fileErr != nil {
		P(ERROR, "write config file Error: %v", fileErr.Error())
		return
	}

	// success
	P(DEFAULT, "Config file [%v] create success.\n", configPath)
	P(WARING, "latest version is [%v], please use 'gnvm update latest' or 'gnvm install latest'.\n", UNKNOWN)

}

func readConfig() {
	if err := config.ReadConfigFile(configPath); err != nil {
		P(ERROR, "read Config file. Error: %v", err.Error())
		return
	}
}

func SetConfig(key string, value interface{}) string {

	if key == "registry" {

		reg := regexp.MustCompile(`(http|https)://([\w-]+\.)+[\w-]+(/[\w- ./?%&=]*)?`)

		switch {
		case !reg.MatchString(value.(string)):
			P(ERROR, "registry value [%v] must url valid.\n", value.(string))
			return ""
		case !strings.HasSuffix(value.(string), "/"):
			value = value.(string) + "/"
		}
	}

	// set new value
	config.Set(key, value)

	// delete old config
	if err := os.Remove(configPath); err != nil {
		P(ERROR, "remove config file Error: %v", err.Error())
	}

	// write new config
	if err := config.WriteConfigFile(configPath, 0777); err != nil {
		P(ERROR, "write config file Error: %v", err.Error())
	}

	return value.(string)

}

func GetConfig(key string) string {
	value, err := config.GetString(key)
	if err != nil {
		P(ERROR, "get config Error: %v", err.Error())
		value = UNKNOWN
	}
	return value
}

func ReSetConfig() {
	SetConfig(REGISTRY, REGISTRY_VAL)
	SetConfig(NODEROOT, util.GlobalNodePath)

	version, err := util.GetNodeVersion(util.GlobalNodePath + "\\")
	if err != nil {
		P(WARING, "not found global node version, please use 'gnvm install x.xx.xx -g'. See 'gnvm help install'.")
		globalversion = GLOBAL_VERSION_VAL
	} else {
		globalversion = version
	}
	SetConfig(GLOBAL_VERSION, globalversion)

	// set url
	url := REGISTRY_VAL + "latest/" + util.SHASUMS
	if latest := util.GetLatestVersion(url); latest != "" {
		latsetversion = latest
	} else {
		latsetversion = LATEST_VERSION_VAL
	}
	SetConfig(LATEST_VERSION, latsetversion)

	P(DEFAULT, "Config file [%v] init success.\n", configPath)
}
