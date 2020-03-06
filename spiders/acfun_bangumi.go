package spiders

import (
	"encoding/json"
	"errors"
	"net/url"
	"regexp"
	"strings"
)

type AcfunBangumiBody struct {
	Title            string `json:"bangumiTitle"`
	Part             string `json:"episodeName"`
	CurrentVideoInfo struct {
		DurationMillis float32 `json:"durationMillis"`
		KsPlayJson     string  `json:"ksPlayJson"`
	} `json:"currentVideoInfo"`
}

type AcfunBangumi struct {
	SpiderObject
}

func (a *AcfunBangumi) Parse(r string) (body Body, err error) {
	bytes, err := a.DownloadWeb(r, url.Values{}, map[string]string{})
	if err != nil {
		return
	}

	regex := regexp.MustCompile(`window.bangumiData = (\{.*\});`)
	complie := regex.FindAllSubmatch(bytes, -1)
	if len(complie) == 0 || len(complie[0]) < 2 {
		err = errors.New("parse web page error")
		return
	}

	var data AcfunBangumiBody
	err = json.Unmarshal(complie[0][1], &data)
	if err != nil {
		return
	}

	var videoData AcfunVideoData
	err = json.Unmarshal([]byte(data.CurrentVideoInfo.KsPlayJson), &videoData)
	if err != nil {
		return
	}

	for index, video := range videoData.AdaptationSet.Representation {
		body.VideoList = append(body.VideoList, &VideoBody{
			ID:      index + 1,
			Title:   strings.TrimSpace(data.Title),
			Part:    data.Part,
			Format:  "mp4",
			Width:   video.Width,
			Height:  video.Height,
			Quality: video.Quality,
			Links: []VideoAttr{
				VideoAttr{
					URL:   video.Url,
					Order: 0,
					Size:  0,
				},
			},
			Duration:         int(video.Duration),
			DownloadProtocol: "hls",
		})
	}

	return

}

func (a *AcfunBangumi) Name() string {
	return "acfun-bangumi"
}

func (a *AcfunBangumi) WebName() string {
	return "Acfun番剧【https://www.acfun.cn/v/list155/index.htm】"
}

func (a *AcfunBangumi) Pattern() string {
	return `https?://(?:www\.)?acfun\.cn/bangumi/`
}
