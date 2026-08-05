package request

import "net/http"

type Request struct {
	raw *http.Request
}

func New(r *http.Request) *Request {
	return &Request{
		raw: r,
	}
}

func (req *Request) Method() string {
	return req.raw.Method
}

func (req *Request) Path() string {
	return req.raw.URL.Path
}

func (req *Request) Query(name string) string {
	return req.raw.URL.Query().Get(name)
}

func (req *Request) Header(name string) string {
	return req.raw.Header.Get(name)
}

func (r *Request) Raw() *http.Request {
	return r.raw
}