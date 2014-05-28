package curl

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
)

type ProcessFunc func(content string, line int)

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

	// err
	if err != nil {
		panic(err)
	}

	// check state code
	if res.StatusCode != 200 {
		fmt.Printf("URL [%v] an [%v] error occurred, please check.\n", url, res.StatusCode)
		return -1, res, err
	}

	return 0, res, err

}

/*
 *
 * Read line from io.ReadCloser
 *
 * return err
 *
 */
func ReadLine(body io.ReadCloser, process ProcessFunc) error {

	var content string
	var err error
	var line int = 1

	// set buff
	buff := bufio.NewReader(body)

	for {
		content, err = buff.ReadString('\n')

		if line > 1 && (err != nil || err == io.EOF) {
			break
		}

		process(content, line)

		line++
	}

	return err
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
 * -2: create file error
 * -3: download node.exe error
 * -4: content length = -1
 *
 */
func New(url, name, dst string) int {

	// try catch
	defer func() {
		if err := recover(); err != nil {
			msg := fmt.Sprintf("CURL Error: Download %v from %v an error has occurred. \nError: %v", name, url, err)
			panic(msg)
		}
	}()

	// get url
	code, res, err := Get(url)
	if code != 0 {
		return code
	}

	// close
	defer res.Body.Close()

	// create file
	file, createErr := os.Create(dst)
	if createErr != nil {
		fmt.Println("Create file error, Error: " + createErr.Error())
		return -2
	}
	defer file.Close()

	if res.ContentLength == -1 {
		fmt.Printf("Download %v fail from %v.\n", name, url)
		return -4
	}

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
