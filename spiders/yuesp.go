package spiders

import (
	"errors"
	url2 "net/url"
	"path/filepath"
)

type YuespIE struct {
	Spider
}

func (tv *YuespIE) CookieName() string {
	return "yuesp"
}

func (tv *YuespIE) Name() string {
	return "yuesp"
}

func (tv *YuespIE) Domain() string {
	return "http://www.yuesp.com/"
}

func (tv *YuespIE) Pattern() string {
	// http://www.yuesp.com/play/index6424-1-0.html
	// http://www.yuesp.com/play/index6424-0-20.html
	// http://www.yuesp.com/play/index6783-0-0.html
	return `https?://(?:www\.)?yuesp\.com/play/index(\d+)`
}

func (tv *YuespIE) Parse(url string) (body []Response, err error) {
	response, err := tv.DownloadWebPage(url, url2.Values{}, map[string]string{})
	if err != nil {
		return
	}

	//<script>var now="http://cx.anna.run/yunmp4/8150324273603888.mp4";var pn="tybf";
	matchResult, err := response.Search(`<script>var now="(.*)";var pn="(.*)"; var next`)
	if len(matchResult) == 0 || len(matchResult[0]) < 3 {
		err = errors.New("parse web page error")
		return
	}

	newURL := matchResult[0][1]
	mark := matchResult[0][2]

	// 获取Title, Part
	// var playn='法证先锋4粤语版', playp='第1集', playerh
	title := filepath.Base(url)
	part := ""
	matchResult, err = response.Search(`var playn='(.*)', playp='(.*)', playerh`)
	if len(matchResult) > 0 && len(matchResult[0]) >= 3 {
		title = matchResult[0][1]
		part = matchResult[0][2]
	}

	var downloadURL string
	var downloadProtocol string
	if mark == "tybf" {
		response, err = tv.DownloadWebPage(newURL, url2.Values{}, map[string]string{"User-Agent": UserAgent, "Referer": url})
		if err != nil {
			return
		}

		// 搜索播放地址
		// var url = 'https://cloud189-shzh-person.oos-gdsz.ctyunapi.cn/456b0df9-c794-4b9b-8d7b-f4e44651bb7d?x-amz-UFID=8150324273603888&x-amz-CLIENTNETWORK=UNKNOWN&x-amz-FSIZE=390231279&response-content-type=video/mp4&Expires=1584455504&x-amz-UID=330954580&response-content-disposition=attachment%3Bfilename%3D%22%255BYueSP.COM%255D%25E6%2596%25B0%25E5%25B0%2581%25E7%25A5%259E%25E6%25BC%2594%25E4%25B9%258902.mp4%22&AWSAccessKeyId=18bd696e8df5d7a48893&x-amz-limitrate=5120&x-amz-CLOUDTYPEIN=PERSON&x-amz-CLIENTTYPEIN=WEB&Signature=Ug3X6SOevZtLZkSB2HK55ecYXQs%3D'；
		matchResult, err = response.Search(`var url = '(.*)';`)
		if len(matchResult) == 0 || len(matchResult[0]) < 2 {
			err = errors.New("parse web page error")
			return
		}

		downloadURL = matchResult[0][1]
		downloadProtocol = "http"
	} else if mark == "m3u8" {
		downloadURL = newURL
		downloadProtocol = "hls"
	} else {
		err = errors.New("this tag is not supported: " + mark)
	}

	body = append(body, Response{
		ID:      1,
		Title:   title,
		Part:    part,
		Format:  "mp4",
		Quality: "720P",
		Links: []URLAttr{
			{
				URL: downloadURL,
			},
		},
		DownloadProtocol: downloadProtocol,
		DownloadHeaders:  map[string]string{"User-Agent": UserAgent, "Referer": newURL},
	})

	return
}
