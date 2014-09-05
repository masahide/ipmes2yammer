package yammer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
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
	id = int(data.([]interface{})[0].(map[string]interface{})["id"].(float64))
	fmt.Printf("--- id:\n%# v\n\n", id)
	return
}

func (y *Yammer) Send(method string, id int, message string) (string, error) {

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

func (y *Yammer) Unfollow(id string) (string, error) {

	req, err := http.NewRequest("DELETE", "https://www.yammer.com/api/v1/threads/"+id+"/follow.json", nil)
	if err != nil {
		log.Fatal("NewRequest:", err)
		return "", err
	}
	r, err := y.transport.Client().Do(req)
	if err != nil {
		log.Fatal("Get:", err)
		return "", err
	}
	defer r.Body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	if r.StatusCode != 200 {
		return buf.String(), fmt.Errorf("Unsubscribe Code:%d, Status:%v", r.StatusCode, r.Status)
	}
	return buf.String(), nil
}
