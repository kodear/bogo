package spiders

import (
	"encoding/json"
	"errors"
	"net/url"
	"regexp"
	"strconv"
)

const BiliBiliBangumiApi = "https://api.bilibili.com/pgc/player/web/playurl"

type BiliBiliBangumiBody struct {
	biliBody
	Data biliData `json:"result"`
}

type bilibangumiHTML struct {
	Loaded    bool `json:"loaded"`
	MediaInfo struct {
		Title string `json:"title"`
	} `json:"mediaInfo"`
	EPInfo info   `json:"EpInfo"`
	EPList []info `json:"EpList"`
}

type info struct {
	Aid   int    `json:"aid"`
	Cid   int    `json:"cid"`
	Title string `json:"title"`
}

type BiliBiliBangumi struct {
	SpiderObject
}

func (b *BiliBiliBangumi) Parse(r string) (body Body, err error) {
	bytes, err := b.DownloadWeb(r, url.Values{}, map[string]string{})
	if err != nil {
		return
	}

	re := regexp.MustCompile(`__INITIAL_STATE__=(.*);\(function\(\)`)
	match := re.FindAllSubmatch(bytes, -1)

	if len(match) == 0 || len(match[0]) < 2 {
		err = errors.New("parse web page error")
		return
	}

	var htmlBody bilibangumiHTML
	err = json.Unmarshal(match[0][1], &htmlBody)
	if err != nil {
		return
	}

	var videoMeta info
	if htmlBody.EPInfo.Aid == -1 || htmlBody.EPInfo.Cid == -1 {
		videoMeta = htmlBody.EPList[0]
	} else {
		videoMeta = htmlBody.EPInfo
	}

	for index, v := range biliQuality {
		var data BiliBiliBangumiBody
		err = b.DownloadJson("GET", BiliBiliBangumiApi, url.Values{
			"qn":   []string{strconv.Itoa(v)},
			"avid": []string{strconv.Itoa(videoMeta.Aid)},
			"cid":  []string{strconv.Itoa(videoMeta.Cid)},
		}, nil, map[string]string{"Referer": r}, &data)

		if err != nil {
			return
		} else if data.Code != 0 {
			err = errors.New(data.Message)
			return
		}

		var repeat bool
		for _, x := range body.VideoList {
			if x.Quality == quality[data.Data.Quality] {
				repeat = true
				break
			}
		}

		if repeat {
			continue
		}

		var size int
		for _, t := range data.Data.Durl {
			size += t.Size
		}

		var links []VideoAttr
		for _, k := range data.Data.Durl {
			links = append(links, VideoAttr{
				URL:   k.Url,
				Order: k.Order,
				Size:  k.Size,
			})
		}

		var downloadProtocol string
		var format string
		if data.Data.Quality == 15 || data.Data.Quality == 16 {
			format = "mp4"
		} else {
			format = "flv"
		}

		if len(links) == 0 {
			err = errors.New("the download address is not matched")
			return
		} else if len(links) == 1 {
			downloadProtocol = "http"
		} else if format == "flv" {
			downloadProtocol = "httpSegFlv"
		} else if format == "mp4" {
			downloadProtocol = "httpSegMp4"
		}

		body.VideoList = append(body.VideoList, &VideoBody{
			ID:               videoMeta.Cid + index,
			Title:            htmlBody.MediaInfo.Title,
			Part:             videoMeta.Title,
			Format:           format,
			Size:             size,
			Duration:         data.Data.Timelength / 1000,
			StreamType:       data.Data.Format,
			Quality:          quality[data.Data.Quality],
			Links:            links,
			DownloadProtocol: downloadProtocol,
			DownloadHeaders:  map[string]string{"Referer": r},
		})
	}

	return
}

func (b *BiliBiliBangumi) Name() string {
	return "bilibili"
}

func (b *BiliBiliBangumi) WebName() string {
	return "哔哩哔哩视频网"
}

func (b *BiliBiliBangumi) Pattern() string {
	return `https?://(?:www\.)?bilibili\.com/bangumi/play/(?:ss|ep)\d+`
}
