package main

import (
	"flag"
	"log"

	"runtime"
	"time"

//"github.com/ActiveState/tail"
)

func main() {
	runtime.GOMAXPROCS(2)

	var err error
	var lsConfig LocalServerConfig

	flag.IntVar(&lsConfig.Port, "p", 16061, "local port: 1024 < ")
	flag.IntVar(&lsConfig.Timeout, "t", 30, "redirect timeout: 0 - 90")

	flag.Parse()

	/*
		//
		// Make the request.
		r, err := yammer.Client().Get(requestURL)
		if err != nil {
			log.Fatal("Get:", err)
		}
		defer r.Body.Close()

		// Write the response to standard output.
		io.Copy(os.Stdout, r.Body)

		// Send final carriage return, just to be neat.
		fmt.Println()
	*/

	id, err := emailToIDYammer(yammer, mail)
	if err != nil {
		log.Fatal("emailtoID:", err)
		return
	}
	time.Sleep(2 * time.Second)
	sendYammer(yammer, id, "テスト")
}

func getNewIPMessage() (message string, err error) {
	message = "hoge"
	err = nil
	return
}

/*
func Decode(r *http.Response, in_err error) (data interface{}, err error) {
	err = in_err
	if err != nil {
		err = json.NewDecoder(r.Body).Decode(&data)
	}
	return
		switch t := x.AuthorRaw.(type) {
		case string:
			x.AuthorEmail = t
		case json.Number:
			var n uint64
			// We would shadow the outer `err` here by using `:=`
			n, err = t.Int64()
			x.AuthorID = n
		}
		return
}
*/
