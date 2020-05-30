package spiders

import (
	"net/http"
	"testing"
)

func TestACFUNBangUmiRequest_Request(t *testing.T) {
	test := ACFUNBangUmiRequest{
		SpiderRequest{
			Header: http.Header{},
			URL:    "https://www.acfun.cn/bangumi/aa6002267",
		},
	}
	err := test.Request()
	if err != nil {
		t.Error(err)
	}
}
