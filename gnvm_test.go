package main

import (
	"gnvm/nodehandle"
	"testing"
)

func TestCurl(t *testing.T) {
	testSearch()
}

func testSearch() {
	nodehandle.Query("x.x.x")
	nodehandle.Query("5.x.x")
	nodehandle.Query("0.10.x")
	nodehandle.Query("5.0.0")
	nodehandle.Query(`/^5(\.([0]|[1-9]\d?)){2}$/`)
	nodehandle.Query("latest")
	nodehandle.Query("latest")
}
