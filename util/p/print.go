package p

import (

	// lib
	"github.com/daviddengcn/go-colortext"

	// go
	"fmt"
)

func Exec() {
	ct.ChangeColor(ct.Green, false, ct.Red, false)
	fmt.Printf("Waring: ")
	ct.ResetColor()
	fmt.Printf("use format error, the correct format is 'gnvm uninstall npm'. See 'gnvm help uninstall'.")
}