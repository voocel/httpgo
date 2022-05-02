package httpgo

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
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
	TEXT            = "text/plain;charset=UTF-8"
	FORM_DATA       = "multipart/form-data;charset=UTF-8"
)

type Request struct {
	*http.Request
	Err        error
	client     *Client
	fileWriter *multipart.Writer
	callback   func(*Response, error) (*Response, error)
}

type Response struct {
	*http.Response
	Req  *Request
	Body []byte
}

type AsyncResponse struct {
	Resp *Response
	Err  error
}

func AsyncGet(rawUrl string, ch chan<- *AsyncResponse) {
	go func() {
		resp, err := Get(rawUrl).Do()
		ch <- &AsyncResponse{
			Resp: resp,
			Err: err,
		}
	}()
}

// GetBody get response body
func (r *Response) GetBody() []byte {
	return r.Body
}

// GetStatusCode get HTTP status code
func (r *Response) GetStatusCode() int {
	return r.StatusCode
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
		callback: func(resp *Response, err error) (*Response, error) {
			return resp, err
		},
	}
}

// Do finish do
func (r *Request) Do() (resp *Response, err error) {
	if r.Err != nil {
		err = r.Err
		return
	}
	resp, err = r.client.do(r)
	return r.callback(resp, err)
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
	v := make(url.Values)
	v.Add(key, value)
	if len(r.URL.RawQuery) > 0 {
		r.URL.RawQuery += "&"
	}
	r.URL.RawQuery += v.Encode()
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

// SetText set the text data
func (r *Request) SetText(v string) *Request {
	r.setBody(strings.NewReader(v))
	r.SetHeader("Content-Type", TEXT)
	return r
}

// SetFile set the file
func (r *Request) SetFile(field, filename string) *Request {
	buf := new(bytes.Buffer)
	if r.fileWriter == nil {
		r.fileWriter = multipart.NewWriter(buf)
	}
	f, err := os.Open(filename)
	if err != nil {
		r.Err = err
		return r
	}
	defer f.Close()

	fw, err := r.fileWriter.CreateFormFile(field, filepath.Base(filename))
	if err != nil {
		r.Err = err
		return r
	}

	_, err = io.Copy(fw, f)
	if err != nil {
		r.Err = err
		return r
	}

	// must be closed multipart before setBody, trail boundary end line by closed
	r.fileWriter.Close()
	r.setBody(buf)
	r.SetHeader("Content-Type", r.fileWriter.FormDataContentType())

	return r
}

// SetTimeout set the request timeout
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

// SetUA set user-agent for the request
func (r *Request) SetUA(ua string) *Request {
	r.Header.Set("User-Agent", ua)
	return r
}

// AddCookie add cookie for the request
func (r *Request) AddCookie(c *http.Cookie) *Request {
	r.Request.AddCookie(c)
	return r
}

// BasicAuth make basic authentication
func (r *Request) BasicAuth(username, password string) *Request {
	r.Request.SetBasicAuth(username, password)
	return r
}

// parseScheme parse request URL
func parseScheme(rawUrl string) string {
	u, err := url.Parse(rawUrl)
	if err != nil {
		panic(err)
	}
	if !u.IsAbs() {
		u.Scheme = "https"
	}
	return u.String()
}
