package main

import (
	"flag"
	"log"

	"runtime"
	"time"

	//"github.com/ActiveState/tail"
	"github.com/masahide/go-yammer/yammer"
)

func main() {
	runtime.GOMAXPROCS(2)

	var err error
	var lsConfig yammer.LocalServerConfig
	var mail string

	flag.IntVar(&lsConfig.Port, "p", 16061, "local port: 1024 < ")
	flag.IntVar(&lsConfig.Timeout, "t", 30, "redirect timeout: 0 - 90")
	flag.StringVar(&mail, "m", "", "email addresss")

	flag.Parse()
	if mail == "" {
		flag.PrintDefaults()
		return
	}

	y := yammer.NewYammer(&lsConfig)
	err = y.YammerAuth()
	if err != nil {
		log.Fatal("Error YammerAuth:", err)
		return
	}
	id, err := y.EmailToIDYammer(mail)
	if err != nil {
		log.Fatal("Erorr emailtoID:", err)
		return
	}

	time.Sleep(2 * time.Second)
	y.SendYammer(id, "テスト")
}

func getNewIPMessage() (message string, err error) {
	message = "hoge"
	err = nil
	return
}
