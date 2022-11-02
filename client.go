package httpgo

import (
	"bytes"
	"io"
	"net"
	"net/http"
	"net/http/cookiejar"
	"time"
)

var DefaultClient = NewClient()

type Handle func(*Request) (*Response, error)

type Client struct {
	*http.Client
	handle Handle
}

// NewClient create a client
func NewClient() *Client {
	jar, _ := cookiejar.New(nil)
	//jar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	c := &Client{
		Client: &http.Client{
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout:   time.Second * 10,
					KeepAlive: time.Second * 30,
				}).DialContext,
				MaxIdleConns:          50,
				IdleConnTimeout:       time.Second * 60,
				TLSHandshakeTimeout:   time.Second * 5,
				ExpectContinueTimeout: time.Second * 1,
				// Limit the size of response headers to avoid excessive use of response headers by dependent services
				MaxResponseHeaderBytes: 1024 * 5,
				DisableCompression:     false,
			},
			CheckRedirect: nil,
			Jar:           jar,
			Timeout:       0,
		},
	}
	c.handle = basicDo(c)
	return c
}

func (c *Client) do(req *Request) (res *Response, err error) {
	res, err = c.handle(req)
	if res == nil {
		return &Response{
			Req: req,
		}, err
	}
	return
}

func basicDo(c *Client) Handle {
	return func(req *Request) (resp *Response, err error) {
		resp = &Response{
			Req: req,
		}
		resp.Response, err = c.Client.Do(req.Request)
		if err != nil {
			return
		}
		defer resp.Response.Body.Close()

		buf := new(bytes.Buffer)
		buf.Grow(1024)
		_, err = io.Copy(buf, resp.Response.Body)
		resp.Body = buf.Bytes()
		buf.Reset()

		return
	}
}
