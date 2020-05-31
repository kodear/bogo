package spider

import (
	"net/http"
	"testing"
)

func TestACFUNClient_Request(t *testing.T) {
	test := ACFUNClient{
		Client{
			Header: http.Header{},
			URL:    "https://www.acfun.cn/v/ac15533372",
		},
	}
	err := test.Request()
	if err != nil {
		t.Error(err)
	}
}
