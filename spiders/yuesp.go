package spiders

import (
	"errors"
	"github.com/zhxingy/bogo/exception"
)

type YUESPClient struct {
	Client
}

func (cls *YUESPClient) Expression() string {
	// http://www.yuesp.com/play/index6424-1-0.html
	// http://www.yuesp.com/play/index6424-0-20.html
	// http://www.yuesp.com/play/index6783-0-0.html
	return `https?://(?:www\.)?yuesp\.com/play/index(\d+)`
}

func (cls *YUESPClient) Args() *Args {
	return &Args{
		"www.yuesp.com",
		"粤视频",
		Cookie{},
	}
}

func (cls *YUESPClient) Request() (err error) {
	response, err := cls.request(cls.URL, nil)
	if err != nil {
		return exception.HTTPHtmlException(err)
	}

	var title, part, mark, newUrl, url, protocol, format string

	err = response.Re(`<script>var now="(.*)";var pn="(.*)"; var next`, &newUrl, &mark)
	if err != nil {
		return exception.HTMLParseException(err)
	}

	err = response.Re(`var playn='(.*)', playp='(.*)', playerh`, &title, &part)
	if err != nil {
		return exception.HTMLParseException(err)
	}

	if mark == "tybf" {
		cls.Header.Add("Referer", cls.URL)
		response, err = cls.request(newUrl, nil)
		if err != nil {
			return exception.HTTPHtmlException(err)
		}

		err = response.Re(`var url = '(.*)';`, &url)
		if err != nil {
			return exception.HTMLParseException(err)
		}

		format = "mp4"
		protocol = "http"
	} else if mark == "m3u8" {
		url = newUrl
		format = "ts"
		protocol = "hls"
	} else {
		return exception.OtherException(errors.New("this tag is not supported " + mark))
	}

	cls.response = append(cls.response, &Response{
		ID:      1,
		Title:   title,
		Part:    part,
		Format:  format,
		Quality: "720P",
		Links: []URLAttr{
			{
				URL: url,
			},
		},
		DownloadProtocol: protocol,
		DownloadHeaders:  map[string]string{"User-Agent": UserAgent, "Referer": newUrl},
	})

	return
}
