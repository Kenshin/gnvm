package config

import (
	// lib
	"github.com/tsuru/config"

	// go
	"fmt"
	"os"
)

const (
	VERSION = "0.1.0"
	CONFIG  = ".gnvmrc"
	NEWLINE = "\n"
	UNKNOWN = "unknown"

	REGISTRY_KEY = "registry: "
	REGISTRY_VAL = "http://nodejs.org/dist/"

	NODEROOT     = "noderoot"
	NODEROOT_KEY = "noderoot: "
	NODEROOT_VAL = "root"

	GLOBAL_VERSION     = "globalversion"
	GLOBAL_VERSION_KEY = "globalversion: "
	GLOBAL_VERSION_VAL = "unknown"

	LATEST_VERSION     = "latestversion"
	LATEST_VERSION_KEY = "latestversion: "
	LATEST_VERSION_VAL = "unknown"

	CURRENT_VERSION     = "currentversion"
	CURRENT_VERSION_KEY = "currentversion: "
	CURRENT_VERSION_VAL = "unknown"
)

func init() {

	// config file is exist
	file, err := os.Open(CONFIG)
	defer file.Close()
	if err != nil && os.IsNotExist(err) {
		// print error
		fmt.Println("Config file is not exist.")

		// create .gnvmrc file and write
		createConfig()
	}

	// read config
	readConfig()

}

func createConfig() {

	// create file
	file, err := os.Create(CONFIG)
	defer file.Close()
	if err != nil {
		fmt.Println("Config file create Error: " + err.Error())
		return
	}

	//write init config
	_, fileErr := file.WriteString(REGISTRY_KEY + REGISTRY_VAL + NEWLINE + NODEROOT_KEY + NODEROOT_VAL + NEWLINE + GLOBAL_VERSION_KEY + GLOBAL_VERSION_VAL + NEWLINE + LATEST_VERSION_KEY + LATEST_VERSION_VAL + NEWLINE + CURRENT_VERSION_KEY + CURRENT_VERSION_VAL)
	if fileErr != nil {
		fmt.Println("Write Config file Error: " + fileErr.Error())
		return
	}

	// success
	fmt.Println("Config file create success.")

}

func readConfig() {
	if err := config.ReadConfigFile(CONFIG); err != nil {
		// print error
		fmt.Println("Read Config file Error: ", err.Error())

		return
	}

	fmt.Println("Read Config file success.")
}

func SetConfig(key string, value interface{}) string {

	// set new value
	config.Set(key, value)

	// delete old config
	if err := os.Remove(CONFIG); err != nil {
		// print error
		fmt.Println("Remove Config Error: ", err.Error())
	}

	// write new config
	if err := config.WriteConfigFile(CONFIG, 0777); err != nil {
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
