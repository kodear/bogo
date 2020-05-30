package spiders

import (
	"net/http"
	"testing"
)

func TestBILIBILIBangUmiRequest_Request(t *testing.T) {
	test := BILIBILIBangUmiRequest{
		SpiderRequest{
			Header: http.Header{},
			URL:    "https://www.bilibili.com/bangumi/play/ep267851",
		},
	}
	err := test.Request()
	if err != nil {
		t.Error(err)
	}

}
