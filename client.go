package httpgo

import (
	"io"
	"net"
	"net/http"
	"net/http/cookiejar"
	"time"
)

var DefaultClient = NewClient()

type Handle func(*Request) *Response

type Client struct {
	*http.Client
	handle Handle
}

// NewClient 构建client结构
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
				MaxIdleConns:           50,
				IdleConnTimeout:        time.Second * 60,
				TLSHandshakeTimeout:    time.Second * 5,
				ExpectContinueTimeout:  time.Second * 1,
				MaxResponseHeaderBytes: 1024 * 5, // 限制响应头的大小，避免依赖的服务过多使用响应头
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

func (c *Client) Do(req *Request) *Response {
	res := c.handle(req)
	if res == nil {
		return &Response{
			Req: req,
		}
	}
	return res
}

func basicDo(c *Client) Handle {
	return func(req *Request) (resp *Response) {
		resp = &Response{
			Req: req,
		}
		resp.Response, resp.Err = c.Client.Do(req.Request)

		if resp.Err != nil {
			return
		}
		defer resp.Response.Body.Close()
		resp.Body, resp.Err = io.ReadAll(resp.Response.Body)
		return
	}
}
