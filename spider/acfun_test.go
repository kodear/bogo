package spider

import (
	"testing"
)

func TestACFUNClient_Request(t *testing.T) {
	test := ACFUNClient{}
	test.Initialization("https://www.acfun.cn/v/ac15533372", nil, nil)
	err := test.Request()
	if err != nil {
		t.Error(err)
	}
}
