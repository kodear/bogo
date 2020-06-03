package spider

import "testing"

func TestNewSpider(t *testing.T) {
	cls, err := NewSpider("https://www.acfun.cn/v/ac15918856")
	if err != nil {
		t.Error(err)
	}

	cls.Initialization(nil, nil)
	err = cls.Request()

	if err != nil {
		t.Error(err)
	}
}
