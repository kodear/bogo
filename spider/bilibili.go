package spider

import (
	"errors"
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
		return DownloadJsonErr(err)
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

	var selector Selector
	var page, part string
	var cid, duration, width, height int
	selector = []byte(cls.URL)
	err = selector.Re(cls.Meta().Expression+`.*\?p=(?P<page>\d+)`, &page)
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

	cls.response = &Response{
		Title:    data.VideoData.Title,
		Part:     part,
		Site:     cls.Meta().Name,
		Duration: duration,
		Stream:   []Stream{},
	}
	for _, qualityID := range qualityIds {
		cls.Header = http.Header{}
		response, err = cls.request("https://api.bilibili.com/x/player/playurl?", url.Values{
			"qn":   []string{strconv.Itoa(qualityID)},
			"avid": []string{strconv.Itoa(data.Aid)},
			"cid":  []string{strconv.Itoa(cid)},
		})
		if err != nil {
			return DownloadJsonErr(err)
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
			return ParseJsonErr(err)
		}

		if json.Code != 0 {
			return ServerAuthErr(errors.New(json.Message))
		}

		var size int
		var urls []string
		var protocol, format string

		for _, video := range json.Data.Durl {
			size += video.Size
			urls = append(urls, video.Url)
		}

		if json.Data.Quality == 15 {
			format = "mp4"
		} else {
			format = "flv"
		}

		if len(urls) == 1 {
			protocol = "http"
		} else if len(urls) > 1 && format == "flv" {
			protocol = "flv"
		} else if len(urls) > 1 && format == "mp4" {
			protocol = "ism"
		}

		cls.response.Stream = append(cls.response.Stream, Stream{
			ID:               json.Data.Quality,
			Format:           format,
			Size:             size,
			Width:            width,
			Height:           height,
			StreamType:       json.Data.Format,
			Quality:          qualityIDByString[json.Data.Quality],
			URLS:             urls,
			DownloadProtocol: protocol,
			DownloadHeaders:  http.Header{"Referer": []string{cls.URL}, "User-Agent": []string{UserAgent}},
		})
	}

	key := make(map[int]Stream)
	for _, stream := range cls.response.Stream {
		key[stream.ID] = stream
	}

	var stream []Stream
	for _, k := range key {
		stream = append(stream, k)
	}

	cls.response.Stream = stream
	return

}
