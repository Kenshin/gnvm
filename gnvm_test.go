package main

import (
	"fmt"
	"os/exec"
	"testing"
)

func Test(t *testing.T) {

	out, err := exec.Command("go", "run", "gnvm.go","version").Output()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(out[:]))
}
