package spiders

import (
	"encoding/json"
	"errors"
	"net/url"
	"regexp"
	"strconv"
)

const BiliBiliApi = "https://api.bilibili.com/x/player/playurl"

var quality = map[int]string{
	15:  "360P",
	16:  "360P",
	32:  "480P",
	48:  "720P",
	64:  "720P",
	74:  "720P60",
	80:  "1080P",
	112: "1080P+",
	116: "1080P60",
	120: "4K",
}

var biliQuality = []int{16, 32, 48, 64, 74, 80, 112, 116, 120}

type biliBody struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type biliData struct {
	Quality           int        `json:"quality"`
	Format            string     `json:"format"`
	Timelength        int        `json:"timelength"` // ms
	AcceptFormat      string     `json:"accept_format"`
	AcceptDescription []string   `json:"accept_description"`
	AcceptQuality     []int      `json:"accept_quality"`
	Durl              []biliDurl `json:"durl"`
}

type biliDurl struct {
	Url    string `json:"url"`
	Order  int    `json:"order"`
	Length int    `json:"length"`
	Size   int    `json:"size"`
}

type biliAvBody struct {
	biliBody
	Data biliData `json:"data"`
}

type biliHTML struct {
	Aid       int `json:"aid"`
	VideoData struct {
		Title string `json:"title"`
		Pages []struct {
			Cid  int    `json:"cid"`
			Page int    `json:"page"`
			Part string `json:"part"`
		} `json:"pages"`
	} `json:"videoData"`
}

type BiliBili struct {
	SpiderObject
}

func (b *BiliBili) Parse(r string) (body Body, err error) {
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

	var htmlBody biliHTML
	err = json.Unmarshal(match[0][1], &htmlBody)
	if err != nil {
		return
	}

	for _, p := range htmlBody.VideoData.Pages {
		for index, v := range biliQuality {
			var data biliAvBody
			err = b.DownloadJson("GET", BiliBiliApi, url.Values{
				"qn":   []string{strconv.Itoa(v)},
				"avid": []string{strconv.Itoa(htmlBody.Aid)},
				"cid":  []string{strconv.Itoa(p.Cid)},
			}, nil, map[string]string{}, &data)

			if err != nil {
				return
			} else if data.Code != 0 {
				err = errors.New(data.Message)
				return
			}

			var repeat bool
			for _, x := range body.VideoList {
				if x.Quality == quality[data.Data.Quality] && x.Part == p.Part {
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
				ID:               p.Cid + index,
				Title:            htmlBody.VideoData.Title,
				Part:             p.Part,
				Format:           format,
				Size:             size,
				Duration:         data.Data.Timelength / 1000,
				StreamType:       data.Data.Format,
				Quality:          quality[data.Data.Quality],
				Links:            links,
				DownloadProtocol: downloadProtocol,
				DownloadHeaders:  map[string]string{"Referer": r, "User-agent": UserAgent},
			})
		}

	}

	return
}

func (b *BiliBili) Name() string {
	return "bilibili"
}

func (b *BiliBili) WebName() string {
	return "哔哩哔哩【https://www.bilibili.com/】"
}

func (b *BiliBili) Pattern() string {
	return `https?://(?:www\.)?bilibili\.com/video/av\d+`
}
