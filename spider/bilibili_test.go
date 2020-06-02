package spider

import (
	"testing"
)

func TestBILIBILIClient_Request(t *testing.T) {
	test := BILIBILIClient{}
	test.Initialization("https://www.bilibili.com/video/BV1mx411E7gn", nil)
	err := test.Request()
	if err != nil {
		t.Error(err)
	}
}
