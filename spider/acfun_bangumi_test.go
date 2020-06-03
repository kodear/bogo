package spider

import (
	"testing"
)

func TestACFUNBangUmiClient_Request(t *testing.T) {
	test := ACFUNBangUmiClient{}
	test.SetURL("https://www.acfun.cn/bangumi/aa6002267")
	test.Initialization(nil, nil)
	err := test.Request()
	if err != nil {
		t.Error(err)
	}
}
