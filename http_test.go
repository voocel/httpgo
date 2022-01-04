package httpgo

import "testing"

func TestGet(t *testing.T) {
	t.Log(Get("https://qq.com").Do())
}

func TestPost(t *testing.T) {
	t.Log(Post("127.0.0.1/post").SetForm(map[string]string{"name": "peter", "address": "unknown"}).Do())
}
