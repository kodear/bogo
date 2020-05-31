package spider

import (
	"net/http"
	"testing"
)

func TestRIJUTVClient_Request(t *testing.T) {
	test := RIJUTVClient{
		Client{
			Header: http.Header{},
			URL:    "https://www.rijutv.com/player/99347.html",
		},
	}
	err := test.Request()
	if err != nil {
		t.Error(err)
	}
}
