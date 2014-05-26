package p

import (

	// lib
	"github.com/daviddengcn/go-colortext"

	// go
	"fmt"
	"os"
	"reflect"
	"strings"
)

const (
	DEFAULT = ""
	WARING  = "Waring"
	ERROR   = "Error"
	NOTICE  = "Notice"
	SPLIT   = "%v"
)

const (
	None = iota
	Black
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

type CC struct {
	FgColor  int
	FgBright bool
	BgColor  int
	BgBright bool
	Value    string
}

/*
 * Print color console message
 *
 * flag   : include 'Waring', 'Error'
 * message: print content
 * args   : variable parameter, include string, CC type
 *          when args last value is "\n", auto new line.
 *
 * e.g. P( "Waring", Remote latest version [%v] = latest version [%v].\n", param1, param2 )
 * e.g. cc := CC{1, true, 2, true, localVersion}
 *      P(DEFAULT, "Current version %v, publish data: ", cc, "2014-05-31")
 * e.g. P(DEFAULT, "Current version %v", localVersion, "\n")
 *
 */
func P(flag string, message interface{}, args ...interface{}) {

	// try catch
	defer func() {
		if err := recover(); err != nil {
			Error(ERROR, "util/print.go an error has occurred. Error: ", err)
			os.Exit(0)
		}
	}()

	// set state
	stateColor(flag)

	// set color message
	msgArr := strings.Split(message.(string), SPLIT)
	for k, v := range msgArr {
		fmt.Print(v)
		if k < len(args) {
			t := reflect.TypeOf(args[k])
			switch t.Name() {
			case "string":
				normalColor(args[k])
			case "CC":
				customColor(args[k])
			default:
				normalColor(args[k])
			}
		}
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

func customColor(cc interface{}) {

	value := reflect.ValueOf(cc)

	fgColor := value.FieldByName("FgColor").Int()
	fgBright := value.FieldByName("FgBright").Bool()
	bgColor := value.FieldByName("BgColor").Int()
	bgBright := value.FieldByName("BgBright").Bool()
	msg := value.FieldByName("Value").String()

	if fgColor > 8 || fgColor < 0 || bgColor > 8 || bgColor < 0 {
		normalColor(msg)
		fmt.Println()
		Error(WARING, "values range error, values range include 0 ~ 8, Error: ", "index out of range")
		return
	}

	ct.ChangeColor(ct.Color(fgColor), fgBright, ct.Color(bgColor), bgBright)
	fmt.Printf(msg)
	ct.ResetColor()
}

func normalColor(msg interface{}) {
	ct.ChangeColor(ct.Green, true, ct.None, false)
	fmt.Printf(msg.(string))
	ct.ResetColor()
}
