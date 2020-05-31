package spider

import (
	"encoding/json"
	"errors"
	"github.com/zhxingy/bogo/exception"
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
		return exception.HTTPHtmlException(err)
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
		return exception.ServerAuthException(errors.New(""))
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
		cls.response = append(cls.response, &Response{
			ID:     index + 1,
			Title:  strings.TrimSpace(video.Title),
			Part:   video.CurrentVideoInfo.Part,
			Format: "ts",
			Width:  v.Width,
			Height: v.Height,
			Links: []URLAttr{
				{
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
