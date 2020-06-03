package spider

import (
	"testing"
)

func TestBILIBILIBangUmiClient_Request(t *testing.T) {
	test := BILIBILIBangUmiClient{}
	test.SetURL("https://www.bilibili.com/bangumi/play/ep15014")
	test.Initialization(nil, nil)
	err := test.Request()
	if err != nil {
		t.Error(err)
	}
}
