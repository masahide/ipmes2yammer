package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
	//"github.com/ActiveState/tail"
	"strconv"

	"github.com/kr/pretty"
	"github.com/masahide/go-yammer/oauth"
)

var (
	cachefile = "cache.json"

	scope             = "https://www.yammer.com/"
	request_token_url = "https://www.yammer.com/dialog/oauth"
	auth_token_url    = "https://www.yammer.com/oauth2/access_token.json"

	//clientId     =
	//clientSecret =
	//
	redirectURL = "http://localhost"
	//
	requestURL = "https://www.yammer.com/api/v1/messages.json"

	by_emailURL = "https://www.yammer.com/api/v1/users/by_email.json" // ?email=user@domain.com
	postURL     = "https://www.yammer.com/api/v1/messages.json"       // body replied_to_id
	//mail        =
)

func main() {
	runtime.GOMAXPROCS(2)

	var err error
	var lsConfig LocalServerConfig

	flag.IntVar(&lsConfig.Port, "p", 16061, "local port: 1024 < ")
	flag.IntVar(&lsConfig.Timeout, "t", 30, "redirect timeout: 0 - 90")

	flag.Parse()

	if lsConfig.Port <= 1024 ||
		lsConfig.Timeout < 0 || lsConfig.Timeout > 90 {
		fmt.Fprintf(os.Stderr, "Usage: \n")
		flag.PrintDefaults()
		os.Exit(2)
	}

	fmt.Println("Start Execute API")

	// 認証コードを引数で受け取る。
	//code := flag.Arg(0)

	config := &oauth.Config{
		ClientId:     clientId,
		ClientSecret: clientSecret,
		RedirectURL:  fmt.Sprintf("%s:%d", redirectURL, lsConfig.Port),
		Scope:        scope,
		AuthURL:      request_token_url,
		TokenURL:     auth_token_url,
		TokenCache:   oauth.CacheFile(cachefile),
	}

	yammer, err := yammerOauth(config, lsConfig)
	if err != nil {
		log.Fatal("Yammer Oauth Error:", err)
	}

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

func sendYammer(yammer *oauth.Transport, id int, message string) (err error) {

	//	func (c *Client) Post(url string, bodyType string, body io.Reader) (resp *Response, err error)

	fmt.Println("\nstart Post")
	r, err := yammer.Client().PostForm(postURL, url.Values{
		"direct_to_id": {strconv.Itoa(id)},
		"body":         {message},
	})
	fmt.Println("\nend Post")
	if err != nil {
		log.Fatal("Get:", err)
		return
	}
	defer r.Body.Close()
	if r.StatusCode != 200 {
		pretty.Printf("--- id:\n%# v\nstatus: %# v\n", id, r.StatusCode)
		//pretty.Printf("sendMessage Response--------\n %# v\n", r)
		fmt.Println()
		fmt.Println("hoge[")
		io.Copy(os.Stdout, r.Body)
		fmt.Println("]fuga")
		fmt.Println()
		return fmt.Errorf("sendMessage %v", r.Status)
	}
	pretty.Printf("--- id:\n%# v\nstatus: %# v\n", id, r.StatusCode)
	pretty.Printf("Response --------\n %# v\n", r)
	fmt.Println("hoge[")
	io.Copy(os.Stdout, r.Body)
	fmt.Println("]fuga")
	return
}

func emailToIDYammer(yammer *oauth.Transport, email string) (id int, err error) {
	r, in_err := yammer.Client().Get(by_emailURL + "?email=" + mail)
	if in_err != nil {
		log.Fatal("Get:", in_err)
		return 0, in_err
	}
	defer r.Body.Close()
	if r.StatusCode != 200 {
		log.Fatal("updateToken %v", r.Status)
		return
	}
	var data interface{}
	if err = json.NewDecoder(r.Body).Decode(&data); err != nil {
		return
	}
	//pretty.Printf("--- data:\n%# v\n\n", data)
	id = int(data.([]interface{})[0].(map[string]interface{})["id"].(float64))
	pretty.Printf("--- id:\n%# v\n\n", id)
	return
}

func yammerOauth(config *oauth.Config, lsConfig LocalServerConfig) (transport *oauth.Transport, err error) {

	transport = &oauth.Transport{Config: config}

	// キャッシュからトークンファイルを取得
	_, err = config.TokenCache.Token()
	if err != nil {
		code, err := getAuthCode(config, lsConfig)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}

		// 認証トークンを取得する。（取得後、キャッシュへ）
		_, err = transport.Exchange(code)
		if err != nil {
			fmt.Printf("Exchange Error: %v\n", err)
			os.Exit(1)
		}
	}
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

func getAuthCode(config *oauth.Config, lsConfig LocalServerConfig) (string, error) {
	url := config.AuthCodeURL("")

	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		url = strings.Replace(url, "&", `^&`, -1)
		cmd = exec.Command("cmd", "/c", "start", url)

	case "darwin":
		//url = strings.Replace(url, "&", `\&`, -1)
		cmd = exec.Command("open", url)

	default:
		return "", fmt.Errorf("ブラウザで以下のURLにアクセスし、認証して下さい。\n%s\n", url)
	}

	redirectResult := make(chan RedirectResult, 1)
	serverStarted := make(chan bool, 1)
	//
	go func(rr chan<- RedirectResult, ss chan<- bool, p int) {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			code := r.URL.Query().Get("code")

			if code == "" {
				rr <- RedirectResult{Err: fmt.Errorf("codeを取得できませんでした。")}
			}

			fmt.Fprintf(w, `<!doctype html>
<html lang="ja">
<head>
<meta charset="utf-8">
</head>
<body onload="window.open('about:blank','_self').close();">
ブラウザが自動で閉じない場合は手動で閉じてください。
</body>
</html>
`)
			rr <- RedirectResult{Code: code}
		})

		host := fmt.Sprintf("localhost:%d", p)

		fmt.Printf("Start Listen: %s\n", host)
		ss <- true

		err := http.ListenAndServe(host, nil)

		if err != nil {
			rr <- RedirectResult{Err: err}
		}
	}(redirectResult, serverStarted, lsConfig.Port)

	<-serverStarted

	// set redirect timeout
	tch := time.After(time.Duration(lsConfig.Timeout) * time.Second)

	fmt.Println("Start your browser after 2sec.")

	time.Sleep(2 * time.Second)

	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("Browser Start Error: %v\n", err)
	}

	var rr RedirectResult

	select {
	case rr = <-redirectResult:
	case <-tch:
		return "", fmt.Errorf("Timeout: waiting redirect.")
	}

	if rr.Err != nil {
		return "", fmt.Errorf("Redirect Error: %v\n", rr.Err)
	}

	fmt.Printf("Got code.\n")

	return rr.Code, nil
}

type RedirectResult struct {
	Code string
	Err  error
}

type LocalServerConfig struct {
	Port    int
	Timeout int
}
