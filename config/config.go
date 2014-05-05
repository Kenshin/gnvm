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
	REGISTRY = "http://nodejs.org/dist/"
	NODEROOT = ""
)

func init() {
	config, err := toml.LoadFile(CONFIG)
	if err != nil {
		fmt.Println("Error ", err.Error())
		createConfigFile()
	} else {
		registry := config.Get("registry").(string)
		fmt.Println(" registry is " + registry)
	}
}

func createConfigFile() {
	file, err := os.Create(CONFIG)
	if err != nil {
		fmt.Println(".gnvmrc create fail.")
	} else {
		file.WriteString("registry=" + REGISTRY)
		file.WriteString("\n")
		file.WriteString("noderoot=" + NODEROOT)
		fmt.Println(".gnvmrc file create success.")
	}

}
