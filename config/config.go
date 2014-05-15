package config

import (
	// lib
	"github.com/tsuru/config"

	// go
	"fmt"
	"os"
	"regexp"
	"strings"

	// local
	"gnvm/util"
)

var configPath string

const (
	VERSION  = "0.1.0"
	CONFIG   = ".gnvmrc"
	NEWLINE  = "\n"
	UNKNOWN  = "unknown"
	LATEST   = "latest"
	NODELIST = "npm-versions.txt"

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
		// print error
		fmt.Printf("Waring: Config file [%v] is not exist.\n", configPath)

		// create .gnvmrc file and write
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
		fmt.Println("Config file create Error: " + err.Error())
		return
	}

	//write init config
	_, fileErr := file.WriteString(REGISTRY_KEY + REGISTRY_VAL + NEWLINE + NODEROOT_KEY + configPath + NEWLINE + GLOBAL_VERSION_KEY + GLOBAL_VERSION_VAL + NEWLINE + LATEST_VERSION_KEY + LATEST_VERSION_VAL)
	if fileErr != nil {
		fmt.Println("Write Config file Error: " + fileErr.Error())
		return
	}

	// success
	fmt.Printf("Config file [%v] create success.\n", configPath)

}

func readConfig() {
	if err := config.ReadConfigFile(configPath); err != nil {
		fmt.Println("Read Config file Error: ", err.Error())
		return
	}
}

func SetConfig(key string, value interface{}) string {

	if key == "registry" {

		reg := regexp.MustCompile(`(http|https)://([\w-]+\.)+[\w-]+(/[\w- ./?%&=]*)?`)

		switch {
		case !reg.MatchString(value.(string)):
			fmt.Printf("Error: registry must url valid, current registry is [%v].\n", value.(string))
			return ""
		case !strings.HasSuffix(value.(string), "/"):
			value = value.(string) + "/"
		}
	}

	// set new value
	config.Set(key, value)

	// delete old config
	if err := os.Remove(configPath); err != nil {
		// print error
		fmt.Println("Remove Config Error: ", err.Error())
	}

	// write new config
	if err := config.WriteConfigFile(configPath, 0777); err != nil {
		// print error
		fmt.Println("Write Config Error: ", err.Error())
	}

	return value.(string)

}

func GetConfig(key string) string {
	value, err := config.GetString(key)
	if err != nil {
		// print error
		fmt.Println("GetConfig Error: " + err.Error())

		// value
		value = "unknown"
	}
	return value
}
