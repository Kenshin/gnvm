package nodehandle

import (

	// go
	"fmt"
	"os"
	"os/exec"
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
	return string(string(out[:])[1:]), err
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

	// nodePath is exist
	/*
		if isDirExist(nodePath) != true {
			fmt.Println("Node.exe version is not exist. Get local node.exe version see 'gnvm ls'.")
			return
		}
	*/

	/*
		app := "ls"
		cmd, err := exec.Command(app, []string{app, "-l"}, nil, "", "", "", "")

		if err != nil {
			fmt.Fprintln(os.Stderr, err.String())
			return
		}

		var b bytes.Buffer
		io.Copy(&b, cmd.Stdout)
		fmt.Println(b.String())

		cmd.Close()
	*/

	/*
		out, err := exec.Command("node", "--version").Output()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("The date is %s\n", out)
	*/

	/*
		//os.Setenv("FOO", "1")
		fmt.Println("FOO:", os.Getenv("FOO"))
		fmt.Println("BAR:", os.Getenv("BAR"))
	*/
	/*
		fmt.Println()
		for _, e := range os.Environ() {
			pair := strings.Split(e, "=")
			fmt.Println(pair[0])
		}
	*/

	// set env
	/*
		if err := os.Setenv(NODEHOME, nodePath); err != nil {
			fmt.Println("Set Env Error: " + err.Error())
			return
		}
		fmt.Println("NODE_HOME is " + os.Getenv(NODEHOME))

		// set path
		if err := os.Setenv("Path", os.Getenv(NODEHOME)); err != nil {
			fmt.Println("Set Env Error: " + err.Error())
			return
		}
		fmt.Println(os.Getenv("Path	"))
	*/

	/*
		cmd := exec.Command("set", NODEHOME, nodePath)
		if err := cmd.Run(); err != nil {
			fmt.Println("Exec set NODE_HOME Error: " + err.Error())
		}
		fmt.Println("NODE_HOME is " + os.Getenv(NODEHOME))
	*/

	/*
		if err := syscall.Setenv(NODEHOME, nodePath); err != nil {
			fmt.Println("Set Env Error: " + err.Error())
			return
		}
		value, _ := syscall.Getenv(NODEHOME)
		fmt.Println("NODE_HOME is " + value)
	*/

	/*
		output, err := exec.Command("gnvm2.bat").CombinedOutput()
		if err != nil {
			fmt.Println("sadfsafaf " + err.Error())
			return
		}
		fmt.Println(string(output))
	*/

	/*
		os.Setenv("FOO", "BAR")
		if err := syscall.Exec(os.Getenv("cmd"), []string{os.Getenv("cmd")}, syscall.Environ()); err != nil {
			fmt.Println("sdfdfafaf " + err.Error())
		}
		//fmt.Println(syscall.Environ())
	*/

	/*
		cmd := exec.Command("cmd", "/C", "del", "D:\\a.txt")
		if err := cmd.Run(); err != nil {
			fmt.Println("sadfsafaf " + err.Error())
		}
	*/
}
