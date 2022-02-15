package httpgo

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
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

// NewRequest create request
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

// Do finish do
func (r *Request) Do() *Response {
	return r.callback(r.client.do(r))
}

// SetQueries set URL query params for the request
func (r *Request) SetQueries(m map[string]string) *Request {
	for k, v := range m {
		r.SetQuery(k, v)
	}
	return r
}

// SetQuery set an URL query parameter for the request
func (r *Request) SetQuery(key, value string) *Request {
	if len(r.URL.RawQuery) > 0 {
		r.URL.RawQuery += "&"
	}
	r.URL.RawQuery += key + "=" + value
	return r
}

// SetForm set the form data from a map
func (r *Request) SetForm(m map[string]string) *Request {
	var payload = url.Values{}
	for k, v := range m {
		payload.Add(k, v)
	}
	r.setBody(strings.NewReader(payload.Encode()))
	r.SetHeader("Content-Type", FORM_URLENCODED)
	return r
}

// SetJSON set the json data
func (r *Request) SetJSON(v string) *Request {
	r.setBody(strings.NewReader(v))
	r.SetHeader("Content-Type", JSON)
	return r
}

func (r *Request) SetTimeout(t time.Duration) *Request {
	ctx, _ := context.WithTimeout(r.Context(), t)
	r.Request = r.WithContext(ctx)
	return r
}

// setBody set the request body, accepts string, []byte, io.Reader, map and struct
func (r *Request) setBody(body io.Reader) {
	rc, ok := body.(io.ReadCloser)
	if !ok && body != nil {
		rc = ioutil.NopCloser(body)
	}
	r.Body = rc

	switch v := body.(type) {
	case *bytes.Buffer:
		r.ContentLength = int64(v.Len())
		buf := v.Bytes()
		r.GetBody = func() (io.ReadCloser, error) {
			r := bytes.NewReader(buf)
			return ioutil.NopCloser(r), nil
		}
	case *bytes.Reader:
		r.ContentLength = int64(v.Len())
		snapshot := *v
		r.GetBody = func() (io.ReadCloser, error) {
			r := snapshot
			return ioutil.NopCloser(&r), nil
		}
	case *strings.Reader:
		r.ContentLength = int64(v.Len())
		snapshot := *v
		r.GetBody = func() (io.ReadCloser, error) {
			r := snapshot
			return ioutil.NopCloser(&r), nil
		}
	default:
	}
	if r.GetBody != nil && r.ContentLength == 0 {
		r.Body = http.NoBody
		r.GetBody = func() (io.ReadCloser, error) { return http.NoBody, nil }
	}
}

// SetHeaders set headers from a map for the request
func (r *Request) SetHeaders(m map[string]string) *Request {
	for k, v := range m {
		r.SetHeader(k, v)
	}
	return r
}

// AddHeaders add headers from a map for the request
func (r *Request) AddHeaders(m map[string]string) *Request {
	for k, v := range m {
		r.AddHeader(k, v)
	}
	return r
}

// SetHeader set a header for the request
func (r *Request) SetHeader(key, value string) *Request {
	r.Header.Set(key, value)
	return r
}

// AddHeader add a header for the request
func (r *Request) AddHeader(key, value string) *Request {
	r.Header.Add(key, value)
	return r
}

// parseScheme parse scheme
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
