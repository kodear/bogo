package spider

import (
	"testing"
)

func TestRIJUTVClient_Request(t *testing.T) {
	test := RIJUTVClient{}
	test.Initialization("https://www.rijutv.com/player/99347.html", nil)
	err := test.Request()
	if err != nil {
		t.Error(err)
	}
}
