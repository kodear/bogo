package spiders

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	url2 "net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	iqiyiApi     = "https://cache.video.iqiyi.com/dash?"
	iqiyiAuthApi = "http://111.59.199.42:9999/authKey"
	iqiyiVfApi   = "http://111.59.199.42:9999/vf"
	iqiyiF4vApi  = "http://111.59.199.42:9999/f4v"
)

var iqiyiQualitys = []int{100, 200, 300, 500, 600, 610}

var iqiyiQuality = map[int]string{
	100: "auto",
	200: "270P",
	300: "480P",
	500: "720P",
	600: "1080P",
	610: "1080P50",
}

type iqiyiResponse struct {
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
	} `json:"data"`
}

type userID struct {
	Uid int `json:"uid"`
}

type us struct {
	L string `json:"l"`
}

type IqiyiIE struct {
	SpiderIE
}

func (tv *IqiyiIE) Parse(url string) (body []Response, err error) {
	response, err := tv.DownloadWebPage(url, url2.Values{}, map[string]string{})
	if err != nil {
		return
	}

	matchResults, err := response.Search(`param\['tvid'\] = \"(?P<tvid>\d+)\";\s+param\['vid'\] = "(?P<vid>[a-zA-Z\d]+)"`)
	if err != nil {
		return
	}

	if len(matchResults) == 0 || len(matchResults[0]) < 3 {
		err = errors.New("parse web page error")
		return
	}

	tvid := matchResults[0][1]
	vid := matchResults[0][2]

	var title string
	matchResults, _ = response.Search(`"tvName":"(?P<title>.*)","isfeizhengpian"`)
	if len(matchResults) == 0 || len(matchResults[0]) < 2 {
		title = vid
	} else {
		title = matchResults[0][1]
	}

	tm := time.Now().Unix() * 1000
	dfp := tv.Cookies["__dfp"]
	if dfp != "" {
		dfpSlice := strings.Split(dfp, "@")
		if len(dfpSlice) > 0 {
			dfp = dfpSlice[0]
		}
	}

	p, err := url2.ParseQuery(tv.Cookies["P00002"])
	if err != nil {
		return
	}

	var uid string
	if len(p) != 0 {
		for k, _ := range p {
			var u userID
			err = json.Unmarshal([]byte(k), &u)
			if err != nil {
				return
			}
			if u.Uid != 0 {
				uid = strconv.Itoa(u.Uid)
				break
			}

		}

	}

	response, err = tv.DownloadWebPage(iqiyiAuthApi, url2.Values{"tvid": []string{tvid}, "tm": []string{strconv.Itoa(int(tm))}}, map[string]string{})
	if err != nil {
		return
	}

	for _, q := range iqiyiQualitys {
		params := url2.Values{
			"tvid":          []string{tvid},
			"vid":           []string{vid},
			"bid":           []string{strconv.Itoa(q)},
			"src":           []string{"01010031010000000000"},
			"vt":            []string{"0"},
			"rs":            []string{"1"},
			"uid":           []string{uid},
			"ori":           []string{"pcw"},
			"ps":            []string{"0"},
			"k_uid":         []string{tv.Cookies["QC005"]},
			"pt":            []string{"0"},
			"d":             []string{"0"},
			"s":             []string{""},
			"lid":           []string{""},
			"cf":            []string{""},
			"ct":            []string{""},
			"authKey":       []string{response.String},
			"k_tag":         []string{"1"},
			"ost":           []string{"0"},
			"ppt":           []string{"0"},
			"dfp":           []string{dfp},
			"locale":        []string{"zh_cn"},
			"prio":          []string{`{"ff":"f4v","code":2}`},
			"pck":           []string{tv.Cookies["P00001"]},
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
		}

		req, err := tv.DownloadWebPage(iqiyiVfApi,
			url2.Values{
				"url": []string{base64.StdEncoding.EncodeToString([]byte(iqiyiApi + params.Encode() + "&ut=1"))},
			}, map[string]string{})

		if err != nil {
			return nil, err
		}

		var data iqiyiResponse
		req, err = tv.DownloadWebPage(req.String, url2.Values{}, map[string]string{"Referer": url})
		if err != nil {
			return nil, err
		}

		err = req.Json(&data)

		if data.Code != "A00000" {
			err = errors.New(data.Msg)
			return nil, err
		}

		var repeated bool
		for _, v := range data.Data.Program.Video {
			for _, b := range body {
				if b.ID == v.Bid && (v.M3u8 != "" || len(v.Fs) > 0) {
					repeated = true
					break
				}
			}

			if repeated {
				break
			}

			if v.M3u8 != "" || len(v.Fs) > 0 {
				width := "0"
				height := "0"
				if v.Scrsz != "" {
					re := regexp.MustCompile(`(\d+)x(\d+)`)
					match := re.FindAllStringSubmatch(v.Scrsz, -1)
					if len(match) > 0 && len(match[0]) == 3 {
						width = match[0][1]
						height = match[0][2]
					}
				}
				w, _ := strconv.Atoi(width)
				h, _ := strconv.Atoi(height)

				var format string
				var urlAttr []URLAttr
				var downloadProtocol string
				if v.M3u8 != "" {
					format = "mp4"
					downloadProtocol = "hlsText"
					urlAttr = []URLAttr{
						{
							URL:   v.M3u8,
							Order: 0,
							Size:  v.Vsize,
						},
					}
				} else {
					format = "f4v"
					if len(v.Fs) > 1 {
						downloadProtocol = "httpSegF4v"
					} else {
						downloadProtocol = "http"
					}

					for _, u := range v.Fs {
						l := data.Data.Dd + u.L
						l += fmt.Sprintf("&cross-domain=1&qyid=%s&qypid=%s&t=%s&cid=afbe8fd3d73448c9&vid=%s&QY00001=%s&su=%s&client=&z=&mi=%s&bt=&ct=5&e=&ib=4&ptime=0&pv=0.1&tn=%v",
							data.Data.Boss.Data.Pt, tvid+"_02020031010000000000", data.Data.Boss.Data.T, vid, data.Data.Boss.Data.U, data.Data.Boss.Data.Pt,
							fmt.Sprintf("tv_%d_%s_%s", data.Data.Aid, tvid, vid), rand.Float32())

						x := strings.Split(u.L, "?")[0]
						z := strings.Split(x, "/")
						r, err := tv.DownloadWebPage(iqiyiF4vApi, url2.Values{"sign": []string{data.Data.Boss.Data.T + strings.Split(z[len(z)-1], ".")[0]}}, map[string]string{})
						if err != nil {
							return nil, err
						}

						p, err := tv.DownloadWebPage(l+"&ibt="+r.String, url2.Values{}, map[string]string{})
						if err != nil {
							return nil, err
						}

						var d us
						err = p.Json(&d)
						if err != nil {
							return nil, err
						}

						urlAttr = append(urlAttr, URLAttr{
							URL: d.L,
						})
					}

				}

				body = append(body, Response{
					ID:               v.Bid,
					Title:            title,
					Part:             v.Name,
					Format:           format,
					Size:             v.Vsize,
					Duration:         v.Duration,
					Width:            w,
					Height:           h,
					Quality:          iqiyiQuality[v.Bid],
					Links:            urlAttr,
					DownloadProtocol: downloadProtocol,
				})
				break
			}
		}
	}

	return

}

func (tv *IqiyiIE) CookieName() string {
	return "iqiyi"
}

func (tv *IqiyiIE) Name() string {
	return "爱奇艺"
}

func (tv *IqiyiIE) Domain() *Cookie {
	return &Cookie{"iqiyi",true, []string{".iqiyi.com"}}
}

func (tv *IqiyiIE) Pattern() string {
	return `https?://(?:www\.)iqiyi\.com/v_(?P<id>[\da-z]+)`
}
