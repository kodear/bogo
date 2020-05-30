package spiders

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/zhxingy/bogo/exception"
	"github.com/zhxingy/bogo/selector"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type MGTVRequest struct {
	SpiderRequest
}

func (cls *MGTVRequest) Expression() string {
	// https://www.mgtv.com/b/332228/6589904.html?fpa=55&fpos=2
	// https://www.mgtv.com/b/304167/3971620.html
	// https://www.mgtv.com/l/99999286/6607635.html?fpa=1173&fpos=2
	return `https?://(?:www\.)?mgtv\.com/(?:b|l)/\d+/(?P<id>\d+)`
}

func (cls *MGTVRequest) Args() *SpiderArgs {
	return &SpiderArgs{
		"www.mgtv.com",
		"芒果TV",
		Cookie{
			"mgtv",
			true,
			[]string{".mgtv.com"},
		},
	}
}

func (cls *MGTVRequest) Request() (err error) {
	var vid string
	var x selector.Selector
	x = []byte(cls.URL)
	err = x.Re(cls.Expression(), &vid)
	if err != nil {
		return exception.TextParseException(err)
	}

	cls.Header.Add("Referer", cls.URL)
	response, err := cls.request("https://pcweb.api.mgtv.com/player/video?", url.Values{
		"tk2":      []string{cls.sign()},
		"video_id": []string{vid},
	})
	if err != nil {
		return exception.HTTPHtmlException(err)
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

		return exception.JSONParseException(err)
	}
	if auth.Code != 200 {
		return exception.ServerAuthException(errors.New(auth.Msg))
	}
	if auth.Data.Atc.Pm2 == "" {
		return exception.ServerAuthException(errors.New("pm not obtained"))
	}
	response, err = cls.request("https://pcweb.api.mgtv.com/player/getSource?", url.Values{
		"_support": []string{"10000000"},
		"tk2":      []string{cls.sign()},
		"video_id": []string{vid},
		"type":     []string{"pch5"},
		"pm2":      []string{auth.Data.Atc.Pm2},
	})
	if err != nil {
		return exception.HTTPJsonException(err)
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
		return exception.JSONParseException(err)
	}
	if json.Code != 200 {
		return exception.ServerAuthException(errors.New(json.Msg))
	}

	var QualityIDByString = map[int]string{
		1: "270P",
		2: "540P",
		3: "720P",
		4: "1080P",
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
			return exception.HTTPJsonException(err)
		}

		err = response.Json(&video)
		if err != nil {
			return exception.JSONParseException(err)
		}
		if video.Status != "ok" {
			return exception.ServerAuthException(errors.New(video.Info))
		}

		id, _ := strconv.Atoi(stream.Def)
		duration, _ := strconv.Atoi(auth.Data.Info.Duration)
		cls.Response = append(cls.Response, &SpiderResponse{
			ID:         id,
			Title:      auth.Data.Info.Title,
			Part:       auth.Data.Info.Series,
			Format:     "ts",
			Duration:   duration,
			StreamType: stream.FileFormat,
			Quality:    QualityIDByString[id],
			Links: []URLAttr{
				{
					URL: video.Info,
				},
			},
			DownloadHeaders:  map[string]string{"Referer": cls.URL},
			DownloadProtocol: "hls",
		})
	}

	return
}

func (cls *MGTVRequest) sign() string {
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
