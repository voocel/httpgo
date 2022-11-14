package httpgo

import (
	"net/http"
)

type HttpClient struct {
	client *Client
}

func NewHttpClient(mfs ...MiddlewareFunc) *HttpClient {
	c := NewClient(mfs...)
	return &HttpClient{
		client: c,
	}
}

func (h *HttpClient) Get(rawUrl string) *Request {
	return h.NewRequest(GET, rawUrl)
}

func (h *HttpClient) Post(rawUrl string) *Request {
	return h.NewRequest(POST, rawUrl)
}

func (h *HttpClient) NewRequest(method, rawUrl string) *Request {
	var err error
	rawUrl, _ = parseScheme(rawUrl)
	req, err := http.NewRequest(method, rawUrl, nil)
	return &Request{
		Request: req,
		Err:     err,
		client:  h.client,
		callback: func(resp *Response, err error) (*Response, error) {
			return resp, err
		},
	}
}
