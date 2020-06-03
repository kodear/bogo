package spider

import (
	"testing"
)

func TestBILIBILIClient_Request(t *testing.T) {
	test := BILIBILIClient{}
	test.SetURL("https://www.bilibili.com/video/BV1mx411E7gn")
	test.Initialization(nil, nil)
	err := test.Request()
	if err != nil {
		t.Error(err)
	}
}
