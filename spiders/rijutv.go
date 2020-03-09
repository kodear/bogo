package spiders

import (
	"errors"
	"net/url"
	"path/filepath"
	"regexp"
)

type RiJuTv struct {
	SpiderObject
}

func (tv *RiJuTv) Parse(r string) (body Body, err error) {
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
	//<h1><a href="/dongman/14599.html" title="魔进战队煌辉者" target="_blank"><span class="drama-name">魔进战队煌辉者</span></a><span>：第01集</span></h1>
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
	/*
	   window.onload = function(){
	       var dp = new DPlayer({
	           container:document.getElementById('pl'),
	           video:{
	               quality:[{
	                   type:'hls',
	                   name:'高清',
	                   url:'https://iqiyi.cdn9-okzy.com/20200308/7227_d6123f42/index.m3u8',
	               }],
	               defaultQuality: 0,
	           }
	       });
	*/
	re = regexp.MustCompile(`url:'(.*)',\n`)
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
		ID:      1,
		Title:   title,
		Part:    part,
		Format:  "mp4",
		Quality: "720P",
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

func (tv *RiJuTv) Name() string {
	return "rijutv"
}

func (tv *RiJuTv) WebName() string {
	return "日剧TV【https://www.rijutv.com/】"
}

func (tv *RiJuTv) Pattern() string {
	// https://www.rijutv.com/player/90485.html
	// https://www.rijutv.com/player/77181.html
	// https://www.rijutv.com/player/90494.html
	return `https://www\.rijutv\.com/player/(?P<id>\d+)\.html`
}
