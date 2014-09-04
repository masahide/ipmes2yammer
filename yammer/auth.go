package yammer

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/masahide/go-yammer/oauth"
)

type LocalServerConfig struct {
	Port    int
	Timeout int
}

type RedirectResult struct {
	Code string
	Err  error
}

type Yammer struct {
	transport *oauth.Transport
	config    *oauth.Config
	lsConfig  *LocalServerConfig
}

const (
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

func (y *Yammer) YammerAuth() {
	runtime.GOMAXPROCS(2)

	y.config = &oauth.Config{
		ClientId:     clientId,
		ClientSecret: clientSecret,
		RedirectURL:  fmt.Sprintf("%s:%d", redirectURL, y.lsConfig.Port),
		Scope:        scope,
		AuthURL:      request_token_url,
		TokenURL:     auth_token_url,
		TokenCache:   oauth.CacheFile(cachefile),
	}

	transport, err := y.yammerOauth()
	if err != nil {
		log.Fatal("Yammer Oauth Error:", err)
	}
	y.transport = transport

}

func (y *Yammer) yammerOauth() (transport *oauth.Transport, err error) {

	transport = &oauth.Transport{Config: y.config}

	// キャッシュからトークンファイルを取得
	_, err = y.config.TokenCache.Token()
	if err != nil {
		code, err := y.getAuthCode()
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

func (y *Yammer) getAuthCode() (string, error) {
	url := y.config.AuthCodeURL("")

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
	}(redirectResult, serverStarted, y.lsConfig.Port)

	<-serverStarted

	// set redirect timeout
	tch := time.After(time.Duration(y.lsConfig.Timeout) * time.Second)

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
