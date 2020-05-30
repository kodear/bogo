package spiders

import (
	"fmt"
	"github.com/zhxingy/bogo/exception"
)

type RIJUTVClient struct {
	Client
}

func (cls *RIJUTVClient) Expression() string {
	// https://www.rijutv.com/player/90485.html
	// https://www.rijutv.com/player/77181.html
	// https://www.rijutv.com/player/90494.html
	return `https://www\.rijutv\.com/player/(?P<id>\d+)\.html`
}

func (cls *RIJUTVClient) Args() *Args {
	return &Args{
		"www.rijutv.com",
		"日剧TV",
		Cookie{},
	}
}

func (cls *RIJUTVClient) Request() (err error) {
	cls.Header.Add("Referer", cls.URL)
	response, err := cls.request(cls.URL, nil)
	if err != nil {
		return exception.HTTPHtmlException(err)
	}

	var urlPath, title, part, url string
	err = response.Re(`<iframe id="playPath" width="100%" height="100%" src="(.*)" frameborder="0" allowfullscreen="true" scrolling="no"></iframe>`, &urlPath)
	if err != nil {
		return exception.HTMLParseException(err)
	}

	err = response.Re(`<span class="drama-name">(.*)</span></a><span>：(.*)</span>(?:<i class=".*"></i>)?</h1>`, &title, &part)
	if err != nil {
		return exception.HTMLParseException(err)
	}

	response, err = cls.request("http:"+urlPath, nil)
	if err != nil {
		return exception.HTTPHtmlException(err)
	}

	err = response.Re(`url:'(.*)',\n`, &url)
	if err != nil {
		fmt.Println("http:" + urlPath)
		return exception.HTMLParseException(err)
	}

	cls.response = append(cls.response, &Response{
		ID:      1,
		Title:   title,
		Part:    part,
		Format:  "ts",
		Quality: "720P",
		Links: []URLAttr{
			{
				URL: url,
			},
		},
		DownloadHeaders:  map[string]string{"User-Agent": UserAgent, "Referer": url},
		DownloadProtocol: "hls",
	})

	return
}
