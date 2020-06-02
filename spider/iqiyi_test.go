package spider

import (
	"testing"
)

func TestIQIYIClient_Request(t *testing.T) {
	test := IQIYIClient{}
	test.Initialization("https://www.iqiyi.com/v_19ry0kjwis.html", nil)
	err := test.Request()
	if err != nil {
		t.Error(err)
	}
}
