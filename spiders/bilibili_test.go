package spiders

import (
	"net/http"
	"testing"
)

func TestBILIBILIRequest_Request(t *testing.T) {
	test := BILIBILIRequest{
		SpiderRequest{
			Header: http.Header{},
			URL:    "https://www.bilibili.com/video/BV1jJ411c7s3?p=20",
		},
	}
	err := test.Request()
	if err != nil {
		t.Error(err)
	}

}
