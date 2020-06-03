package spider

import (
	"testing"
)

func TestQQClient_Request(t *testing.T) {
	test := QQClient{}
	test.Initialization("https://v.qq.com/x/cover/mzc00200uig074r/o00348wcv07.html", nil, nil)
	err := test.Request()
	if err != nil {
		t.Error(err)
	}
}
