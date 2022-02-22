package httpgo

import (
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

func TestGetHeader(t *testing.T) {
	t.Log(Get("https://qq.com").Do())
}

func TestPost(t *testing.T) {
	t.Log(Post("127.0.0.1/post").SetForm(map[string]string{"name": "peter", "address": "unknown"}).Do())
}

func TestName(t *testing.T) {
	r := Post("https://qq.com").SetUA("test-ua").
		AddCookie(&http.Cookie{
			Name:  "http_go_cookie",
			Value: "http_go",
		}).
		AddHeader("test", "value").Do()
	t.Log(r)
}
