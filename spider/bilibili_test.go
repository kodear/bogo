package spider

import (
	"net/http"
	"testing"
)

func TestBILIBILIClient_Request(t *testing.T) {
	test := BILIBILIClient{
		Client{
			Header: http.Header{},
			URL:    "https://www.bilibili.com/video/BV1pt41127FZ?p=1",
		},
	}
	err := test.Request()
	if err != nil {
		t.Error(err)
	}

}
