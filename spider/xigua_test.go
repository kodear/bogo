package spider

import (
	"testing"
)

func TestXIGUAClient_Request(t *testing.T) {
	test := XIGUAClient{}
	test.Initialization("https://www.ixigua.com/cinema/album/81Lx4cDshC4/?logTag=Gi5jdJoRqWJDECwhxPt__", nil, nil)
	err := test.Request()
	if err != nil {
		t.Error(err)
	}
}
