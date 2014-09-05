package main

import (
	"flag"
	"log"

	"runtime"

	//"github.com/ActiveState/tail"
	"github.com/masahide/go-yammer/yammer"
)

func main() {
	runtime.GOMAXPROCS(2)

	var err error
	var lsConfig yammer.LocalServerConfig
	var mail string
	var id int

	flag.IntVar(&lsConfig.Port, "p", 16061, "local port: 1024 < ")
	flag.IntVar(&lsConfig.Timeout, "t", 30, "redirect timeout: 0 - 90")
	flag.StringVar(&mail, "m", "", "email addresss")
	flag.IntVar(&id, "id", 0, "threadId")

	flag.Parse()
	if mail == "" && id == 0 {
		flag.PrintDefaults()
		return
	}

	y := yammer.NewYammer(&lsConfig)
	err = y.YammerAuth()
	if err != nil {
		log.Fatal("Error YammerAuth:", err)
		return
	}
	if mail != "" {
		id, err := y.EmailToIDYammer(mail)
		if err != nil {
			log.Fatal("Erorr emailtoID:", err)
			return
		}
		y.Send("direct_to_id", id, "テスト")
	} else if id != 0 {
		y.Send("replied_to_id", id, "テスト")
	}

}

func getNewIPMessage() (message string, err error) {
	message = "hoge"
	err = nil
	return
}
