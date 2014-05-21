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
 * name: download file name e.g. node.exe
 * dst: downlaod path
 *
 * return code
 * 0: success
 * -1: status code != 200
 * -2: create folder error
 * -3: download node.exe error
 *
 */
func New(url, name, dst string) int {

	// get res
	res, err := http.Get(url)

	// close
	defer res.Body.Close()

	// err
	if err != nil {
		panic(err)
	}

	// check state code
	if res.StatusCode != 200 {
		fmt.Printf("Downlaod url [%v] an [%v] error occurred, please check.", url, res.StatusCode)
		return -1
	}

	// create file
	file, createErr := os.Create(dst)
	if createErr != nil {
		fmt.Println("Create file error, Error: " + createErr.Error())
		return -2
	}
	defer file.Close()

	fmt.Printf("Start download [%v] from %v.\n", name, url)

	// loop buff to file
	buf := make([]byte, res.ContentLength)
	var m float32
	isShow, oldCurrent := false, 0
	for {
		n, err := res.Body.Read(buf)

		// write complete
		if n == 0 {
			fmt.Println("100% \nEnd download.")
			break
		}

		//error
		if err != nil {
			panic(err)
		}

		/* show console e.g.
		 * Start download node.exe version [x.xx.xx] from http://nodejs.org/dist/.
		 * 10% 20% 30% 40% 50% 60% 70% 80% 90% 100%
		 * End download.
		 */
		m = m + float32(n)
		current := int(m / float32(res.ContentLength) * 100)

		if current > oldCurrent {
			switch current {
			case 10, 20, 30, 40, 50, 60, 70, 80, 90:
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
			fmt.Printf("Error: Downlaod [%v] size error, please check your network and run 'gnvm uninstall %v'.\n", name, name)
			return -3
		}
	}

	return 0
}
