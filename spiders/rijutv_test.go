package spiders

import (
	"net/http"
	"testing"
)

func TestRIJUTVRequest_Request(t *testing.T) {
	test := RIJUTVRequest{
		SpiderRequest{
			Header: http.Header{},
			URL:    "https://www.rijutv.com/player/99347.html",
		},
	}
	err := test.Request()
	if err != nil {
		t.Error(err)
	}
}
