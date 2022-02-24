package httpgo

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	t.Log(Get("https://qq.com").Do())
}

func TestGetTimeout(t *testing.T) {
	t.Log(Get("https://voocel.com").SetTimeout(1 * time.Second).Do())
}

func TestPost(t *testing.T) {
	t.Log(Post("127.0.0.1/post").SetForm(map[string]string{"name": "peter", "address": "unknown"}).Do())
}

func TestHeader(t *testing.T) {
	r := Post("https://qq.com").SetUA("test-ua").
		AddCookie(&http.Cookie{
			Name:  "http_go_cookie",
			Value: "http_go",
		}).
		AddHeader("test", "value").Do()
	t.Log(r)
}

func TestAsync(t *testing.T) {
	ch := make(chan *AsyncResponse)
	AsyncGet("https://qq.com", ch)

	timer := time.NewTimer(5*time.Second)
	for {
		select {
		case res := <-ch:
			fmt.Println(res.Resp)
		case <-timer.C:
			return
		}
	}
}

func TestRequest_SetFile(t *testing.T) {
	Post("http://127.0.0.1:3333/file").SetFile("file", "a.txt").Do()
}