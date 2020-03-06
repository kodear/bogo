package spiders

import (
	"errors"
	"net/url"
	"regexp"
	"strconv"
)

const (
	YoukuApi     = "http://ups.youku.com/ups/get.json"
	YoukuReferer = "http://player.youku.com/embed/"
)

var YouKuQuality = map[string][]string{
	"1080P_hdr": []string{"hls4hd3_sdr"},
	"1080P":     []string{"mp4hd3v2", "mp4hd3", "cmaf4hd3"},
	"720P":      []string{"mp4hd2v2", "mp4hd2", "cmaf4hd2"},
	"540P":      []string{"mp4hd", "cmaf4hd"},
	"360P":      []string{"mp4sd", "3gphd", "flvhd", "cmaf4sd", "cmaf4ld"},
	"auto":      []string{"h264"},
}

type YouKuBody struct {
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

type YouKu struct {
	SpiderObject
}

func (y *YouKu) Parse(r string) (body Body, err error) {
	vid, err := vid(r)
	if err != nil {
		return
	}
	params := url.Values{"ccode": []string{"0512"}, "client_ip": []string{"192.168.1.1"}, "client_ts": []string{"1569140694"}, "vid": []string{vid}, "utid": []string{y.Cookies["cna"]}}

	headers := map[string]string{"Referer": YoukuReferer + vid}
	var data YouKuBody
	err = y.DownloadJson("GET", YoukuApi, params, nil, headers, &data)
	if err != nil {
		return
	}

	if data.Data.Eroor.Code != 0 {
		err = errors.New(data.Data.Eroor.Note)
		return
	}

	var videoList []*VideoBody
	for index, stream := range data.Data.Stream {
		var quality string
		for key, value := range YouKuQuality {
			for _, v := range value {
				if v == stream.StreamType {
					quality = key
					break
				}
			}
		}

		if quality == "" {
			quality = "auto"
		}

		videoList = append(videoList, &VideoBody{
			ID:         index + 1,
			Title:      data.Data.Show.Title,
			Part:       strconv.Itoa(data.Data.Show.Stage),
			Format:     "mp4",
			Size:       stream.StreamExt.Size,
			Duration:   stream.StreamExt.Duration / 1000,
			Width:      stream.Width,
			Height:     stream.Height,
			StreamType: stream.StreamType,
			Quality:    quality,
			Links: []VideoAttr{
				VideoAttr{
					URL:   stream.Url,
					Order: 0,
					Size:  stream.StreamExt.Size,
				},
			},
			DownloadHeaders:  nil,
			DownloadProtocol: "hls",
		})

	}

	body.VideoList = videoList

	return
}

func (y *YouKu) Name() string {
	return "youku"
}

func (y *YouKu) WebName() string {
	return "优酷视频"
}

func (y *YouKu) Pattern() string {
	return `https?://(?:v\.|player\.|video\.)?(?:youku|tudou)\.com/(?:v_show|v_nextstage|embed|v)/(?:id_)?(?P<vid>[a-zA-Z\d]+={0,2})`
}

func vid(link string) (vid string, err error) {
	pattern := `https?://(?:v\.|player\.|video\.)?(?:youku|tudou)\.com/(?:v_show|v_nextstage|embed|v)/(?:id_)?(?P<vid>[a-zA-Z\d]+={0,2})`
	re := regexp.MustCompile(pattern)
	match := re.FindStringSubmatch(link)

	if len(match) < 2 {
		err = errors.New("not matched to this youku url")
	} else {
		vid = match[1]
	}
	return
}
