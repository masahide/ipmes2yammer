package yammer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
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

// method	direct_to_id, replied_to_id
func (y *Yammer) Send(method string, id int, message string) (string, error) {

	//	func (c *Client) Post(url string, bodyType string, body io.Reader) (resp *Response, err error)

	r, err := y.transport.Client().PostForm(postURL, url.Values{
		method: {strconv.Itoa(id)},
		"body": {message},
	})
	if err != nil {
		log.Fatal("Get:", err)
		return "", err
	}
	defer r.Body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	if r.StatusCode != 200 {
		return buf.String(), fmt.Errorf("sendMessage Code:%d, Status:%v", r.StatusCode, r.Status)
	}
	return buf.String(), nil
}
