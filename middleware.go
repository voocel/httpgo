package httpgo

import (
	"bytes"
	"io"
	"log"
	"net/http"
)

type TransportFunc func(*http.Request) (*http.Response, error)

func (tf TransportFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return tf(r)
}

type MiddlewareFunc func(http.RoundTripper) http.RoundTripper

func Middleware(t http.RoundTripper, mfs ...MiddlewareFunc) http.RoundTripper {
	rt := t
	for _, mf := range mfs {
		rt = mf(rt)
	}
	return rt
}

func WithLogger(l *log.Logger) MiddlewareFunc {
	return func(rt http.RoundTripper) http.RoundTripper {
		return TransportFunc(func(req *http.Request) (*http.Response, error) {
			buf := new(bytes.Buffer)
			io.Copy(buf, req.Body)
			req.Body = io.NopCloser(buf)
			l.Printf("method: %v, requests: %v", req.URL.Path, buf.String())

			resp, err := rt.RoundTrip(req)

			buf.Reset()
			if err == nil {
				io.Copy(buf, resp.Body)
				resp.Body = io.NopCloser(buf)
			}
			l.Printf("method: %v, response: %v", req.URL.Path, buf.String())

			return resp, err
		})
	}
}

func WithBasicAuth(username, password string) MiddlewareFunc {
	return func(rt http.RoundTripper) http.RoundTripper {
		return TransportFunc(func(req *http.Request) (*http.Response, error) {
			req.SetBasicAuth(username, password)
			return rt.RoundTrip(req)
		})
	}
}
