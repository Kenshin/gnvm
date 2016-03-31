package main

import (
	"fmt"
	"gnvm/nodehandle"
	"gnvm/util"
	"testing"
)

func TestCurl(t *testing.T) {
	//testSearch()
	//testNodist()
	//testNPManage()
	//testGetNPMVer()
	testIsDirExist()
}

func testSearch() {
	nodehandle.Query("x.x.x")
	nodehandle.Query("0.10.x")
	nodehandle.Query("5.x.x")
	nodehandle.Query("5.0.0")
	nodehandle.Query(`/^5(\.([0]|[1-9]\d?)){2}$/`)
	nodehandle.Query("latest")
	nodehandle.Query("1.x.x")
	nodehandle.Query("1.1.x")
	nodehandle.Query("3.x.x")
	nodehandle.Query("3.3.x")
}

func testNodist() {
	if nl, err, code := nodehandle.New("http://npm.taobao.org/mirrors/iojs/index.json", nil); err != nil {
		fmt.Println(err)
		fmt.Println(code)
	} else {
		nl.Detail(0)
	}
}

func testNPManage() {
	name := `v3.8.5.zip`
	npm := new(nodehandle.NPMange)
	npm.New().CleanAll()
	npm.SetZip(name)
	npm.Unzip()
	npm.Install()
	fmt.Println(npm)
}

func testGetNPMVer() {
	url := "http://npm.taobao.org/mirrors/node/index.json"
	ver := "5.9.0"
	if nd, err := nodehandle.FindNodeDetailByVer(url, ver); err == nil {
		fmt.Println(nd)
	}
}

func testIsDirExist() {
	// empty
	fmt.Println(util.IsDirExist(""))
	// no exist
	fmt.Println(util.IsDirExist("/Users/kenshin/Work/28-GO/01-work/src/gnvm/node_modules"))
	fmt.Println(util.IsDirExist("/Users/kenshin/Work/28-GO/01-work/src/gnvm/node_modules/npm"))
	// exist
	fmt.Println(util.IsDirExist("/Users/kenshin/Work/28-GO/01-work/src/gnvm/"))
	fmt.Println(util.IsDirExist("/Users/kenshin/Work/28-GO/01-work/src/gnvm"))
	// not valid path
	fmt.Println(util.IsDirExist("gnvm"))
}
