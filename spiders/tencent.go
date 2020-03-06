package spiders

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	tencentCkeyApi = "http://111.59.199.42:9999/ckey?"
	tencentApi     = "https://vd.l.qq.com/proxyhttp"
	tencentAuthApi = "https://access.video.qq.com/user/auth_refresh"
)

type params struct {
	Buid       string `json:"buid"`
	Vinfoparam string `json:"vinfoparam"`
}

type tencentBody struct {
	ErrCode int    `json:"errCode"`
	Vinfo   string `json:"vinfo"`
}

type vinfo struct {
	Msg string `json:"msg"`
	Fl  struct {
		Fi []fi `json:"fi"`
	} `json:"fl"`
	Vl struct {
		Vi []vi `json:"vi"`
	} `json:"vl"`
}

type fi struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Resolution string `json:"resolution"`
	Fs         int    `json:"fs"`
}

type vi struct {
	Ti string `json:"ti"`
	Vw int    `json:"vw"`
	Vt int    `json:"vt"`
	Vh int    `json:"vh"`
	Fs int    `json:"fs"`
	Td string `json:"td"`
	Ul struct {
		Ui []struct {
			URL string `json:"url"`
			Hls struct {
				Ftype string `json:"ftype"`
				Pt    string `json:"pt"`
			} `json:"hls"`
		} `json:"ui"`
	} `json:"ul"`
}

type auth struct {
	Errcode     int    `json:"errcode"`
	AccessToken string `json:"access_token"`
	Vuserid     int    `json:"vuserid"`
	Vusession   string `json:"vusession"`
}

var tencentQuality = []string{"sd", "hd", "shd", "fhd"}

type Tencent struct {
	SpiderObject
}

