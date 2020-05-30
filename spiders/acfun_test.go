package spiders

import (
	"net/http"
	"testing"
)

func TestACFUNRequest_Request(t *testing.T) {
	test := ACFUNRequest{
		SpiderRequest{
			Header: http.Header{},
			URL:    "https://www.acfun.cn/v/ac15533372",
		},
	}
	err := test.Request()
	if err != nil {
		t.Error(err)
	}
}
