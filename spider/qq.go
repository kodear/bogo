package spider

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type QQClient struct {
	Client
}

func (cls *QQClient) Meta() *Meta {
	return &Meta{
		Domain:     "https://v.qq.com/",
		Name:       "腾讯视频",
		Expression: `https?://(?:v\.|www\.)?qq\.(com|cn)/x/(?:cover|page)/(?:[a-z\d]+/)?(?P<id>[a-z\d]+)`,
		Cookie: Cookie{
			Name:   "qq",
			Enable: true,
			Domain: []string{".v.qq.com"},
		},
	}
}

func (cls *QQClient) Request() (err error) {
	response, err := cls.request(cls.URL, nil)
	if err != nil {
		return DownloadHtmlErr(err)
	}

	var vid string
	err = response.Re(`vid=(?P<vid>[a-zA-Z\d]+)&`, &vid)
	if err != nil {
		return ParseJsonErr(err)
	}

	tm := int(time.Now().Unix())
	uid := cls.guid()

	cresponse, err := cls.request("http://111.59.199.42:9999/ckey?", url.Values{
		"vid":  []string{vid},
		"guid": []string{uid},
		"tm":   []string{strconv.Itoa(tm)},
	})
	if err != nil {
		return ServerAuthKeyErr(err)
	}
	key := cresponse.String()

	cls.Header.Add("Referer", cls.URL)
	response, err = cls.request("https://access.video.qq.com/user/auth_refresh?", url.Values{
		"vappid":  []string{"11059694"},
		"vsecret": []string{"fdf61a6be0aad57132bc5cdf78ac30145b6cd2c1470b0cfe"},
		"type":    []string{"qq"},
		"g_tk":    []string{},
		"g_vstk":  []string{strconv.Itoa(cls.sign(cls.CookieJar.Name("vqq_vusession")))},
		"g_actk":  []string{strconv.Itoa(cls.sign(cls.CookieJar.Name("vqq_access_token")))},
		"_":       []string{strconv.Itoa(tm * 1000)},
	})
	if err == nil {
		var auth struct {
			Errcode     int    `json:"errcode"`
			AccessToken string `json:"access_token"`
			Vuserid     int    `json:"vuserid"`
			Vusession   string `json:"vusession"`
		}

		err = response.ReByJson(`data=(.*)`, &auth)
		//if auth.AccessToken == "" || auth.Vusession == "" || auth.Vuserid == 0{
		//	return exception.ServerAuthException(errors.New("qq video auth error."))
		//}
		if err != nil && auth.Errcode == 0 {
			cls.CookieJar.SetValue("vqq_vusession", auth.Vusession)
			cls.CookieJar.SetValue("vqq_access_token", auth.AccessToken)
			cls.CookieJar.SetValue("vqq_vuserid", strconv.Itoa(auth.Vuserid))
		}
	}

	type postJson struct {
		Buid       string `json:"buid"`
		Vinfoparam string `json:"vinfoparam"`
	}

	cls.CookieJar = nil
	cls.response = &Response{
		Site:   cls.Meta().Name,
		Stream: []Stream{},
	}
	for _, quality := range []string{"sd", "hd", "shd", "fhd"} {
		body, _ := json.Marshal(postJson{
			"onlyvinfo",
			url.Values{
				"adsid":       []string{""},
				"vid":         []string{vid},
				"dtype":       []string{"3"},
				"show1080p":   []string{"1"},
				"guid":        []string{uid},
				"sdtfrom":     []string{"v1010"},
				"adpinfo":     []string{""},
				"spgzip":      []string{"1"},
				"spau":        []string{"1"},
				"ehost":       []string{cls.URL},
				"fp2p":        []string{"1"},
				"spaudio":     []string{"15"},
				"appVer":      []string{"3.5.57"},
				"sphttps":     []string{"1"},
				"platform":    []string{"10201"},
				"charge":      []string{"0"},
				"tm":          []string{strconv.Itoa(tm)},
				"cKey":        []string{key},
				"sphls":       []string{"2"},
				"defaultfmt":  []string{"auto"},
				"refer":       []string{"v.qq.com"},
				"spwm":        []string{"4"},
				"flowid":      []string{uid + "_10201"},
				"host":        []string{"v.qq.com"},
				"hdcp":        []string{"hdcp"},
				"defn":        []string{quality},
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
			}.Encode(),
		})
		response, err := cls.fromRequest("https://vd.l.qq.com/proxyhttp", nil, body)
		if err != nil {
			return DownloadJsonErr(err)
		}

		var vjson struct {
			ErrCode int    `json:"errCode"`
			Vinfo   string `json:"vinfo"`
		}
		err = response.Json(&vjson)
		if err != nil {
			return ParseJsonErr(err)
		}
		if vjson.ErrCode != 0 && len(cls.response.Stream) > 0 {
			return nil
		} else if vjson.ErrCode != 0 {
			return ServerAuthErr(errors.New("qq video code: " + strconv.Itoa(vjson.ErrCode)))
		}

		var video struct {
			Msg string `json:"msg"`
			Fl  struct {
				Fi []struct {
					ID         int    `json:"id"`
					Name       string `json:"name"`
					Resolution string `json:"resolution"`
					Fs         int    `json:"fs"`
				} `json:"fi"`
			} `json:"fl"`
			Vl struct {
				Vi []struct {
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
				} `json:"vi"`
			} `json:"vl"`
		}
		err = json.Unmarshal([]byte(vjson.Vinfo), &video)
		if err != nil {
			return ParseJsonErr(err)
		}
		if video.Msg != "" && len(cls.response.Stream) > 0 {
			return nil
		} else if video.Msg != "" {
			return ServerAuthErr(errors.New(video.Msg))
		}

		var id int
		var streamType, format, link string
		duration, _ := strconv.ParseFloat(video.Vl.Vi[0].Td, 32)
		for _, p := range video.Fl.Fi {
			if video.Vl.Vi[0].Fs == p.Fs {
				id = p.ID
				streamType = p.Name
				quality = p.Resolution
				break
			}
		}

		if strings.HasSuffix(video.Vl.Vi[0].Ul.Ui[0].URL, "/") {
			if video.Vl.Vi[0].Ul.Ui[0].Hls.Pt == "" {
				continue
			}
			format = video.Vl.Vi[0].Ul.Ui[0].Hls.Ftype
			link = video.Vl.Vi[0].Ul.Ui[0].URL + video.Vl.Vi[0].Ul.Ui[0].Hls.Pt
		} else {
			format = "ts"
			link = video.Vl.Vi[0].Ul.Ui[0].URL
		}

		cls.response.Title = video.Vl.Vi[0].Ti
		cls.response.Duration = int(duration)
		cls.response.Stream = append(cls.response.Stream, Stream{
			ID:               id,
			Format:           format,
			Size:             video.Vl.Vi[0].Fs,
			Width:            video.Vl.Vi[0].Vw,
			Height:           video.Vl.Vi[0].Vh,
			StreamType:       streamType,
			Quality:          quality,
			URLS:             []string{link},
			DownloadProtocol: "hls",
		})
	}

	return
}

func (cls *QQClient) guid() string {
	var uid string
	rand.Seed(time.Now().UnixNano())
	for i := 1; i <= 32; i++ {
		uid += fmt.Sprintf("%x", int(math.Floor(rand.Float64()*16)))
	}
	return uid
}

func (cls *QQClient) sign(t string) int {
	e := 0
	n := len(t)
	i := 5381
	for e < n {
		i += i<<5 + int(rune(t[e]))
		e += 1
	}

	return 2147483647 & i
}
