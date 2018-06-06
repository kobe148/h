package h

import (
	"io"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

type MiddlewareFunc func(r *Request, res *http.Response, err error) (*http.Response, error)

type Client struct {
	Client      http.Client      // the client interface
	BaseURL     string           // base url
	Header      http.Header      // per Client common headers
	Middlewares []MiddlewareFunc // middleware functions
}

// NewClient properly initialize a Client for the user.
func NewClient() *Client {
	return &Client{
		Header: make(http.Header),
	}
}

func (c *Client) SetBaseURL(url string) *Client {
	c.BaseURL = url
	return c
}

func (c *Client) SetHeader(k, v string) *Client {
	c.Header.Set(k, v)
	return c
}

func (c *Client) SetTimeout(duration time.Duration) *Client {
	c.Client.Timeout = duration
	return c
}

func (c *Client) SetTransport(t http.RoundTripper) *Client {
	c.Client.Transport = t
	return c
}

// Use registers a middleware function for this client.
// A middleware is useful for common logging logic, error handling logic and retry logic.
func (c *Client) Use(middlewareFunc MiddlewareFunc) *Client {
	c.Middlewares = append(c.Middlewares, middlewareFunc)
	return c
}

// Request constructs a Request with given methods and path.
func (c *Client) Request(method string, path string) *Request {
	r, e := http.NewRequest(method, c.BaseURL+path, nil)
	r.Header = c.Header
	return &Request{
		Client:  c,
		Request: r,
		Header:  make(http.Header),
		E:       e,
	}
}

func (c *Client) Run(r *Request) (*http.Response, error) {
	if r.E != nil {
		return nil, errors.WithStack(r.E)
	}
	// Set all client common headers to the underlying http.Request
	for k, v := range r.Client.Header {
		r.Request.Header.Set(k, v[0])
	}
	// Set all request specific headers to the underlying http.Request
	for k, v := range r.Header {
		r.Request.Header.Set(k, v[0])
	}
	res, e := r.Client.Client.Do(r.Request)

	return res, e
}

// Request wraps a http.Request and is more powerful.
type Request struct {
	Client          *Client
	Header          http.Header
	Request         *http.Request
	MiddlewareIndex int
	E               error
}

func (r *Request) SetHeader(k, v string) *Request {
	if r.E != nil {
		return nil
	}
	r.Header.Set(k, v)
	return r
}

func (r *Request) SetBody(body io.Reader) *Request {
	if r.E != nil {
		return nil
	}
	req, e := http.NewRequest(r.Request.Method, r.Request.URL.String(), body)
	r.E = e
	r.Request = req
	return r
}

// Run makes the request over network.
// This is the only public method for a user to do IO.
func (r *Request) Run() (*http.Response, error) {
	if r.E != nil {
		return nil, errors.WithStack(r.E)
	}

	var res *http.Response
	var e error

	// Let the client run the request.
	for r.MiddlewareIndex <= len(r.Client.Middlewares) {
		if r.MiddlewareIndex == 0 {
			res, e = r.Client.Run(r)
			r.MiddlewareIndex++
		} else {
			mid := r.Client.Middlewares[r.MiddlewareIndex-1]
			res, e = mid(r, res, e)
			if r.MiddlewareIndex != 0 {
				r.MiddlewareIndex++
			}
		}
	}

	return res, errors.WithStack(e)
}
