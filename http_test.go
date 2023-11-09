package httpgo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
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
	t.Log(Post("127.0.0.1:3333/post").SetForm(map[string]string{"name": "peter", "address": "unknown"}).Do())
}

func TestPostJSON(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			m := map[string]string{
				"name": fmt.Sprint("tony-", i),
				"addr": "unknown",
			}
			b, _ := json.Marshal(m)
			t.Log(Post("127.0.0.1:3333/post").SetJSON(string(b)).Do())
		}(i)
	}
	wg.Wait()
}

func TestHeader(t *testing.T) {
	r, err := Post("https://qq.com").SetUA("test-ua").
		AddCookie(&http.Cookie{
			Name:  "http_go_cookie",
			Value: "http_go",
		}).
		AddHeader("test", "value").Do()
	t.Log(r, err)
}

func TestAsync(t *testing.T) {
	ch := make(chan *AsyncResponse)
	AsyncGet("https://qq.com", ch)

	timer := time.NewTimer(5 * time.Second)
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
	r, err := Post("http://127.0.0.1:3333/file").SetFile("file", "abc.txt").Do()
	t.Log(r, err)
}
