package spider

import (
	"testing"
)

func TestACFUNClient_Request(t *testing.T) {
	test := ACFUNClient{}
	test.SetURL("https://www.acfun.cn/v/ac15533372")
	test.Initialization(nil, nil)
	err := test.Request()
	if err != nil {
		t.Error(err)
	}
}
