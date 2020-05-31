package spider

import (
	"errors"
	"github.com/zhxingy/bogo/exception"
	"github.com/zhxingy/bogo/selector"
	"net/url"
	"strconv"
)

type YOUKUClient struct {
	Client
}

func (cls *YOUKUClient) Meta() *Meta {
	return &Meta{
		"www.youku.com",
		"优酷视频",
		`https?://(?:v\.|player\.|video\.)?(?:youku|tudou)\.com/(?:v_show|v_nextstage|embed|v)/(?:id_)?(?P<vid>[a-zA-Z\d]+={0,2})`,
		Cookie{
			"youku",
			true,
			[]string{".youku.com", "user.youku.com"},
		},
	}
}

func (cls *YOUKUClient) Request() (err error) {
	var vid string
	var x selector.Selector
	x = []byte(cls.URL)
	err = x.Re(cls.Meta().Expression, &vid)
	if err != nil {
		return exception.TextParseException(err)
	}

	cls.Header.Add("Referer", "http://player.youku.com/embed/"+vid)
	response, err := cls.request("http://ups.youku.com/ups/get.json?", url.Values{
		"ccode":     []string{"0512"},
		"client_ip": []string{"192.168.1.1"},
		"client_ts": []string{"1569140694"},
		"vid":       []string{vid},
		"utid":      []string{cls.CookieJar.Name("cna")},
	})
	if err != nil {
		return exception.HTTPJsonException(err)
	}

	var json struct {
		Cost float32 `json:"cost"`
		Data struct {
			Eroor struct {
				Code int    `json:"code"`
				Ok   bool   `json:"ok"`
				Note string `json:"note"`
			} `json:"error"`
			Show struct {
				Title string `json:"title"`
				Stage int    `json:"stage"`
			} `json:"show"`
			Stream []struct {
				Width      int    `json:"width"`
				Height     int    `json:"height"`
				StreamType string `json:"stream_type"`
				Url        string `json:"m3u8_url"`
				StreamExt  struct {
					Size     int `json:"hls_size"`
					Duration int `json:"hls_duration"`
				} `json:"stream_ext"`
			} `json:"stream"`
		} `json:"data"`
	}
	err = response.Json(&json)
	if err != nil {
		return exception.JSONParseException(err)
	}
	if json.Data.Eroor.Code != 0 {
		return exception.ServerAuthException(errors.New(json.Data.Eroor.Note))
	}

	for id, stream := range json.Data.Stream {
		var quality string
		for key, value := range map[string][]string{
			"1080P_hdr": []string{"hls4hd3_sdr"},
			"1080P":     []string{"mp4hd3v2", "mp4hd3", "cmaf4hd3"},
			"720P":      []string{"mp4hd2v2", "mp4hd2", "cmaf4hd2"},
			"540P":      []string{"mp4hd", "cmaf4hd"},
			"360P":      []string{"mp4sd", "3gphd", "flvhd", "cmaf4sd", "cmaf4ld"},
			"auto":      []string{"h264"},
		} {
			for _, v := range value {
				if v == stream.StreamType {
					quality = key
					break
				}
			}
		}

		cls.response = append(cls.response, &Response{
			ID:         id + 1,
			Title:      json.Data.Show.Title,
			Part:       strconv.Itoa(json.Data.Show.Stage),
			Format:     "ts",
			Size:       stream.StreamExt.Size,
			Duration:   stream.StreamExt.Duration / 1000,
			Width:      stream.Width,
			Height:     stream.Height,
			StreamType: stream.StreamType,
			Quality:    quality,
			Links: []URLAttr{
				{
					URL:   stream.Url,
					Order: 0,
					Size:  stream.StreamExt.Size,
				},
			},
			DownloadHeaders:  nil,
			DownloadProtocol: "hls",
		})
	}

	return
}
