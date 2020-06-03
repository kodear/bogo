package spider

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type MGTVClient struct {
	Client
}

func (cls *MGTVClient) Meta() *Meta {
	return &Meta{
		Domain:     "https://www.mgtv.com//",
		Name:       "芒果TV",
		Expression: `https?://(?:www\.)?mgtv\.com/(?:b|l)/\d+/(?P<id>\d+)`,
		Cookie: Cookie{
			Name:   "mgtv",
			Enable: true,
			Domain: []string{".mgtv.com"},
		},
	}
}

func (cls *MGTVClient) Request() (err error) {
	var vid string
	var selector Selector
	selector = []byte(cls.URL)
	err = selector.Re(cls.Meta().Expression, &vid)
	if err != nil {
		return ParseTextErr(err)
	}

	cls.Header.Add("Referer", cls.URL)
	response, err := cls.request("https://pcweb.api.mgtv.com/player/video?", url.Values{
		"tk2":      []string{cls.sign()},
		"video_id": []string{vid},
	})
	if err != nil {
		return DownloadHtmlErr(err)
	}

	var auth struct {
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
	err = response.Json(&auth)
	if err != nil {
		return ParseJsonErr(err)
	}
	if auth.Code != 200 {
		return ServerAuthErr(errors.New(auth.Msg))
	}
	if auth.Data.Atc.Pm2 == "" {
		return ServerAuthErr(errors.New("pm not obtained"))
	}
	response, err = cls.request("https://pcweb.api.mgtv.com/player/getSource?", url.Values{
		"_support": []string{"10000000"},
		"tk2":      []string{cls.sign()},
		"video_id": []string{vid},
		"type":     []string{"pch5"},
		"pm2":      []string{auth.Data.Atc.Pm2},
	})
	if err != nil {
		return ParseJsonErr(err)
	}

	var json struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			StreamDomain []string `json:"stream_domain"`
			Stream       []struct {
				Def        string `json:"def"`
				Name       string `json:"name"`
				FileFormat string `json:"fileformat"`
				URL        string `json:"url"`
			} `json:"stream"`
		} `json:"data"`
	}
	err = response.Json(&json)
	if err != nil {
		return ParseJsonErr(err)
	}
	if json.Code != 200 {
		return ServerAuthErr(errors.New(json.Msg))
	}

	var QualityIDByString = map[int]string{
		1: "270P",
		2: "540P",
		3: "720P",
		4: "1080P",
	}

	cls.response = &Response{
		Title:  auth.Data.Info.Title,
		Part:   auth.Data.Info.Series,
		Site:   cls.Meta().Name,
		Stream: []Stream{},
	}
	for _, stream := range json.Data.Stream {
		if stream.URL == "" {
			continue
		}
		var video struct {
			Info   string `json:"info"`
			Status string `json:"status"`
		}
		response, err = cls.request(json.Data.StreamDomain[0]+stream.URL, nil)
		if err != nil {
			return DownloadJsonErr(err)
		}

		err = response.Json(&video)
		if err != nil {
			return ParseJsonErr(err)
		}
		if video.Status != "ok" {
			return ServerAuthErr(errors.New(video.Info))
		}

		id, _ := strconv.Atoi(stream.Def)
		duration, _ := strconv.Atoi(auth.Data.Info.Duration)
		cls.response.Duration = duration
		cls.response.Stream = append(cls.response.Stream, Stream{
			ID:               id,
			Format:           "ts",
			StreamType:       stream.FileFormat,
			Quality:          QualityIDByString[id],
			URLS:             []string{video.Info},
			DownloadHeaders:  http.Header{"Referer": []string{cls.URL}, "User-Agent": []string{UserAgent}},
			DownloadProtocol: "hls",
		})
	}

	return
}

func (cls *MGTVClient) sign() string {
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
