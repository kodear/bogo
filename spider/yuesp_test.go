package spider

import (
	"testing"
)

func TestYUESPClient_Request(t *testing.T) {
	test := YUESPClient{}
	test.SetURL("http://www.yuesp.com/play/index6652-0-14.html")
	test.Initialization(nil, nil)
	err := test.Request()
	if err != nil {
		t.Error(err)
	}
}
