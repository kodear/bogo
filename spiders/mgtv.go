package spiders

import (
	"encoding/base64"
	"errors"
	"fmt"
	url2 "net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	mgtvAuthApi = "https://pcweb.api.mgtv.com/player/video"
	mgtvApi     = "https://pcweb.api.mgtv.com/player/getSource"
)

var mgtvQuality = map[int]string{
	1: "270P",
	2: "540P",
	3: "720P",
	4: "1080P",
}

type MgtvIE struct {
	SpiderIE
}

type mgtvAuthResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Atc struct {
			Pm2 string `json:"pm2"`
			Tk2 string `json:"tk2"`
		} `json:"atc"`
		Info struct {
			Title    string `json:"title"`
			Series   string `json:"series"`
			Duration string `json:"duration"`
		}
	} `json:"data"`
}

type mgtvResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		StreamDomain []string `json:"stream_domain"`
		Stream       []stream `json:"stream"`
	} `json:"data"`
}

type stream struct {
	Def        string `json:"def"`
	Name       string `json:"name"`
	FileFormat string `json:"fileformat"`
	URL        string `json:"url"`
}

type mgtvURL struct {
	Info   string `json:"info"`
	Status string `json:"status"`
}

func (tv *MgtvIE) Parse(url string) (body []Response, err error) {
	re := regexp.MustCompile(tv.Pattern())
	match := re.FindAllStringSubmatch(url, -1)
	if len(match) < 1 || len(match[0]) < 2 {
		err = errors.New("parse url vid error")
		return
	}

	vid := match[0][1]
	var authBody mgtvAuthResponse
	resposne, err := tv.DownloadWebPage(mgtvAuthApi, url2.Values{
		"tk2":      []string{sign()},
		"video_id": []string{vid},
	}, map[string]string{"Referer": url, "User-agent": UserAgent})

	if err != nil {
		return
	}

	err = resposne.Json(&authBody)
	if err != nil {
		return
	} else if authBody.Code != 200 {
		err = errors.New(authBody.Msg)
		return
	}

	pm := authBody.Data.Atc.Pm2
	if pm == "" {
		err = errors.New("pm not obtained")
		return
	}

	resposne, err = tv.DownloadWebPage(mgtvApi, url2.Values{
		"_support": []string{"10000000"},
		"tk2":      []string{sign()},
		"video_id": []string{vid},
		"type":     []string{"pch5"},
		"pm2":      []string{pm},
	}, map[string]string{"Referer": url, "User-agent": UserAgent})

	if err != nil {
		return
	}

	var mgtvBody mgtvResponse
	err = resposne.Json(&mgtvBody)
	if err != nil {
		return
	} else if mgtvBody.Code != 200 {
		err = errors.New(mgtvBody.Msg)
		return
	}

	for _, x := range mgtvBody.Data.Stream {
		if x.URL == "" {
			continue
		}

		var URL mgtvURL
		resposne, err = tv.DownloadWebPage(mgtvBody.Data.StreamDomain[0]+x.URL, url2.Values{}, map[string]string{})
		if err != nil {
			return
		}

		err = resposne.Json(&URL)
		if err != nil || URL.Status != "ok" {
			return
		}

		id, _ := strconv.Atoi(x.Def)
		duration, _ := strconv.Atoi(authBody.Data.Info.Duration)
		quality := mgtvQuality[id]

		body = append(body, Response{
			ID:         id,
			Title:      authBody.Data.Info.Title,
			Part:       authBody.Data.Info.Series,
			Format:     "mp4",
			Duration:   duration,
			StreamType: x.FileFormat,
			Quality:    quality,
			Links: []URLAttr{
				{
					URL: URL.Info,
				},
			},
			DownloadHeaders:  map[string]string{"Referer": url},
			DownloadProtocol: "hls",
		})
	}
	return
}

func (tv *MgtvIE) CookieName() string {
	return "mgtv"
}

func (tv *MgtvIE) Name() string {
	return "芒果TV"
}

func (tv *MgtvIE) Domain() *Cookie {
	return &Cookie{"mgtv",true, []string{".mgtv.com"}}
}

func (tv *MgtvIE) Pattern() string {
	//    https://www.mgtv.com/b/332228/6589904.html?fpa=55&fpos=2
	//    https://www.mgtv.com/b/304167/3971620.html
	//    https://www.mgtv.com/l/99999286/6607635.html?fpa=1173&fpos=2

	return `https?://(?:www\.)?mgtv\.com/(?:b|l)/\d+/(?P<id>\d+)`
}

func sign() string {
	str := fmt.Sprintf("did=aaaa|pno=1030|ver=0.3.0301|clit=%d", int(time.Now().Unix()))
	b64str := base64.StdEncoding.EncodeToString([]byte(str))
	newB64str := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(b64str, "+", "-"), "/", "~"), "=", "-")

	bytes := []rune(newB64str)
	for from, to := 0, len(bytes)-1; from < to; from, to = from+1, to-1 {
		bytes[from], bytes[to] = bytes[to], bytes[from]
	}
	str = string(bytes)
	return str
}
