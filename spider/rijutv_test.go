package spider

import (
	"testing"
)

func TestRIJUTVClient_Request(t *testing.T) {
	test := RIJUTVClient{}
	test.SetURL("https://www.rijutv.com/player/99347.html")
	test.Initialization(nil, nil)
	err := test.Request()
	if err != nil {
		t.Error(err)
	}
}
