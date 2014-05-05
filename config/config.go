package config

import (
	// lib
	"github.com/pelletier/go-toml"

	// go
	"fmt"
	"os"
)

const (
	NEWLINE = "\n"
	VERSION = "0.1.0"
	CONFIG  = ".gnvmrc"

	REGISTRY_KEY = "registry="
	REGISTRY_VAL = `"http://nodejs.org/dist/"`
	NODEROOT_KEY = `noderoot=`
	NODEROOT_VAL = `""`

	GLOBAL_VERSION_KEY = "globalversion="
	GLOBAL_VERSION_VAL = `"unknown"`

	LATEST_VERSION_KEY = "latestversion="
	LATEST_VERSION_VAL = `"unknown"`

	CURRENT_VERSION_KEY = "currentversion="
	CURRENT_VERSION_VAL = `"unknown"`
)

func init() {
	config, err := toml.LoadFile(CONFIG)
	if err != nil {

		// print error
		fmt.Println("Error ", err.Error())

		// create .gnvmrc file and write
		createConfigFile()
	} else {

		// get registry
		registry := config.Get("registry").(string)
		fmt.Println("registry is " + registry)

		// get nodeversion
		noderoot := config.Get("noderoot").(string)
		fmt.Println("noderoot is " + noderoot)
	}
}

func createConfigFile() {

	// create file
	file, err := os.Create(CONFIG)
	if err != nil {
		fmt.Println(".gnvmrc create fail, error is" + err.Error())
		return
	}

	//write file
	_, fileErr := file.WriteString(REGISTRY_KEY + REGISTRY_VAL + NEWLINE + NODEROOT_KEY + NODEROOT_VAL + NEWLINE + GLOBAL_VERSION_KEY + GLOBAL_VERSION_VAL + NEWLINE + LATEST_VERSION_KEY + LATEST_VERSION_VAL + NEWLINE + CURRENT_VERSION_KEY + CURRENT_VERSION_VAL)
	if fileErr != nil {
		fmt.Println("write .gnvmrc fail" + fileErr.Error())
		return
	}

	// close file
	file.Close()

	// defear
	defer file.Close()

	// success
	fmt.Println(".gnvmrc file create success.")

}
