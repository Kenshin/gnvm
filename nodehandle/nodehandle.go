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

/**
 * rootPath is gnvm.exe root path,     e.g <root>
 * rootNode is rootPath + "/node.exe", e.g. <root>/node.exe
 *
 * usePath  is use node version path,  e.g. <root>/x.xx.xx
 * useNode  is usePath + "/node.exe",  e.g. <root>/x.xx.xx/node.exe
 *
 * rootVersion is <root>/node.exe version
 * rootFolder  is <root>/rootVersion
 */
func Use(folder string, global bool) {

	// set rootPath and rootNode
	rootPath := getCurrentPath() + DIVIDE
	rootNode := rootPath + NODE
	//fmt.Println("Current path is " + rootPath)

	// set usePath and useNode
	usePath := rootPath + folder + DIVIDE
	useNode := usePath + NODE
	//fmt.Println("Node.exe path is " + usePath)

	// <root>/folder is exist
	if isDirExist(usePath) != true {
		fmt.Printf("%v version is not exist. Get local node.exe version see 'gnvm ls'.", folder)
		return
	}

	// <root>/node.exe is exist
	if isDirExist(rootNode) != true {
		fmt.Println("Not found global node version, please checkout. If not exist node.exe, See 'gnvm install latest'.")
		return
	}

	// get <root>/node.exe version
	rootVersion, err := getNodeVersion(rootPath)
	if err != nil {
		fmt.Println("Not found global node version, please checkout. If not exist node.exe, See 'gnvm install latest'.")
	}
	//fmt.Printf("root node.exe verison is: %v", rootVersion)
	rootFolder := rootPath + rootVersion

	// <root>/rootVersion is exist
	if isDirExist(rootPath+rootVersion) != true {

		// create rootVersion folder
		if err := cmd("md", rootVersion); err != nil {
			fmt.Printf("Create %v folder Error: %v", rootVersion, err.Error())
			return
		}

	}

	// copy rootNode to <root>/rootVersion
	if err := copy(rootNode, rootPath+rootVersion); err != nil {
		fmt.Printf("copy %v to %v folder Error: ", rootNode, rootPath+rootVersion, err.Error())
		return
	}

	// delete <root>/node.exe
	if err := os.Remove(rootNode); err != nil {
		fmt.Printf("remove %v to %v folder Error: ", rootNode, err.Error())
		return
	}

	// copy useNode to rootPath
	if err := copy(useNode, rootPath); err != nil {
		fmt.Printf("copy %v to %v folder Error: ", useNode, rootPath, err.Error())
		return
	}

	fmt.Printf("Set success, Current Node.exe version is [%v].", folder)

}
