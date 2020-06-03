package cmd

import (
	"github.com/zhxingy/bogo/spider"
	"net/http"
)

func extract(url string, header http.Header, jar []*http.Cookie) (response *spider.Response, err error) {
	cls, err := spider.NewSpider(url)
	if err != nil {
		return
	}

	cls.Initialization(jar, header)
	err = cls.Request()
	if err != nil{
		return
	}

	return cls.Response(), nil
}



