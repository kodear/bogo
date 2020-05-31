package spider

import (
	"net/http"
	"testing"
)

func TestIQIYIClient_Request(t *testing.T) {
	test := IQIYIClient{
		Client{
			Header: http.Header{},
			URL:    "https://www.iqiyi.com/v_19rte5wm40.html#vfrm=19-9-0-1",
		},
	}
	err := test.Request()
	if err != nil {
		t.Error(err)
	}

}
