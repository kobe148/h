package retry

import (
	"net/http"
	"testing"

	"github.com/CreatCodeBuild/h"
	"github.com/stretchr/testify/assert"
)

type testTransport struct{}

type timeout struct {
	t func() bool
}

func (t timeout) Timeout() bool {
	return t.t()
}

func (t timeout) Error() string {
	return ""
}

var i = 0

func (t testTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	if i < 3 {
		i++
		return nil, timeout{t: func() bool { return true }}
	}
	return &http.Response{
		StatusCode: 400,
		Body:       r.Body,
	}, nil
}

func TestRetry(t *testing.T) {
	t.Run("succeed", func(t *testing.T) {
		theOtherMiddleware := 0

		// It's very easy to test a h.Client
		// Simply implement your own http.RoundTripper as testTransport does.
		// You don't have to start a http test server or mock anything.
		client := h.NewClient().
			SetTransport(testTransport{}).
			Use(Retry(3, 1)).
			Use(func(r *h.Request, res *http.Response, err error) (*http.Response, error) {
				theOtherMiddleware++
				return res, nil
			})

		res, err := client.Request(http.MethodGet, "www.google.com").Run()

		assert.Nil(t, err)
		assert.Equal(t, 400, res.StatusCode)
		assert.Nil(t, res.Body)
		assert.Equal(t, 3, i)

		assert.Equalf(t, 1, theOtherMiddleware,
			"the other middleware should only be called once, even though the retry middleware is called 3 times.")
	})
}
