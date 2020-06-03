package spider

import (
	"errors"
	"net/url"
	"strconv"
)

type YOUKUClient struct {
	Client
}

func (cls *YOUKUClient) Meta() *Meta {
	return &Meta{
		Domain:     "https://www.youku.com/",
		Name:       "优酷视频",
		Expression: `https?://(?:v\.|player\.|video\.)?(?:youku|tudou)\.com/(?:v_show|v_nextstage|embed|v)/(?:id_)?(?P<vid>[a-zA-Z\d]+={0,2})`,
		Cookie: Cookie{
			Name:   "youku",
			Enable: true,
			Domain: []string{".youku.com"},
		},
	}
}

func (cls *YOUKUClient) Request() (err error) {
	var vid string
	var selector Selector
	selector = []byte(cls.URL)
	err = selector.Re(cls.Meta().Expression, &vid)
	if err != nil {
		return ParseTextErr(err)
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
		return DownloadJsonErr(err)
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
		return ParseJsonErr(err)
	}
	if json.Data.Eroor.Code != 0 {
		return ServerAuthErr(errors.New(json.Data.Eroor.Note))
	}

	cls.response = &Response{
		Title:  json.Data.Show.Title,
		Part:   strconv.Itoa(json.Data.Show.Stage),
		Site:   cls.Meta().Name,
		Stream: []Stream{},
	}
	for id, stream := range json.Data.Stream {
		var quality string
		for key, value := range map[string][]string{
			"1080P_hdr": {"hls4hd3_sdr"},
			"1080P":     {"mp4hd3v2", "mp4hd3", "cmaf4hd3"},
			"720P":      {"mp4hd2v2", "mp4hd2", "cmaf4hd2"},
			"540P":      {"mp4hd", "cmaf4hd"},
			"360P":      {"mp4sd", "3gphd", "flvhd", "cmaf4sd", "cmaf4ld"},
			"auto":      {"h264"},
		} {
			for _, v := range value {
				if v == stream.StreamType {
					quality = key
					break
				}
			}
		}

		cls.response.Duration = stream.StreamExt.Duration / 1000
		cls.response.Stream = append(cls.response.Stream, Stream{
			ID:               id + 1,
			Format:           "ts",
			Size:             stream.StreamExt.Size,
			Width:            stream.Width,
			Height:           stream.Height,
			StreamType:       stream.StreamType,
			Quality:          quality,
			URLS:             []string{stream.Url},
			DownloadHeaders:  nil,
			DownloadProtocol: "hls",
		})
	}

	return
}
