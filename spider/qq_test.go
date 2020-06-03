package spider

import (
	"testing"
)

func TestQQClient_Request(t *testing.T) {
	test := QQClient{}
	test.SetURL("https://v.qq.com/x/cover/mzc00200uig074r/o00348wcv07.html")
	test.Initialization(nil, nil)
	err := test.Request()
	if err != nil {
		t.Error(err)
	}
}
