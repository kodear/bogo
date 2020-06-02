package spider

import (
	"encoding/json"
	"strings"
)

type ACFUNBangUmiClient struct {
	Client
}

func (cls *ACFUNBangUmiClient) Meta() *Meta {
	return &Meta{
		Domain:     "https://www.acfun.cn/",
		Name:       "acfun番剧",
		Expression: `https?://(?:www\.)?acfun\.cn/bangumi/`,
		Cookie:     Cookie{},
	}
}

func (cls *ACFUNBangUmiClient) Request() (err error) {
	selector, err := cls.request(cls.URL, nil)
	if err != nil {
		return DownloadHtmlErr(err)
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
		return ParseJsonErr(err)
	}

	cls.response = &Response{
		Title:  strings.TrimSpace(video.Title),
		Part:   video.Part,
		Stream: []Stream{},
	}
	for index, v := range currentVideo.AdaptationSet.Representation {
		cls.response.Stream = append(cls.response.Stream, Stream{
			ID:               index + 1,
			Format:           "ts",
			Quality:          v.Quality,
			Duration:         int(v.Duration),
			DownloadProtocol: "hls",
			Width:            v.Width,
			Height:           v.Height,
			URLS:             []string{v.Url},
		})
	}

	return
}
