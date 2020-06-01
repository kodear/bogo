package spider

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zhxingy/bogo/exception"
	"github.com/zhxingy/bogo/selector"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type IQIYIClient struct {
	Client
}

func (cls *IQIYIClient) Meta() *Meta {
	return &Meta{
		Domain:     "https://www.iqiyi.com/",
		Name:       "爱奇艺",
		Expression: `https?://(?:www\.)iqiyi\.com/v_(?P<id>[\da-z]+)`,
		Cookie: Cookie{
			Name:   "iqiyi",
			Enable: true,
			Domain: []string{".iqiyi.com"},
		},
	}
}

func (cls *IQIYIClient) Request() (err error) {
	response, err := cls.request(cls.URL, nil)
	if err != nil {
		return exception.HTTPHtmlException(err)
	}

	var vid, tvid, title string
	err = response.Re(`param\['tvid'\] = \"(?P<tvid>\d+)\";\s+param\['vid'\] = "(?P<vid>[a-zA-Z\d]+)"`, &tvid, &vid)
	if err != nil {
		return exception.HTMLParseException(err)
	}
	err = response.Re(`"tvName":"(?P<title>.*)","isfeizhengpian"`, &title)
	if err != nil {
		return exception.HTMLParseException(err)
	}

	tm := time.Now().Unix() * 1000
	dfp := cls.CookieJar.Name("__dfp")
	if dfp != "" && len(strings.Split(dfp, "@")) > 0 {
		dfp = strings.Split(dfp, "@")[0]
	}

	var uid string
	p2, _ := url.ParseQuery(cls.CookieJar.Name("P00002"))
	if len(p2) > 0 {
		type user struct {
			Uid int `json:"uid"`
		}
		for key := range p2 {
			var u user
			_ = json.Unmarshal([]byte(key), &u)
			if u.Uid != 0 {
				uid = strconv.Itoa(u.Uid)
				break
			}

		}
	}

	key, err := cls.authKey(tvid, strconv.Itoa(int(tm)))
	if err != nil {
		return exception.AuthKeyException(err)
	}

	cls.Header.Add("Referer", cls.URL)
	for _, qualityID := range []int{100, 200, 300, 500, 600, 610} {
		uri, err := cls.cmd5x("https://cache.video.iqiyi.com/dash?" + url.Values{
			"tvid":          []string{tvid},
			"vid":           []string{vid},
			"bid":           []string{strconv.Itoa(qualityID)},
			"src":           []string{"01010031010000000000"},
			"vt":            []string{"0"},
			"rs":            []string{"1"},
			"uid":           []string{uid},
			"ori":           []string{"pcw"},
			"ps":            []string{"0"},
			"k_uid":         []string{cls.CookieJar.Name("QC005")},
			"pt":            []string{"0"},
			"d":             []string{"0"},
			"s":             []string{""},
			"lid":           []string{""},
			"cf":            []string{""},
			"ct":            []string{""},
			"authKey":       []string{key},
			"k_tag":         []string{"1"},
			"ost":           []string{"0"},
			"ppt":           []string{"0"},
			"dfp":           []string{dfp},
			"locale":        []string{"zh_cn"},
			"prio":          []string{`{"ff":"f4v","code":2}`},
			"pck":           []string{cls.CookieJar.Name("P00001")},
			"k_err_retries": []string{"0"},
			"up":            []string{""},
			"qd_v":          []string{"2"},
			"tm":            []string{strconv.Itoa(int(tm))},
			"qdy":           []string{"a"},
			"qds":           []string{"0"},
			"k_ft1":         []string{"143486267424772"},
			"k_ft4":         []string{"1581060"},
			"k_ft5":         []string{"1"},
			"bop":           []string{fmt.Sprintf(`{"version":"10.0","dfp":"{%s}"}`, dfp)},
		}.Encode() + "&ut=1")
		if err != nil {
			return exception.AuthKeyException(err)
		}

		response, err = cls.request(uri, nil)
		if err != nil {
			if len(cls.response) > 0 {
				goto End
			}
			return exception.HTTPJsonException(err)
		}

		var vjson struct {
			Code string `json:"code"`
			Msg  string `json:"msg"`
			Data struct {
				Aid  int    `json:"aid"`
				Dd   string `json:"dd"`
				Boss struct {
					Data struct {
						Pt string `json:"pt"`
						T  string `json:"t"`
						U  string `json:"u"`
					} `json:"data"`
				} `json:"boss"`
				Program struct {
					Video []struct {
						Ff       string `json:"ff"`
						Name     string `json:"name"`
						Vsize    int    `json:"vsize"`
						Duration int    `json:"duration"`
						Bid      int    `json:"bid"`
						M3u8     string `json:"m3u8"`
						Scrsz    string `json:"scrsz"`
						Fs       []struct {
							L string `json:"l"`
						} `json:"fs"`
					} `json:"video"`
				} `json:"program"`
				BossTs struct {
					Msg  string `json:"msg"`
					Code string `json:"code"`
				} `json:"boss_ts"`
			} `json:"data"`
		}

		err = response.Json(&vjson)
		if err != nil {
			return exception.JSONParseException(err)
		}

		if vjson.Code != "A00000" || len(vjson.Data.Program.Video) == 0 {
			if len(cls.response) > 0 {
				goto End
			}
			return exception.ServerAuthException(errors.New(vjson.Data.BossTs.Msg))
		}

		for _, video := range vjson.Data.Program.Video {
			if video.M3u8 != "" || len(video.Fs) > 0 {
				var width, height int
				var w, h string
				if video.Scrsz != "" {
					var x selector.Selector
					x = []byte(video.Scrsz)
					_ = x.Re(`(\d+)x(\d+)`, &w, &h)
				}
				if w != "" && h != "" {
					width, _ = strconv.Atoi(w)
					height, _ = strconv.Atoi(h)
				}

				var format, protocol string
				var urlAttrs []URLAttr
				if video.M3u8 != "" {
					format = "ts"
					protocol = "hls_native"
					urlAttrs = []URLAttr{
						{
							URL:   video.M3u8,
							Order: 0,
							Size:  video.Vsize,
						},
					}
				} else {
					format = "f4v"
					if len(video.Fs) > 1 {
						protocol = "f4v"
					} else {
						protocol = "http"
					}

					for _, v := range video.Fs {
						uri = vjson.Data.Dd + v.L + fmt.Sprintf("&cross-domain=1&qyid=%s&qypid=%s&t=%s&cid=afbe8fd3d73448c9&vid=%s&QY00001=%s&su=%s&client=&z=&mi=%s&bt=&ct=5&e=&ib=4&ptime=0&pv=0.1&tn=%v", vjson.Data.Boss.Data.Pt, tvid+"_02020031010000000000", vjson.Data.Boss.Data.T, vid, vjson.Data.Boss.Data.U, vjson.Data.Boss.Data.Pt, fmt.Sprintf("tv_%d_%s_%s", vjson.Data.Aid, tvid, vid), rand.Float32())
						z := strings.Split(strings.Split(v.L, "?")[0], "/")
						sign, err := cls.f4v(vjson.Data.Boss.Data.T + strings.Split(z[len(z)-1], ".")[0])
						if err != nil {
							return exception.AuthKeyException(err)
						}

						response, err := cls.request(uri+"&ibt="+sign, nil)
						if err != nil {
							return exception.HTTPJsonException(err)
						}

						var us struct {
							L string `json:"l"`
						}
						err = response.Json(&us)
						if err != nil {
							return exception.JSONParseException(err)
						}

						urlAttrs = append(urlAttrs, URLAttr{
							URL: us.L,
						})
					}
				}

				cls.response = append(cls.response, &Response{
					ID:       video.Bid,
					Title:    title,
					Part:     video.Name,
					Format:   format,
					Size:     video.Vsize,
					Duration: video.Duration,
					Width:    width,
					Height:   height,
					Quality: map[int]string{
						100: "auto",
						200: "270P",
						300: "480P",
						500: "720P",
						600: "1080P",
						610: "1080P50",
					}[video.Bid],
					Links:            urlAttrs,
					DownloadProtocol: protocol,
				})

				break
			}
		}
	}

	if len(cls.response) < 1 {
		return exception.OtherException(errors.New(""))
	}

End:
	x := make(map[int]*Response)
	for _, response := range cls.response {
		x[response.ID] = response
	}

	var Response []*Response
	for _, response := range x {
		Response = append(Response, response)
	}

	cls.response = Response
	return

}

func (cls *IQIYIClient) authKey(id, tm string) (key string, err error) {
	response, err := cls.request("http://111.59.199.42:9999/authKey?", url.Values{
		"tvid": []string{id},
		"tm":   []string{tm},
	})
	if err != nil {
		return
	}

	key = response.String()
	return
}

func (cls *IQIYIClient) cmd5x(oldUrl string) (newUrl string, err error) {
	response, err := cls.request("http://111.59.199.42:9999/vf?", url.Values{"url": []string{base64.StdEncoding.EncodeToString([]byte(oldUrl))}})
	if err != nil {
		return
	}

	newUrl = response.String()
	return
}

func (cls *IQIYIClient) f4v(key string) (sign string, err error) {
	response, err := cls.request("http://111.59.199.42:9999/f4v", url.Values{"sign": []string{key}})
	if err != nil {
		return
	}

	sign = response.String()
	return
}
