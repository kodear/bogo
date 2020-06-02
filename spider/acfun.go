package spider

import (
	"encoding/json"
	"errors"
	"strings"
)

type ACFUNClient struct {
	Client
}

func (cls *ACFUNClient) Meta() *Meta {
	return &Meta{
		Domain:     "https://www.acfun.cn/",
		Name:       "acfun",
		Expression: `https?://(?:www\.)?acfun\.cn/v/ac\d+`,
		Cookie:     Cookie{},
	}
}

func (cls *ACFUNClient) Request() (err error) {
	x, err := cls.request(cls.URL, nil)
	if err != nil {
		return DownloadHtmlErr(err)
	}

	var video struct {
		Result           int    `json:"result"`
		Title            string `json:"title"`
		CurrentVideoInfo struct {
			Part           string  `json:"title"`
			DurationMillis float32 `json:"durationMillis"`
			KsPlayJson     string  `json:"ksPlayJson"`
		} `json:"currentVideoInfo"`
	}
	err = x.ReByJson(`window.videoInfo = (\{.*\});`, &video)
	if err != nil {
		return err
	}
	if video.Result != 0 {
		return ServerAuthErr(errors.New(""))
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
		Part:   video.CurrentVideoInfo.Part,
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
