package spiders

import (
	"net/http"
	"testing"
)

func TestYUESPRequest_Request(t *testing.T) {
	test := YUESPRequest{
		SpiderRequest{
			Header: http.Header{},
			URL:    "http://www.yuesp.com/play/index6652-0-14.html",
		},
	}
	err := test.Request()
	if err != nil {
		t.Error(err)
	}
}
