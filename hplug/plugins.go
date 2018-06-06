package hplug

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/CreatCodeBuild/h"
	"github.com/segmentio/objconv/json"
)

// Retry x times per y seconds on timeout error.
// Does not change timeout duration.
// Because it is implemented as recursion, there are overhead.
func Retry(times int, second int) h.MiddlewareFunc {
	i := 0
	return func(r *h.Request, res *http.Response, err error) (*http.Response, error) {
		type timeout interface {
			Timeout() bool
		}
		if netErr, ok := err.(timeout); ok && netErr.Timeout() && i < times {
			i++
			time.Sleep(time.Duration(second) * time.Second)
			r.MiddlewareIndex = 0
			return nil, nil
		} else {
			return res, err
		}
	}
}

// JSON dumps response body to data, assuming it's of JSON format.
// data should be a pointer.
// It closes the body on success.
// It does not close the body on error.
func JSON(res http.Response, data interface{}) error {
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, data)
	if err != nil {
		return err
	}
	return res.Body.Close()
}
