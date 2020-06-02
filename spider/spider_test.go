package spider

import "testing"

func TestDo(t *testing.T) {
	_, err := Do("https://www.bilibili.com/bangumi/play/ep15014", nil)
	if err != nil {
		t.Error(err)
	}
}
