package spiders

import (
	"encoding/json"
	"errors"
	url2 "net/url"
	"strings"
)

type videoInfo struct {
	Result           int    `json:"result"`
	Title            string `json:"title"`
	CurrentVideoInfo struct {
		Part           string  `json:"title"`
		DurationMillis float32 `json:"durationMillis"`
		KsPlayJson     string  `json:"ksPlayJson"`
	} `json:"currentVideoInfo"`
}

type bangumiData struct {
	Title            string `json:"bangumiTitle"`
	Part             string `json:"episodeName"`
	CurrentVideoInfo struct {
		DurationMillis float32 `json:"durationMillis"`
		KsPlayJson     string  `json:"ksPlayJson"`
	} `json:"currentVideoInfo"`
}

type currentVideoInfo struct {
	AdaptationSet struct {
		Representation []struct {
			Url      string  `json:"url"`
			Width    int     `json:"width"`
			Height   int     `json:"height"`
			Quality  string  `json:"qualityType"`
			Duration float32 `json:"duration"`
		} `json:"representation"`
	} `json:"adaptationSet"`
}

// Acfun 弹幕视频网

type AcfunIE struct {
	Spider
}

func (tv *AcfunIE) CookieName() string {
	return "acfun"
}

func (tv *AcfunIE) Name() string {
	return "Acfun"
}

func (tv *AcfunIE) Domain() string {
	return "https://www.acfun.cn/"
}

func (tv *AcfunIE) Pattern() string {
	return `https?://(?:www\.)?acfun\.cn/v/ac\d+`
}

func (tv *AcfunIE) Parse(url string) (body []Response, err error) {
	response, err := tv.DownloadWebPage(url, url2.Values{}, map[string]string{})
	if err != nil {
		return
	}

	matchResult, err := response.Search(`window.videoInfo = (\{.*\});`)
	if len(matchResult) == 0 || len(matchResult[0]) < 2 {
		err = errors.New("parse web page error")
		return
	}

	var data videoInfo
	err = json.Unmarshal([]byte(matchResult[0][1]), &data)
	if err != nil {
		return
	}

	if data.Result != 0 {
		err = errors.New("access denied")
	}

	var videoData currentVideoInfo
	err = json.Unmarshal([]byte(data.CurrentVideoInfo.KsPlayJson), &videoData)
	if err != nil {
		return
	}

	for index, video := range videoData.AdaptationSet.Representation {
		body = append(body, Response{
			ID:     index + 1,
			Title:  strings.TrimSpace(data.Title),
			Part:   data.CurrentVideoInfo.Part,
			Format: "mp4",
			Width:  video.Width,
			Height: video.Height,
			Links: []URLAttr{
				URLAttr{
					URL: video.Url,
				},
			},
			Quality:          video.Quality,
			Duration:         int(video.Duration),
			DownloadProtocol: "hls",
		})
	}

	return
}

// Acfun 番剧

type AcfunBangumiIE struct {
	Spider
}

func (tv *AcfunBangumiIE) CookieName() string {
	return "acfun-bangumi"
}

func (tv *AcfunBangumiIE) Name() string {
	return "Acfun番剧"
}

func (tv *AcfunBangumiIE) Domain() string {
	return "https://www.acfun.cn/v/list155/index.htm"
}

func (tv *AcfunBangumiIE) Pattern() string {
	return `https?://(?:www\.)?acfun\.cn/bangumi/`
}

func (tv *AcfunBangumiIE) Parse(url string) (body []Response, err error) {
	response, err := tv.DownloadWebPage(url, url2.Values{}, map[string]string{})
	if err != nil {
		return
	}

	matchResult, err := response.Search(`window.bangumiData = (\{.*\});`)
	if len(matchResult) == 0 || len(matchResult[0]) < 2 {
		err = errors.New("parse web page error")
		return
	}

	var data bangumiData
	err = json.Unmarshal([]byte(matchResult[0][1]), &data)
	if err != nil {
		return
	}

	var videoData currentVideoInfo
	err = json.Unmarshal([]byte(data.CurrentVideoInfo.KsPlayJson), &videoData)
	if err != nil {
		return
	}

	for index, video := range videoData.AdaptationSet.Representation {
		body = append(body, Response{
			ID:      index + 1,
			Title:   strings.TrimSpace(data.Title),
			Part:    data.Part,
			Format:  "mp4",
			Width:   video.Width,
			Height:  video.Height,
			Quality: video.Quality,
			Links: []URLAttr{
				URLAttr{
					URL: video.Url,
				},
			},
			Duration:         int(video.Duration),
			DownloadProtocol: "hls",
		})
	}

	return

}
