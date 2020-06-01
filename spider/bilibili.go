package spider

import (
	"errors"
	"github.com/zhxingy/bogo/exception"
	"github.com/zhxingy/bogo/selector"
	"net/http"
	"net/url"
	"strconv"
)

type BILIBILIClient struct {
	Client
}

func (cls *BILIBILIClient) Meta() *Meta {
	return &Meta{
		Domain:     "https://www.bilibili.com/",
		Name:       "哔哩哔哩",
		Expression: `https?://(?:www\.)?bilibili\.com/video/[A-Za-z\d]+`,
		Cookie: Cookie{
			Name:   "bilibili",
			Enable: true,
			Domain: []string{".bilibili.com"},
		},
	}
}

func (cls *BILIBILIClient) Request() (err error) {
	response, err := cls.request(cls.URL, nil)
	if err != nil {
		return exception.HTTPHtmlException(err)
	}

	var data struct {
		Aid       int `json:"aid"`
		VideoData struct {
			Title string `json:"title"`
			Pages []struct {
				Cid       int    `json:"cid"`
				Page      int    `json:"page"`
				Part      string `json:"part"`
				Duration  int    `json:"duration"`
				Dimension struct {
					Width  int `json:"width"`
					Height int `json:"height"`
				} `json:"dimension"`
			} `json:"pages"`
		} `json:"videoData"`
	}

	err = response.ReByJson(`__INITIAL_STATE__=(.*);\(function\(\)`, &data)
	if err != nil {
		return
	}

	var x selector.Selector
	var page, part string
	var cid, duration, width, height int
	x = []byte(cls.URL)
	err = x.Re(cls.Meta().Expression+`.*\?p=(?P<page>\d+)`, &page)
	if err != nil {
		page = "1"
	}
	pageInt, _ := strconv.Atoi(page)

	for _, pageInfo := range data.VideoData.Pages {
		if pageInfo.Page == pageInt {
			part = pageInfo.Part
			duration = pageInfo.Duration
			width = pageInfo.Dimension.Width
			height = pageInfo.Dimension.Height
			cid = pageInfo.Cid
			break
		}
	}

	qualityIDByString := map[int]string{
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
	var qualityIds []int
	for qualityId := range qualityIDByString {
		qualityIds = append(qualityIds, qualityId)
	}

	for _, qualityID := range qualityIds {
		response, err = cls.request("https://api.bilibili.com/x/player/playurl?", url.Values{
			"qn":   []string{strconv.Itoa(qualityID)},
			"avid": []string{strconv.Itoa(data.Aid)},
			"cid":  []string{strconv.Itoa(cid)},
		})
		if err != nil {
			return exception.HTTPJsonException(err)
		}

		var json struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
			Data    struct {
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
			} `json:"data"`
		}
		err = response.Json(&json)
		if err != nil {
			return exception.JSONParseException(err)
		}

		if json.Code != 0 {
			return exception.ServerAuthException(errors.New(json.Message))
		}

		var size int
		var links []URLAttr
		var protocol, format string

		for _, video := range json.Data.Durl {
			size += video.Size
			links = append(links, URLAttr{
				URL:   video.Url,
				Order: video.Order,
				Size:  video.Size,
			})
		}

		if json.Data.Quality == 15 || json.Data.Quality == 16 {
			format = "mp4"
		} else {
			format = "flv"
		}

		if len(links) == 1 {
			protocol = "http"
		} else if len(links) > 1 && format == "flv" {
			protocol = "flv"
		} else if len(links) > 1 && format == "mp4" {
			protocol = "ism"
		}

		cls.response = append(cls.response, &Response{
			ID:               json.Data.Quality,
			Title:            data.VideoData.Title,
			Part:             part,
			Format:           format,
			Size:             size,
			Duration:         duration,
			Width:            width,
			Height:           height,
			StreamType:       json.Data.Format,
			Quality:          qualityIDByString[json.Data.Quality],
			Links:            links,
			DownloadProtocol: protocol,
			DownloadHeaders:  http.Header{"Referer": []string{cls.URL}, "User-Agent": []string{UserAgent}},
		})
	}

	key := make(map[int]*Response)
	for _, response := range cls.response {
		key[response.ID] = response
	}

	var Response []*Response
	for _, response := range key {
		Response = append(Response, response)
	}

	cls.response = Response
	return

}
