package spider

import (
	"net/http"
	"testing"
)

func TestXIGUAClient_Request(t *testing.T) {
	test := XIGUAClient{
		Client{
			Header: http.Header{},
			URL:    "https://www.ixigua.com/cinema/album/81Lx4cDshC4/?logTag=Gi5jdJoRqWJDECwhxPt__",
		},
	}
	err := test.Request()
	if err != nil {
		t.Error(err)
	}
}
