package response

import (
	"io/ioutil"
	"net/http"

	"github.com/segmentio/objconv/json"
)

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
