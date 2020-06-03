package spider

import (
	"testing"
)

func TestXIGUAClient_Request(t *testing.T) {
	test := XIGUAClient{}
	test.SetURL("https://www.ixigua.com/cinema/album/81Lx4cDshC4/?logTag=Gi5jdJoRqWJDECwhxPt__")
	test.Initialization(nil, nil)
	err := test.Request()
	if err != nil {
		t.Error(err)
	}
}
