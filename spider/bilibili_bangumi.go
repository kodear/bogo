package spider

import (
	"errors"
	"github.com/zhxingy/bogo/exception"
	"net/url"
	"strconv"
)

type BILIBILIBangUmiClient struct {
	Client
}

func (cls *BILIBILIBangUmiClient) Meta() *Meta {
	return &Meta{
		Domain:     "https://www.bilibili.com/",
		Name:       "哔哩哔哩番剧",
		Expression: `https?://(?:www\.)?bilibili\.com/bangumi/play/(?:ss|ep)\d+`,
		Cookie: Cookie{
			Name:   "bilibili",
			Enable: true,
			Domain: []string{".bilibili.com"},
		},
	}
}

func (cls *BILIBILIBangUmiClient) Request() (err error) {
	response, err := cls.request(cls.URL, nil)
	if err != nil {
		return exception.HTTPHtmlException(err)
	}

	var data struct {
		Loaded    bool `json:"loaded"`
		MediaInfo struct {
			Title string `json:"title"`
		} `json:"mediaInfo"`
		EPInfo struct {
			Aid   int    `json:"aid"`
			Cid   int    `json:"cid"`
			Title string `json:"title"`
		} `json:"EpInfo"`
		EPList []struct {
			Aid   int    `json:"aid"`
			Cid   int    `json:"cid"`
			Title string `json:"title"`
		} `json:"EpList"`
	}

	err = response.ReByJson(`__INITIAL_STATE__=(.*);\(function\(\)`, &data)
	if err != nil {
		return
	}

	var aid, cid int
	var part string
	if data.EPInfo.Aid == -1 || data.EPInfo.Cid == -1 {
		aid = data.EPList[0].Aid
		cid = data.EPList[0].Cid
		part = data.EPList[0].Title
	} else {
		aid = data.EPInfo.Aid
		cid = data.EPInfo.Cid
		part = data.EPInfo.Title
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
		cls.Header.Add("Referer", cls.URL)
		response, err = cls.request("https://api.bilibili.com/pgc/player/web/playurl?", url.Values{
			"qn":   []string{strconv.Itoa(qualityID)},
			"avid": []string{strconv.Itoa(aid)},
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
			} `json:"result"`
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

		if len(links) > 1 {
			protocol = "http"
		} else if len(links) >= 1 && format == "flv" {
			protocol = "flv"
		} else if len(links) >= 1 && format == "mp4" {
			protocol = "ism"
		}

		cls.response = append(cls.response, &Response{
			ID:               json.Data.Quality,
			Title:            data.MediaInfo.Title,
			Part:             part,
			Format:           format,
			Size:             size,
			Duration:         json.Data.Timelength / 1000,
			StreamType:       json.Data.Format,
			Quality:          qualityIDByString[json.Data.Quality],
			Links:            links,
			DownloadProtocol: protocol,
			DownloadHeaders:  map[string]string{"Referer": cls.URL, "User-Agent": UserAgent},
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