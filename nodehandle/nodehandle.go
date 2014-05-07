package nodehandle

import (

	// go
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const (
	DIVIDE = "\\"
	//NODEHOME = "NODE_HOME_2"
	PATH = "path"
	NODE = "node.exe"
)

var rootPath string

func getCurrentPath() string {
	path, err := os.Getwd()
	if err != nil {
		fmt.Println("Get current path Error: " + err.Error())
		return ""
	}
	return path
}

func isDirExist(path string) bool {
	fmt.Println(path)
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	} else {
		// return file.IsDir()
		return true
	}
}

func getNodeVersion(path string) (string, error) {
	out, err := exec.Command(path+"node", "--version").Output()
	//string(out[:]) bytes to string
	// replace \r\n
	newout := strings.Replace(string(string(out[:])[1:]), "\r\n", "", -1)
	return newout, err
}

func cmd(name, arg string) error {
	_, err := exec.Command("cmd", "/C", name, arg).Output()
	return err
}

func copy(src, dest string) error {
	_, err := exec.Command("cmd", "/C", "copy", "/y", src, dest).Output()
	return err
}

func Use(folderName string, global bool) {

	// get current path
	rootPath := getCurrentPath() + DIVIDE
	//fmt.Println("Current path is " + rootPath)

	// get use node path
	nodePath := rootPath + folderName + DIVIDE
	//fmt.Println("Node.exe path is " + nodePath)

	// <root>/foldName/ is exist
	if isDirExist(nodePath) != true {
		fmt.Printf("%v version is not exist. Get local node.exe version see 'gnvm ls'.", folderName)
		return
	}

	// <root>/node.exe is exist
	if isDirExist(rootPath+NODE) != true {
		fmt.Println("Not found global node version, please checkout. If not exist node.exe, See 'gnvm install latest'.")
		return
	}

	// current node.exe
	currentVersion, err := getNodeVersion(rootPath)
	if err != nil {
		fmt.Println("Not found global node version, please checkout. If not exist node.exe, See 'gnvm install latest'.")
	}
	//fmt.Printf("root node.exe verison is: %v", currentVersion)

	// <root>/currentVersion is exist?
	if isDirExist(rootPath+currentVersion+DIVIDE) != true {

		// create currentversion folder
		if err := cmd("md", currentVersion); err != nil {
			fmt.Printf("Create %v folder Error: %v", currentVersion, err.Error())
			return
		}

	}

	// copy currentVersion to <root>/currentVersion
	if err := copy(rootPath+NODE, rootPath+currentVersion); err != nil {
		fmt.Printf("copy %v to %v folder Error: ", rootPath+NODE, rootPath+currentVersion, err.Error())
		return
	}

	// delete <root>/node.exe
	if err := os.Remove(rootPath + NODE); err != nil {
		fmt.Printf("remove %v to %v folder Error: ", rootPath+NODE, err.Error())
		return
	}

	// copy currentVersion to <root>/currentVersion
	if err := copy(nodePath+NODE, rootPath); err != nil {
		fmt.Printf("copy %v to %v folder Error: ", nodePath+NODE, rootPath, err.Error())
		return
	}

	fmt.Printf("Set success, Current Node.exe version is [%v].", folderName)

}
