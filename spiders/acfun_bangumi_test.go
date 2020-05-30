package spiders

import (
	"net/http"
	"testing"
)

func TestACFUNBangUmiClient_Request(t *testing.T) {
	test := ACFUNBangUmiClient{
		Client{
			Header: http.Header{},
			URL:    "https://www.acfun.cn/bangumi/aa6002267",
		},
	}
	err := test.Request()
	if err != nil {
		t.Error(err)
	}
}
