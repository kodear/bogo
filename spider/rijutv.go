package spider

import "net/http"

type RIJUTVClient struct {
	Client
}

func (cls *RIJUTVClient) Meta() *Meta {
	return &Meta{
		Domain:     "https://www.rijutv.com/",
		Name:       "日剧TV",
		Expression: `https://www\.rijutv\.com/player/(?P<id>\d+)\.html`,
		Cookie:     Cookie{},
	}
}

func (cls *RIJUTVClient) Request() (err error) {
	cls.Header.Add("Referer", cls.URL)
	response, err := cls.request(cls.URL, nil)
	if err != nil {
		return DownloadHtmlErr(err)
	}

	var urlPath, title, part, url string
	err = response.Re(`<iframe id="playPath" width="100%" height="100%" src="(.*)" frameborder="0" allowfullscreen="true" scrolling="no"></iframe>`, &urlPath)
	if err != nil {
		return ParseHtmlErr(err)
	}

	err = response.Re(`<span class="drama-name">(.*)</span></a><span>：(.*)</span>(?:<i class=".*"></i>)?</h1>`, &title, &part)
	if err != nil {
		return ParseHtmlErr(err)
	}

	response, err = cls.request("http:"+urlPath, nil)
	if err != nil {
		return DownloadHtmlErr(err)
	}

	err = response.Re(`url:'(.*)',\n`, &url)
	if err != nil {
		return ParseHtmlErr(err)
	}

	cls.response = &Response{
		Title: title,
		Part:  part,
		Stream: []Stream{
			{
				ID:               1,
				Format:           "ts",
				Quality:          "720P",
				URLS:             []string{url},
				DownloadHeaders:  http.Header{"Referer": []string{cls.URL}, "User-Agent": []string{UserAgent}},
				DownloadProtocol: "hls",
			},
		},
	}

	return
}