func (t *Tencent) Parse(r string) (body Body, err error) {
	bytes, err := t.DownloadWeb(r, url.Values{}, map[string]string{})
	if err != nil {
		return
	}
	re := regexp.MustCompile(`vid=(?P<vid>[a-zA-Z\d]+)&`)
	match := re.FindAllStringSubmatch(string(bytes), -1)
	if len(match) == 0 || len(match[0]) < 2 {
		err = errors.New("parse web page error")
		return
	}

	vid := match[0][1]
	tm := int(time.Now().Unix())
	uid := guid()

	bytes, err = t.DownloadWeb(tencentCkeyApi, url.Values{"vid": []string{vid}, "guid": []string{uid}, "tm": []string{strconv.Itoa(tm)}}, map[string]string{})
	if err != nil {
		return
	}
	ckey := string(bytes)

	authParams := url.Values{
		"vappid":  []string{"11059694"},
		"vsecret": []string{"fdf61a6be0aad57132bc5cdf78ac30145b6cd2c1470b0cfe"},
		"type":    []string{"qq"},
		"g_tk":    []string{},
		"g_vstk":  []string{strconv.Itoa(time33(t.Cookies["vqq_vusession"]))},
		"g_actk":  []string{strconv.Itoa(time33(t.Cookies["vqq_access_token"]))},
		"_":       []string{strconv.Itoa(tm * 1000)},
	}

	bytes, err = t.DownloadWeb(tencentAuthApi, authParams, map[string]string{"Referer": r})
	if err == nil {
		re := regexp.MustCompile(`data=(.*)`)
		match := re.FindAllStringSubmatch(string(bytes), -1)
		if len(match) > 0 && len(match[0]) >= 2 {
			var at auth
			err = json.Unmarshal([]byte(match[0][1]), &at)
			if err == nil && at.Errcode == 0 {
				t.SetCookies(fmt.Sprintf("vqq_vusession=%s;vqq_access_token=%s;vqq_vuserid=%d", at.Vusession, at.AccessToken, at.Vuserid))
			}
		}
	}

	for _, q := range tencentQuality {
		p := url.Values{
			"adsid":       []string{""},
			"vid":         []string{vid},
			"dtype":       []string{"3"},
			"show1080p":   []string{"1"},
			"guid":        []string{uid},
			"sdtfrom":     []string{"v1010"},
			"adpinfo":     []string{""},
			"spgzip":      []string{"1"},
			"spau":        []string{"1"},
			"ehost":       []string{r},
			"fp2p":        []string{"1"},
			"spaudio":     []string{"15"},
			"appVer":      []string{"3.5.57"},
			"sphttps":     []string{"1"},
			"platform":    []string{"10201"},
			"charge":      []string{"0"},
			"tm":          []string{strconv.Itoa(tm)},
			"cKey":        []string{ckey},
			"sphls":       []string{"2"},
			"defaultfmt":  []string{"auto"},
			"refer":       []string{"v.qq.com"},
			"spwm":        []string{"4"},
			"flowid":      []string{uid + "_10201"},
			"host":        []string{"v.qq.com"},
			"hdcp":        []string{"hdcp"},
			"defn":        []string{q},
			"defnpayver":  []string{"1"},
			"defsrc":      []string{"2"},
			"isHLS":       []string{"1"},
			"spadseg":     []string{"1"},
			"onlyGetinfo": []string{"true"},
			"encryptVer":  []string{"9.1"},
			"otype":       []string{"ojson"},
			"drm":         []string{"32"},
			"fhdswitch":   []string{"1"},
			"dlver":       []string{"2"},
		}

		bytes, err = json.Marshal(params{Buid: "onlyvinfo", Vinfoparam: p.Encode()})
		if err != nil {
			return
		}

		var data tencentBody

		err = t.DownloadJson("POST", tencentApi, url.Values{}, bytes, map[string]string{"Referer": r}, &data)

		if err != nil {
			return
		}

		if data.ErrCode != 0 {
			err = errors.New("tencent video errcode %d" + strconv.Itoa(data.ErrCode))
			return
		}

		var x vinfo
		err = json.Unmarshal([]byte(data.Vinfo), &x)
		if err != nil {
			return
		}

		if x.Msg != "" {
			continue
		}

		duration, _ := strconv.ParseFloat(x.Vl.Vi[0].Td, 32)

		var (
			id         int
			streamType string
			quality    string
			format     string
			link       string
		)

		for _, p := range x.Fl.Fi {
			if x.Vl.Vi[0].Fs == p.Fs {
				id = p.ID
				streamType = p.Name
				quality = p.Resolution
				break
			}
		}

		if strings.HasSuffix(x.Vl.Vi[0].Ul.Ui[0].URL, "/") {
			if x.Vl.Vi[0].Ul.Ui[0].Hls.Pt == "" {
				continue
			}
			format = x.Vl.Vi[0].Ul.Ui[0].Hls.Ftype
			link = x.Vl.Vi[0].Ul.Ui[0].URL + x.Vl.Vi[0].Ul.Ui[0].Hls.Pt
		} else {
			format = "mp4"
			link = x.Vl.Vi[0].Ul.Ui[0].URL
		}

		body.VideoList = append(body.VideoList, &VideoBody{
			ID:         id,
			Title:      x.Vl.Vi[0].Ti,
			Part:       "",
			Format:     format,
			Size:       x.Vl.Vi[0].Fs,
			Duration:   int(duration),
			Width:      x.Vl.Vi[0].Vw,
			Height:     x.Vl.Vi[0].Vh,
			StreamType: streamType,
			Quality:    quality,
			Links: []VideoAttr{
				VideoAttr{
					URL:   link,
					Order: 0,
					Size:  x.Vl.Vi[0].Fs,
				},
			},
			DownloadProtocol: "hls",
		})
	}

	return
}

func (t *Tencent) Name() string {
	return "tencent"
}

func (t *Tencent) WebName() string {
	return "腾讯视频"
}

func (t *Tencent) Pattern() string {
	return `https?://(?:v\.|www\.)?qq\.(com|cn)/x/(?:cover|page)/(?:[a-z\d]+/)?(?P<id>[a-z\d]+)`
}

func guid() string {
	var uid string
	rand.Seed(time.Now().UnixNano())
	for i := 1; i <= 32; i++ {
		uid += fmt.Sprintf("%x", int(math.Floor(rand.Float64()*16)))
	}
	return uid
}

func time33(t string) int {
	e := 0
	n := len(t)
	i := 5381
	for e < n {
		i += i<<5 + int(rune(t[e]))
		e += 1
	}

	return 2147483647 & i
}
