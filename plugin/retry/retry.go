package retry

import (
	"net/http"
	"time"

	"github.com/CreatCodeBuild/h"
	"github.com/pkg/errors"
)

// Retry x times per y seconds on timeout error.
// Does not change timeout duration.
// Because it is implemented as recursion, there are overhead.
func Retry(times int, second time.Duration) h.MiddlewareFunc {
	return func(r *h.Request, res *http.Response, err error) (*http.Response, error) {
		type timeout interface {
			Timeout() bool
		}

		err = errors.Cause(err)
		netErr, ok := err.(timeout)
		for i := 0; ok && netErr.Timeout() && i < times; i++ {
			time.Sleep(second * time.Second)
			res, err = r.Client.Client.Do(r.Request)
			netErr, ok = err.(timeout)
		}
		return res, err
	}
}
