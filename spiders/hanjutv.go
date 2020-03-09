package spiders

import (
	"errors"
	"net/url"
	"path/filepath"
	"regexp"
)

type HanJuTv struct {
	SpiderObject
}

func (tv *HanJuTv) Parse(r string) (body Body, err error) {
	bytes, err := tv.DownloadWeb(r, url.Values{}, map[string]string{})
	if err != nil {
		return
	}
	// <iframe id="playPath" width="100%" height="100%" src="//jiexi.rijutv.com/index.php?path=3d450%2BUdx1TKRbzcj8OSSK%2FGCyPnR3mTkwf5Ko3Kdb3dK1fyuAHfKZHmJF75rzUeNvCED%2Fx145DYTvyTuott%2FrowC1HBi50teifqTA4MHZIax1%2Bq0DlxeJ5raaATCrZR7ShLNA" frameborder="0" allowfullscreen="true" scrolling="no"></iframe>
	re := regexp.MustCompile(`<iframe id="playPath" width="100%" height="100%" src="(.*)" frameborder="0" allowfullscreen="true" scrolling="no"></iframe>`)
	match := re.FindAllStringSubmatch(string(bytes), -1)
	if len(match) == 0 || len(match[0]) < 2 {
		err = errors.New("parse web page error")
		return
	}

	urlpath := match[0][1]

	var title string
	var part string

	re = regexp.MustCompile(`<span class="drama-name">(.*)</span></a><span>：(.*)</span>(?:<i class=".*"></i>)?</h1>`)
	match = re.FindAllStringSubmatch(string(bytes), -1)
	if len(match) == 0 || len(match[0]) < 3 {
		title = filepath.Base(r)
	} else {
		title = match[0][1]
		part = match[0][2]
	}

	bytes, err = tv.DownloadWeb("http:"+urlpath, url.Values{}, map[string]string{})
	if err != nil {
		return
	}

	re = regexp.MustCompile(`url: '(.*)',\n`)
	match = re.FindAllStringSubmatch(string(bytes), -1)
	if len(match) == 0 || len(match[0]) < 2 {
		err = errors.New("parse web page error")
		return
	}

	if match[0][1] == "" {
		err = errors.New("the playback address was not found")
		return
	}

	body.VideoList = append(body.VideoList, &VideoBody{
		ID:     1,
		Title:  title,
		Part:   part,
		Format: "mp4",
		Links: []VideoAttr{
			VideoAttr{
				URL: match[0][1],
			},
		},
		DownloadHeaders:  map[string]string{"User-Agent": UserAgent, "Referer": r},
		DownloadProtocol: "hls",
	})

	return
}

func (tv *HanJuTv) Name() string {
	return "hanjutv"
}

func (tv *HanJuTv) WebName() string {
	return "韩剧TV【https://www.hanjutv.com/】"
}

func (tv *HanJuTv) Pattern() string {
	return `https://www\.hanjutv\.com/player/(?P<id>\d+)\.html`
}
