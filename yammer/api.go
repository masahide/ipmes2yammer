package yammer

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"strconv"

	"github.com/kr/pretty"
)

func (y *Yammer) EmailToIDYammer(email string) (id int, err error) {
	r, in_err := y.transport.Client().Get(by_emailURL + "?email=" + email)
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

func (y *Yammer) SendYammer(id int, message string) (err error) {

	//	func (c *Client) Post(url string, bodyType string, body io.Reader) (resp *Response, err error)

	fmt.Println("\nstart Post")
	r, err := y.transport.Client().PostForm(postURL, url.Values{
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
