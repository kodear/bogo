package spiders

import (
	"errors"
	url2 "net/url"
	"path/filepath"
)

type RiJuTvIE struct {
	Spider
}

func (tv *RiJuTvIE) Parse(url string) (body []Response, err error) {
	response, err := tv.DownloadWebPage(url, url2.Values{}, map[string]string{})
	if err != nil {
		return
	}
	// <iframe id="playPath" width="100%" height="100%" src="//jiexi.rijutv.com/index.php?path=3d450%2BUdx1TKRbzcj8OSSK%2FGCyPnR3mTkwf5Ko3Kdb3dK1fyuAHfKZHmJF75rzUeNvCED%2Fx145DYTvyTuott%2FrowC1HBi50teifqTA4MHZIax1%2Bq0DlxeJ5raaATCrZR7ShLNA" frameborder="0" allowfullscreen="true" scrolling="no"></iframe>

	matchResult, err := response.Search(`<iframe id="playPath" width="100%" height="100%" src="(.*)" frameborder="0" allowfullscreen="true" scrolling="no"></iframe>`)
	if err != nil {
		return
	}

	if len(matchResult) == 0 || len(matchResult[0]) < 2 {
		err = errors.New("parse web page error")
		return
	}

	urlpath := matchResult[0][1]

	var title string
	var part string
	//<h1><a href="/dongman/14599.html" title="魔进战队煌辉者" target="_blank"><span class="drama-name">魔进战队煌辉者</span></a><span>：第01集</span></h1>
	matchResult, err = response.Search(`<span class="drama-name">(.*)</span></a><span>：(.*)</span>(?:<i class=".*"></i>)?</h1>`)
	if len(matchResult) == 0 || len(matchResult[0]) < 3 {
		title = filepath.Base(url)
	} else {
		title = matchResult[0][1]
		part = matchResult[0][2]
	}

	response, err = tv.DownloadWebPage("http:"+urlpath, url2.Values{}, map[string]string{})
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

	matchResult, err = response.Search(`url:'(.*)',\n`)
	if err != nil {
		return
	}

	if len(matchResult) == 0 || len(matchResult[0]) < 2 {
		err = errors.New("parse web page error")
		return
	}

	if matchResult[0][1] == "" {
		err = errors.New("the playback address was not found")
		return
	}

	body = append(body, Response{
		ID:      1,
		Title:   title,
		Part:    part,
		Format:  "mp4",
		Quality: "720P",
		Links: []URLAttr{
			URLAttr{
				URL: matchResult[0][1],
			},
		},
		DownloadHeaders:  map[string]string{"User-Agent": UserAgent, "Referer": url},
		DownloadProtocol: "hls",
	})

	return
}

func (tv *RiJuTvIE) CookieName() string {
	return "rijutv"
}

func (tv *RiJuTvIE) Name() string {
	return "日剧TV"
}

func (tv *RiJuTvIE) Domain() string {
	return "https://www.rijutv.com/"
}

func (tv *RiJuTvIE) Pattern() string {
	// https://www.rijutv.com/player/90485.html
	// https://www.rijutv.com/player/77181.html
	// https://www.rijutv.com/player/90494.html
	return `https://www\.rijutv\.com/player/(?P<id>\d+)\.html`
}
