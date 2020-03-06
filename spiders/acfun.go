package spiders

import (
	"encoding/json"
	"errors"
	"net/url"
	"regexp"
	"strings"
)

type AcfunVideoBody struct {
	Url      string  `json:"url"`
	Width    int     `json:"width"`
	Height   int     `json:"height"`
	Quality  string  `json:"qualityType"`
	Duration float32 `json:"duration"`
}

type AcfunVideoData struct {
	AdaptationSet struct {
		Representation []AcfunVideoBody `json:"representation"`
	} `json:"adaptationSet"`
}

type AcfunBody struct {
	Result           int    `json:"result"`
	Title            string `json:"title"`
	CurrentVideoInfo struct {
		Part           string  `json:"title"`
		DurationMillis float32 `json:"durationMillis"`
		KsPlayJson     string  `json:"ksPlayJson"`
	} `json:"currentVideoInfo"`
}

type Acfun struct {
	SpiderObject
}

func (a *Acfun) Parse(r string) (body Body, err error) {
	bytes, err := a.DownloadWeb(r, url.Values{}, map[string]string{})
	if err != nil {
		return
	}

	regex := regexp.MustCompile(`window.videoInfo = (\{.*\});`)
	complie := regex.FindAllSubmatch(bytes, -1)
	if len(complie) == 0 || len(complie[0]) < 2 {
		err = errors.New("parse web page error")
		return
	}

	var data AcfunBody
	err = json.Unmarshal(complie[0][1], &data)
	if err != nil {
		return
	}

	if data.Result != 0 {
		err = errors.New("access denied")
	}

	var videoData AcfunVideoData
	err = json.Unmarshal([]byte(data.CurrentVideoInfo.KsPlayJson), &videoData)
	if err != nil {
		return
	}

	for index, video := range videoData.AdaptationSet.Representation {
		body.VideoList = append(body.VideoList, &VideoBody{
			ID:     index + 1,
			Title:  strings.TrimSpace(data.Title),
			Part:   data.CurrentVideoInfo.Part,
			Format: "mp4",
			Width:  video.Width,
			Height: video.Height,
			Links: []VideoAttr{
				VideoAttr{
					URL:   video.Url,
					Order: 0,
					Size:  0,
				},
			},
			Quality:          video.Quality,
			Duration:         int(video.Duration),
			DownloadProtocol: "hls",
		})
	}

	return
}

func (a *Acfun) Name() string {
	return "acfun"
}

func (a *Acfun) WebName() string {
	return "Acfun【https://www.acfun.cn/】"
}

func (a *Acfun) Pattern() string {
	return `https?://(?:www\.)?acfun\.cn/v/ac\d+`
}
