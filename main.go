package main

import (
	"flag"
	"log"
	"strings"
	"time"

	"runtime"

	//"github.com/ActiveState/tail"
	"github.com/masahide/ipmes2yammer/ipmes"
	"github.com/masahide/ipmes2yammer/yammer"
)

func main() {
	runtime.GOMAXPROCS(2)

	var err error
	var lsConfig yammer.LocalServerConfig
	var mail, file, unfollow string
	var id int

	flag.IntVar(&lsConfig.Port, "p", 16061, "local port: 1024 < ")
	flag.IntVar(&lsConfig.Timeout, "t", 30, "redirect timeout: 0 - 90")
	//flag.StringVar(&mail, "to_mail", "", "email addresss")
	flag.IntVar(&id, "to_id", 0, "threadId")
	flag.StringVar(&file, "file", "", "tail file name")
	flag.StringVar(&unfollow, "unfollow", "", "auto unfollow threadid(csv)")

	flag.Parse()
	if (mail == "" && id == 0) || file == "" {
		flag.PrintDefaults()
		return
	}

	y := yammer.NewYammer(&lsConfig)
	err = y.YammerAuth()
	if err != nil {
		log.Fatal("Error YammerAuth:", err)
		return
	}
	t, err := ipmes.NewTailFile(file)
	if err != nil {
		log.Fatal("import.NewTailFile:", err)
	}
	if unfollow != "" {
		go func() {
			for {
				for _, uf := range strings.Split(unfollow, ",") {
					_, err := y.Unfollow(uf)
					if err != nil {
						log.Printf("Error Unfollow:%v, err:%v", uf, err)
					}
					time.Sleep(60 * time.Second)
				}
			}
		}()
	}
	for {

		time.Sleep(1 * time.Second)
		s, err := t.TailMessage()
		if err != nil {
			log.Fatal("Error ipmes.TailMessage:", err)
		}
		if s != "" {
			//fmt.Println(s)
			y.Send("replied_to_id", id, s)
		}
	}

}
