package spider

import (
	"net/http"
	"testing"
)

func TestYUESPClient_Request(t *testing.T) {
	test := YUESPClient{
		Client{
			Header: http.Header{},
			URL:    "http://www.yuesp.com/play/index6652-0-14.html",
		},
	}
	err := test.Request()
	if err != nil {
		t.Error(err)
	}
}
