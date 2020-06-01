package spider

import (
	"net/http"
	"testing"
)

func TestBILIBILIBangUmiClient_Request(t *testing.T) {
	test := BILIBILIBangUmiClient{
		Client{
			Header: http.Header{},
			URL:    "https://www.bilibili.com/bangumi/play/ep15014",
		},
	}
	err := test.Request()
	if err != nil {
		t.Error(err)
	}

}
