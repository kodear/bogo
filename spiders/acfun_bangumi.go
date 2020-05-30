package spiders

import (
	"encoding/json"
	"github.com/zhxingy/bogo/exception"
	"strings"
)

type ACFUNBangUmiRequest struct {
	SpiderRequest
}

func (cls *ACFUNBangUmiRequest) Expression() string {
	return `https?://(?:www\.)?acfun\.cn/bangumi/`
}

func (cls *ACFUNBangUmiRequest) Args() *SpiderArgs {
	return &SpiderArgs{
		"www.acfun.com",
		"acfun番剧",
		Cookie{},
	}
}

func (cls *ACFUNBangUmiRequest) Request() (err error) {
	selector, err := cls.request(cls.URL, nil)
	if err != nil {
		return exception.HTTPHtmlException(err)
	}

	var video struct {
		Title            string `json:"bangumiTitle"`
		Part             string `json:"episodeName"`
		CurrentVideoInfo struct {
			DurationMillis float32 `json:"durationMillis"`
			KsPlayJson     string  `json:"ksPlayJson"`
		} `json:"currentVideoInfo"`
	}
	err = selector.ReByJson(`window.bangumiData = (\{.*\});`, &video)
	if err != nil {
		return err
	}

	var currentVideo struct {
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
	err = json.Unmarshal([]byte(video.CurrentVideoInfo.KsPlayJson), &currentVideo)
	if err != nil {
		return exception.JSONParseException(err)
	}

	for index, v := range currentVideo.AdaptationSet.Representation {
		cls.Response = append(cls.Response, &SpiderResponse{
			ID:     index + 1,
			Title:  strings.TrimSpace(video.Title),
			Part:   video.Part,
			Format: "ts",
			Width:  v.Width,
			Height: v.Height,
			Links: []URLAttr{
				URLAttr{
					URL: v.Url,
				},
			},
			Quality:          v.Quality,
			Duration:         int(v.Duration),
			DownloadProtocol: "hls",
		})
	}

	return
}
