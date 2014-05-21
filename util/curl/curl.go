package curl

import (
	"fmt"
	"net/http"
	"os"
)

/*
 *
 * parameter
 * url: download url e.g. http://nodejs.org/dist/v0.10.0/node.exe
 *
 * return code
 * 0: success
 * -1: status code != 200
 *
 * return res
 * return err
 *
 */
func Get(url string) (code int, res *http.Response, err error) {

	// get res
	res, err = http.Get(url)

	// close
	defer res.Body.Close()

	// err
	if err != nil {
		panic(err)
	}

	// check state code
	if res.StatusCode != 200 {
		fmt.Printf("URL [%v] an [%v] error occurred, please check.", url, res.StatusCode)
		return -1, res, err
	}

	return 0, res, err

}

/*
 *
 * parameter
 * url: download url e.g. http://nodejs.org/dist/v0.10.0/node.exe
 * name: download file name e.g. node.exe
 * dst: download path
 *
 * return code
 * 0: success
 * -1: status code != 200 ( from Get() method )
 * -2: create folder error
 * -3: download node.exe error
 *
 */
func New(url, name, dst string) int {

	// get url
	code, res, err := Get(url)
	if code != 0 {
		return code
	}

	// create file
	file, createErr := os.Create(dst)
	if createErr != nil {
		fmt.Println("Create file error, Error: " + createErr.Error())
		return -2
	}
	defer file.Close()

	fmt.Printf("Start download [%v] from %v.\n%v", name, url, "1% ")

	// loop buff to file
	buf := make([]byte, res.ContentLength)
	var m float32
	isShow, oldCurrent := false, 0
	for {
		n, err := res.Body.Read(buf)

		// write complete
		if n == 0 && err.Error() == "EOF" {
			fmt.Println("100% \nEnd download.")
			break
		}

		//error
		if err != nil {
			panic(err)
		}

		/* show console e.g.
		 * Start download [x.xx.xx] from http://nodejs.org/dist/.
		 * 1% 5% 10% 20% 30% 40% 50% 60% 70% 80% 90% 100%
		 * End download.
		 */
		m = m + float32(n)
		current := int(m / float32(res.ContentLength) * 100)

		switch {
		case current > 0 && current < 6:
			current = 5
		case current > 5 && current < 11:
			current = 10
		case current > 10 && current < 21:
			current = 20
		case current > 20 && current < 31:
			current = 30
		case current > 30 && current < 41:
			current = 40
		case current > 40 && current < 51:
			current = 50
		case current > 60 && current < 71:
			current = 60
		case current > 70 && current < 81:
			current = 70
		case current > 80 && current < 91:
			current = 80
		case current > 90 && current < 101:
			current = 90
		}

		if current > oldCurrent {
			switch current {
			case 5, 10, 20, 30, 40, 50, 60, 70, 80, 90:
				isShow = true
			}

			if isShow {
				fmt.Printf("%d%v", current, "% ")
			}

			isShow = false
		}

		oldCurrent = current

		file.WriteString(string(buf[:n]))
	}

	// valid download exe
	fi, err := file.Stat()
	if err == nil {
		if fi.Size() != res.ContentLength {
			fmt.Printf("Error: Downlaod [%v] size error, please check your network.\n", name)
			return -3
		}
	}

	return 0
}
