package h

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testTransport struct{}

func (t testTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: 400,
		Body:       r.Body,
	}, nil
}

func TestNewClient(t *testing.T) {
	t.Run("succeed", func(t *testing.T) {
		// It's very easy to test a h.Client
		// Simply implement your own http.RoundTripper as testTransport does.
		// You don't have to start a http test server or mock anything.
		client := NewClient().SetTransport(testTransport{}).
			SetBaseURL("baseURL").
			SetHeader("Header", "Value").
			Use(func(r *Request, res *http.Response, err error) (*http.Response, error) {
				return res, errors.New("middleware is called")
			}).
			SetTimeout(2000 * time.Millisecond)

		assert.Equal(t, "baseURL", client.BaseURL)                       // make sure base url is set
		assert.Equal(t, 2*time.Second, client.Client.Timeout)            // make sure timeout is set
		assert.Equal(t, http.Header{"Header": {"Value"}}, client.Header) // make sure headers are set

		res, err := client.Request(http.MethodGet, "www.google.com").
			SetHeader("x", "x").
			SetBody(strings.NewReader("123")).
			Run()

		assert.Equal(t, "middleware is called", err.Error())
		assert.Equal(t, 400, res.StatusCode)
		b, err := ioutil.ReadAll(res.Body)
		assert.Equal(t, "123", string(b))
	})
}
