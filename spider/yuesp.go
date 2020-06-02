package spider

import (
	"errors"
	"net/http"
)

type YUESPClient struct {
	Client
}

func (cls *YUESPClient) Meta() *Meta {
	return &Meta{
		Domain:     "http://www.yuesp.com/",
		Name:       "粤视频",
		Expression: `https?://(?:www\.)?yuesp\.com/play/index(\d+)`,
		Cookie:     Cookie{},
	}
}

func (cls *YUESPClient) Request() (err error) {
	response, err := cls.request(cls.URL, nil)
	if err != nil {
		return DownloadHtmlErr(err)
	}

	var title, part, mark, newUrl, url, protocol, format string

	err = response.Re(`<script>var now="(.*)";var pn="(.*)"; var next`, &newUrl, &mark)
	if err != nil {
		return ParseHtmlErr(err)
	}

	err = response.Re(`var playn='(.*)', playp='(.*)', playerh`, &title, &part)
	if err != nil {
		return ParseHtmlErr(err)
	}

	if mark == "tybf" {
		format = "mp4"
		protocol = "http"
		cls.Header.Add("Referer", cls.URL)
		response, err = cls.request(newUrl, nil)
		if err != nil {
			return DownloadHtmlErr(err)
		}

		err = response.Re(`var url = '(.*)';`, &url)
		if err != nil {
			return ParseHtmlErr(err)
		}
	} else if mark == "m3u8" {
		url = newUrl
		format = "ts"
		protocol = "hls"
	} else {
		return UnknownErr(errors.New("this tag is not supported " + mark))
	}

	cls.response = &Response{
		Title: title,
		Part:  part,
		Stream: []Stream{
			{
				ID:               1,
				Format:           format,
				Quality:          "720P",
				URLS:             []string{url},
				DownloadProtocol: protocol,
				DownloadHeaders:  http.Header{"Referer": []string{cls.URL}, "User-Agent": []string{UserAgent}},
			},
		},
	}

	return
}
