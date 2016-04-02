package nodehandle

import (

	// lib
	"github.com/Kenshin/curl"
	"github.com/bitly/go-simplejson"

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

	Nodist struct {
		nl    map[string]NodeDetail
		Sorts []string
	}
)

/*
 Create nodist( map[string]NodeDetail )

 Param:
    - url:    index.json url, e.g. http://npm.taobao.org/mirrors/node/index.json
    - filter: regexp when regexp == nil, filter all NodeDetail

 Return:
    - nodist: nodedetail collection
    - error:  error
    - code:   error flag

      Code:
        - -1: get url error
        - -2: read res.body error
        - -3: create json error
        - -4: parse json error

*/
func New(url string, filter *regexp.Regexp) (*Nodist, error, int) {
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

	nodist, idx := new(Nodist), 0
	nodist.nl = make(map[string]NodeDetail, 0)
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
			exe := formatExe(ver[1:])
			nodist.Sorts = append(nodist.Sorts, ver)
			nodist.nl[ver] = NodeDetail{idx, date, Node{ver, exe}, NPM{npm}}
			idx++
		}
	}
	return nodist, nil, 0
}

/*
 Find NodeDetail by node version

 Param:
    - url: index.json url, e.g. http://npm.taobao.org/mirrors/node/index.json
    - ver: node version. e.g. 5.9.0

 Return:
    - *NodeDetail: nodedetail struct
    - error

*/
func FindNodeDetailByVer(url, ver string) (*NodeDetail, error) {
	filter, err := util.FormatWildcard(ver, url)
	if err != nil {
		return nil, err
	}
	nodist, err, _ := New(url, filter)
	if err != nil {
		return nil, err
	}
	if len(nodist.nl) == 1 {
		nd := nodist.nl["v"+ver]
		return &nd, nil
	}
	return nil, nil
}

/*
 Print NodeDetail collection

 Param:
    - limit: print lines, when limit == 0, print all nodedetail

*/
func (this *Nodist) Detail(limit int) {
	table := `+--------------------------------------------------+
| No.   date         node ver    exec      npm ver |
+--------------------------------------------------+`
	if limit == 0 || limit > len(this.Sorts) {
		limit = len(this.Sorts)
	}
	for idx, v := range this.Sorts {
		if idx == 0 {
			fmt.Println(table)
		}
		if idx >= limit {
			break
		}
		value := this.nl[v]
		id := leftpad(strconv.Itoa(value.ID+1), 6)
		date := leftpad(value.Date, 13)
		ver := leftpad(value.Node.Version[1:], 12)
		exe := leftpad(value.Node.Exec, 10)
		npm := leftpad(value.NPM.Version, 9)
		fmt.Println("  " + id + date + ver + exe + npm)
		if idx == limit-1 {
			fmt.Println("+--------------------------------------------------+")
		}
	}
}

/*
 Format exe

 Param:
 	- version: node.exe version

 Return:
 	- exec: formatting string, e.g. '[x]'

*/
func formatExe(version string) (exec string) {
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

/*
  Format label, e.g.
     aa:
    bbb:

 Param:
 	- value: format str
 	- max  : max empty

 Return:
 	- Format label

*/
func leftpad(value string, max int) string {
	if len(value) > max {
		max = len(value)
	}
	newValue := strings.Repeat(" ", max-len(value))
	return value + newValue
}
