package spider

import (
	"errors"
	"net/http"
)

type Spider interface {
	Request() error
	Response() *Response
	Initialization(uri string, jar CookiesJar, header http.Header)
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

func Do(uri string, jar []*http.Cookie, header http.Header) (*Response, error) {
	var ie Spider
	for _, spider := range Spiders {
		if match(uri, spider.Meta().Expression) {
			ie = spider
			break
		}
	}

	if ie == nil {
		return nil, errors.New("not matched to extractor")
	}

	ie.Initialization(uri, jar, header)
	err := ie.Request()
	if err != nil {
		return nil, err
	}

	return ie.Response(), nil
}
