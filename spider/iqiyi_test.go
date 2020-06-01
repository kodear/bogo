package spider

import (
	"net/http"
	"testing"
)

func TestIQIYIClient_Request(t *testing.T) {
	test := IQIYIClient{
		Client{
			Header: http.Header{},
			URL:    "http://www.iqiyi.com/v_19rrjbzvbc.html",
		},
	}
	err := test.Request()
	if err != nil {
		t.Error(err)
	}
}
