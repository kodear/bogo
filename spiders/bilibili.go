package spiders

import (
	"encoding/json"
	"errors"
	url2 "net/url"
	"strconv"
)

const (
	biliApi        = "https://api.bilibili.com/x/player/playurl"
	biliBangumiApi = "https://api.bilibili.com/pgc/player/web/playurl"
)

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

type avResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    v      `json:"data"`
}

type bangumiResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    v      `json:"result"`
}

type v struct {
	Quality           int      `json:"quality"`
	Format            string   `json:"format"`
	Timelength        int      `json:"timelength"` // ms
	AcceptFormat      string   `json:"accept_format"`
	AcceptDescription []string `json:"accept_description"`
	AcceptQuality     []int    `json:"accept_quality"`
	Durl              []struct {
		Url    string `json:"url"`
		Order  int    `json:"order"`
		Length int    `json:"length"`
		Size   int    `json:"size"`
	} `json:"durl"`
}

type avWebJson struct {
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

type bangumiWebJson struct {
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

type BiliBiliIE struct {
	Spider
}

func (tv *BiliBiliIE) Parse(url string) (body []Response, err error) {
	response, err := tv.DownloadWebPage(url, url2.Values{}, map[string]string{})
	if err != nil {
		return
	}

	matchResult, err := response.Search(`__INITIAL_STATE__=(.*);\(function\(\)`)
	if err != nil {
		return
	}

	if len(matchResult) == 0 || len(matchResult[0]) < 2 {
		err = errors.New("parse web page error")
		return
	}

	var htmlBody avWebJson
	err = json.Unmarshal([]byte(matchResult[0][1]), &htmlBody)
	if err != nil {
		return
	}

	for _, p := range htmlBody.VideoData.Pages {
		for index, v := range biliQuality {
			var data avResponse
			response, err = tv.DownloadWebPage(biliApi, url2.Values{
				"qn":   []string{strconv.Itoa(v)},
				"avid": []string{strconv.Itoa(htmlBody.Aid)},
				"cid":  []string{strconv.Itoa(p.Cid)},
			}, map[string]string{})

			if err != nil {
				return
			}

			err = response.Json(&data)
			if err != nil {
				return
			} else if data.Code != 0 {
				err = errors.New(data.Message)
				return
			}

			var repeat bool
			for _, x := range body {
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

			var links []URLAttr
			for _, k := range data.Data.Durl {
				links = append(links, URLAttr{
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

			body = append(body, Response{
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
				DownloadHeaders:  map[string]string{"Referer": url, "User-agent": UserAgent},
			})
		}

	}

	return
}

func (tv *BiliBiliIE) CookieName() string {
	return "bilibili"
}

func (tv *BiliBiliIE) Name() string {
	return "哔哩哔哩"
}

func (tv *BiliBiliIE) Domain() string {
	return "https://www.bilibili.com/"
}

func (tv *BiliBiliIE) Pattern() string {
	return `https?://(?:www\.)?bilibili\.com/video/([A-Z])+`
}

// 哔哩哔哩番剧

type BiliBiliBangumiIE struct {
	Spider
}

func (tv *BiliBiliBangumiIE) Parse(url string) (body []Response, err error) {
	response, err := tv.DownloadWebPage(url, url2.Values{}, map[string]string{})
	if err != nil {
		return
	}

	matchResult, err := response.Search(`__INITIAL_STATE__=(.*);\(function\(\)`)
	if err != nil {
		return
	}

	if len(matchResult) == 0 || len(matchResult[0]) < 2 {
		err = errors.New("parse web page error")
		return
	}

	var htmlBody bangumiWebJson
	err = json.Unmarshal([]byte(matchResult[0][1]), &htmlBody)
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
		var data bangumiResponse
		response, err = tv.DownloadWebPage(biliBangumiApi, url2.Values{
			"qn":   []string{strconv.Itoa(v)},
			"avid": []string{strconv.Itoa(videoMeta.Aid)},
			"cid":  []string{strconv.Itoa(videoMeta.Cid)},
		}, map[string]string{"Referer": url})

		if err != nil {
			return
		}

		err = response.Json(&data)
		if err != nil {
			return
		} else if data.Code != 0 {
			err = errors.New(data.Message)
			return
		}

		var repeat bool
		for _, x := range body {
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

		var links []URLAttr
		for _, k := range data.Data.Durl {
			links = append(links, URLAttr{
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

		body = append(body, Response{
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
			DownloadHeaders:  map[string]string{"Referer": url},
		})
	}

	return
}

func (tv *BiliBiliBangumiIE) CookieName() string {
	return "bilibili"
}

func (tv *BiliBiliBangumiIE) Name() string {
	return "哔哩哔哩番剧"
}

func (tv *BiliBiliBangumiIE) Domain() string {
	return "https://www.bilibili.com/anime/"
}

func (tv *BiliBiliBangumiIE) Pattern() string {
	return `https?://(?:www\.)?bilibili\.com/bangumi/play/(?:ss|ep)\d+`
}
