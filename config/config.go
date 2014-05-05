package config

import (
	// lib
	"github.com/pelletier/go-toml"

	// go
	"fmt"
	"os"
)

const (
	VERSION  = "0.1.0"
	CONFIG   = ".gnvmrc"
	REGISTRY = `"http://nodejs.org/dist/"`
	NODEROOT = `""`
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
		registry := config.Get("noderoot").(string)
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
	_, fileErr := file.WriteString("registry=" + REGISTRY + "\n" + "noderoot=" + NODEROOT)
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
