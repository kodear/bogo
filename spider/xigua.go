package spider

import (
	"encoding/base64"
	"strconv"
)

type XIGUAClient struct {
	Client
}

func (cls *XIGUAClient) Meta() *Meta {
	return &Meta{
		Domain:     "https://www.ixigua.com/",
		Name:       "西瓜视频",
		Expression: `https://www.ixigua.com/cinema/album/([A-Za-z\d]+)`,
		Cookie:     Cookie{},
	}
}

func (cls *XIGUAClient) Request() (err error) {
	response, err := cls.request(cls.URL, nil)
	if err != nil {
		return DownloadHtmlErr(err)
	}

	var json struct {
		AlbumInfo struct {
			Title     string `json:"title"`
			LatestSeq int    `json:"latestSeq"`
			Duration  string `json:"duration"`
		} `json:"albumInfo"`
		VideoResource struct {
			Normal struct {
				VideoList map[string]struct {
					Definition string `json:"definition"`
					Vtype      string `json:"vtype"`
					Vwidth     int    `json:"vwidth"`
					Vheight    int    `json:"vheight"`
					Size       int    `json:"size"`
					URL        string `json:"main_url"`
				} `json:"video_list"`
			} `json:"normal"`
		} `json:"videoResource"`
	}
	err = response.ReByJson(`"Teleplay":(\{"albumInfo".*\}),"Projection"`, &json)
	if err != nil {
		return err
	}

	var index int
	cls.response = &Response{
		Title:  json.AlbumInfo.Title,
		Part:   strconv.Itoa(json.AlbumInfo.LatestSeq),
		Stream: []Stream{},
	}
	for _, video := range json.VideoResource.Normal.VideoList {
		index += 1
		duration, _ := strconv.Atoi(json.AlbumInfo.Duration)
		decodeBytes, _ := base64.StdEncoding.DecodeString(video.URL)
		cls.response.Stream = append(cls.response.Stream, Stream{
			ID:               index,
			Format:           video.Vtype,
			Width:            video.Vwidth,
			Height:           video.Vheight,
			Size:             video.Size,
			StreamType:       video.Vtype,
			URLS:             []string{string(decodeBytes)},
			Quality:          video.Definition,
			Duration:         duration,
			DownloadProtocol: "http",
		})
	}

	return
}
