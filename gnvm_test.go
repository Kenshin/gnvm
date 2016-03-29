package main

import (
	"fmt"
	"gnvm/nodehandle"
	"testing"
)

func TestCurl(t *testing.T) {
	//testSearch()
	//testNodist()
	testUnzip()
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

func testUnzip() {
	path := `C:\Users\Kenshin\Documents\DevTools\aaa\`
	name := `npm-3.8.5.zip`
	dest := `npm-3.8.5`
	if code, err := nodehandle.Unzip(path+name, path+dest); err != nil {
		fmt.Println(code)
		fmt.Println(err)
	}
}
