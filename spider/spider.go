package spider

import (
	"errors"
	"net/http"
)

type Spider interface {
	Request() error
	Response() *Response
	Initialization(jar CookiesJar, header http.Header)
	SetURL(string)
	Meta() *Meta
}

var Spiders = []Spider{
	&ACFUNClient{}, &ACFUNBangUmiClient{},
	&BILIBILIClient{}, &BILIBILIBangUmiClient{},
	&IQIYIClient{},
	&MGTVClient{},
	&QQClient{},
	&RIJUTVClient{},
	&XIGUAClient{},
	&YOUKUClient{},
	&YUESPClient{},
	&HLSClient{},
	&HTTPClient{},
}

func NewSpider(url string) (cls Spider, err error) {
	for _, spider := range Spiders {
		if match(url, spider.Meta().Expression) {
			cls = spider
			cls.SetURL(url)
			break
		}
	}

	if cls == nil {
		return nil, errors.New("not matched to extractor")
	}

	return
}
