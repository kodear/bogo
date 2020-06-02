package spider

import (
	"errors"
	"net/http"
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
		return DownloadHtmlErr(err)
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

	cls.response = &Response{
		Title:  data.MediaInfo.Title,
		Part:   part,
		Stream: []Stream{},
	}
	for _, qualityID := range qualityIds {
		cls.Header.Add("Referer", cls.URL)
		response, err = cls.request("https://api.bilibili.com/pgc/player/web/playurl?", url.Values{
			"qn":   []string{strconv.Itoa(qualityID)},
			"avid": []string{strconv.Itoa(aid)},
			"cid":  []string{strconv.Itoa(cid)},
		})
		if err != nil {
			return ParseJsonErr(err)
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
			Duration:         json.Data.Timelength / 1000,
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
