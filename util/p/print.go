package p

import (

	// lib
	"github.com/daviddengcn/go-colortext"

	// go
	"fmt"
	"os"
	"strings"
)

const (
	DEFAULT = ""
	WARING  = "Waring"
	ERROR   = "Error"
	NOTICE  = "Notice"
	SPLIT   = "%v"
)

/*
 * Print color console message
 *
 * flag   : include: 'Waring', 'Error'
 * message: print content
 * args   : variable parameter
 *
 * e.g. P( "Waring", Remote latest version [%v] = latest version [%v].\n", param1, param2 )
 *
 */
func P(flag, message string, args ...interface{}) {

	// try catch
	defer func() {
		if err := recover(); err != nil {
			Error(ERROR, "util/print.go an error has occurred. Error: ", err)
			os.Exit(0)
		}
	}()

	// set state
	stateColor(flag)

	// set key message
	msgArr := strings.Split(message, SPLIT)
	for k, v := range msgArr {
		fmt.Print(v)
		if k < len(args) {
			normalColor(args[k])
		}
	}

	if !strings.HasSuffix(message, "\n") {
		fmt.Printf("\n")
	}

}

/*
 * Print error color console message
 *
 * flag   : include: 'Waring', 'Error'
 * message: print content
 * err    : err content
 *
 * e.g. Error(ERROR, "util/print.go an error has occurred. Error: ", err)
 *
 */
func Error(flag, message string, err interface{}) {

	// set flag
	stateColor(flag)

	// color message
	ct.ChangeColor(ct.Red, false, ct.Green, false)
	fmt.Printf(message)

	// print err
	fmt.Println(err)

	// reset color
	ct.ResetColor()
}

func stateColor(state string) {
	switch state {
	case NOTICE:
		ct.ChangeColor(ct.Blue, false, ct.White, false)
		fmt.Printf("Notice: ")
	case WARING:
		ct.ChangeColor(ct.Green, false, ct.Red, false)
		fmt.Printf("Waring: ")
	case ERROR:
		ct.ChangeColor(ct.Red, false, ct.Green, false)
		fmt.Printf("Error: ")
	default:
		//ct.ChangeColor(ct.Blue, false, ct.White, false)
		//fmt.Printf("Notice: ")
	}
	ct.ResetColor()
}

func normalColor(msg interface{}) {
	ct.ChangeColor(ct.Green, true, ct.None, false)
	fmt.Printf(msg.(string))
	ct.ResetColor()
}
