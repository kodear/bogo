package spider

import (
	"testing"
)

func TestBILIBILIBangUmiClient_Request(t *testing.T) {
	test := BILIBILIBangUmiClient{}
	test.Initialization("https://www.bilibili.com/bangumi/play/ep15014", nil, nil)
	err := test.Request()
	if err != nil {
		t.Error(err)
	}
}
