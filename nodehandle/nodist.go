package nodehandle

import (

	// lib
	"github.com/bitly/go-simplejson"
	//"github.com/Kenshin/curl"
	"curl"

	// go
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"

	// local
	"gnvm/util"
)

type (
	Node struct {
		Version string
		Exec    string
	}

	NPM struct {
		Version string
	}

	NodeDetail struct {
		ID   int
		Date string
		Node
		NPM
	}

	Nodeist map[string]NodeDetail
)

var sorts []string

/*
 Create nl( map[string]NodeDetail )

 Param:
    - url:    index.json url, e.g. http://npm.taobao.org/mirrors/node/index.json
    - filter: regexp when regexp == nil, filter all NodeDetail

 Return:
    - nl:     nodedetail collection
    - error:  error
    - code:   error flag

      Code:
        - -1: get url error
        - -2: read res.body error
        - -3: create json error
        - -4: parse json error
*/
func New(url string, filter *regexp.Regexp) (*Nodeist, error, int) {
	code, res, err := curl.Get(url)
	if err != nil {
		return nil, err, code
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err, -2
	}

	json, err := simplejson.NewJson(body)
	if err != nil {
		return nil, err, -3
	}
	arr, err := json.Array()
	if err != nil {
		return nil, err, -4
	}

	nl, idx := make(Nodeist, 0), 0
	sorts = make([]string, 0)
	for _, element := range arr {
		if value, ok := element.(map[string]interface{}); ok {
			ver, _ := value["version"].(string)
			if filter != nil {
				if ok := filter.MatchString(ver[1:]); !ok {
					continue
				}
			}
			date, _ := value["date"].(string)
			npm, _ := value["npm"].(string)
			if npm == "" {
				npm = "[x]"
			}
			exe := parse(ver[1:])
			sorts = append(sorts, ver)
			nl[ver] = NodeDetail{idx, date, Node{ver, exe}, NPM{npm}}
			idx++
		}
	}
	return &nl, nil, 0
}

/*
 Print NodeDetail collection

 Param:
    - limit: print lines, when limit == 0, print all nodedetail

*/
func (this *Nodeist) Detail(limit int) {
	table := `+--------------------------------------------------+
| No.   date         node ver    exec      npm ver |
+--------------------------------------------------+`
	if limit == 0 || limit > len(sorts) {
		limit = len(sorts)
	}
	for idx, v := range sorts {
		if idx == 0 {
			fmt.Println(table)
		}
		if idx >= limit {
			break
		}
		value := (*this)[v]
		id := format(strconv.Itoa(value.ID+1), 6)
		date := format(value.Date, 13)
		ver := format(value.Node.Version[1:], 12)
		exe := format(value.Node.Exec, 10)
		npm := format(value.NPM.Version, 9)
		fmt.Println("  " + id + date + ver + exe + npm)
		if idx == limit-1 {
			fmt.Println("+--------------------------------------------------+")
		}
	}
}

func parse(version string) (exec string) {
	switch util.GetNodeVerLev(util.FormatNodeVer(version)) {
	case 0:
		exec = "[x]"
	case 1:
		exec = "x86"
	default:
		exec = "x86 x64"
	}
	return
}

func format(value string, max int) string {
	if len(value) > max {
		max = len(value)
	}
	newValue := strings.Repeat(" ", max-len(value))
	return value + newValue
}
