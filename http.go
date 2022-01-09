package httpgo

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// HTTP methods support
const (
	POST    = "POST"
	GET     = "GET"
	HEAD    = "HEAD"
	PUT     = "PUT"
	DELETE  = "DELETE"
	PATCH   = "PATCH"
	OPTIONS = "OPTIONS"
)

// content type support
const (
	JSON            = "application/json;charset=UTF-8"
	FORM_URLENCODED = "application/x-www-form-urlencoded;charset=UTF-8"
	FORM_DATA       = "multipart/form-data;charset=UTF-8"
)

type Request struct {
	*http.Request
	Err      error
	client   *Client
	callback func(*Response) *Response
}

type Response struct {
	*http.Response
	Req  *Request
	Body []byte
	Err  error
}

// Get http request
func Get(rawUrl string) *Request {
	return NewRequest(GET, rawUrl)
}

// Post http request
func Post(rawUrl string) *Request {
	return NewRequest(POST, rawUrl)
}

// NewRequest 构建请求体
func NewRequest(method, rawUrl string) *Request {
	req, err := http.NewRequest(method, parseScheme(rawUrl), nil)
	return &Request{
		Request: req,
		Err:     err,
		client:  DefaultClient,
		callback: func(resp *Response) *Response {
			return resp
		},
	}
}

// Do 最终执行请求
func (r *Request) Do() *Response {
	return r.callback(r.client.Do(r))
}

// SetQuery 设置GET请求参数
func (r *Request) SetQuery(m map[string]string) *Request {
	if len(r.URL.RawQuery) > 0 {
		r.URL.RawQuery += "&"
	}
	for k, v := range m {
		r.URL.RawQuery += k + "=" + v
	}
	return r
}

// SetForm 设置POST表单参数
func (r *Request) SetForm(m map[string]string) *Request {
	var payload = url.Values{}
	for k, v := range m {
		payload.Add(k, v)
	}
	r.setBody(strings.NewReader(payload.Encode()))
	r.Header.Set("Content-Type", FORM_URLENCODED)
	return r
}

// SetJSON 设置JSON参数
func (r *Request) SetJSON(v string) *Request {
	r.setBody(strings.NewReader(v))
	r.Header.Set("Content-Type", JSON)
	return r
}

func (r *Request)SetTimeout(t time.Duration) *Request {
	ctx, _ := context.WithTimeout(r.Context(), t)
	r.Request = r.WithContext(ctx)
	return r
}

// 设置body
func (r *Request) setBody(body io.Reader) {
	rc, ok := body.(io.ReadCloser)
	if !ok && body != nil {
		rc = io.NopCloser(body)
	}
	r.Body = rc

	switch v := body.(type) {
	case *bytes.Buffer:
		r.ContentLength = int64(v.Len())
		buf := v.Bytes()
		r.GetBody = func() (io.ReadCloser, error) {
			r := bytes.NewReader(buf)
			return io.NopCloser(r), nil
		}
	case *bytes.Reader:
		r.ContentLength = int64(v.Len())
		snapshot := *v
		r.GetBody = func() (io.ReadCloser, error) {
			r := snapshot
			return io.NopCloser(&r), nil
		}
	case *strings.Reader:
		r.ContentLength = int64(v.Len())
		snapshot := *v
		r.GetBody = func() (io.ReadCloser, error) {
			r := snapshot
			return io.NopCloser(&r), nil
		}
	default:
	}
	if r.GetBody != nil && r.ContentLength == 0 {
		r.Body = http.NoBody
		r.GetBody = func() (io.ReadCloser, error) { return http.NoBody, nil }
	}

}

// 解析协议头
func parseScheme(url string) string {
	if strings.HasPrefix(url, "https://") || strings.HasPrefix(url, "http://") {
		return url
	}
	if strings.HasPrefix(url, ":") {
		return fmt.Sprintf("http://localhost%s", url)
	}
	if strings.HasPrefix(url, "/") {
		return fmt.Sprintf("http://localhost%s", url)
	}
	return fmt.Sprintf("http://%s", url)
}
