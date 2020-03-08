package spiders

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	iqiyiApi     = "https://cache.video.iqiyi.com/dash?"
	iqiyiAuthApi = "http://111.59.199.42:9999/authKey?"
	iqiyiVfApi   = "http://111.59.199.42:9999/vf?"
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

type IqiyiBody struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Program struct {
			Video []IqiyiVideoBody `json:"video"`
		} `json:"program"`
	} `json:"data"`
}

type IqiyiVideoBody struct {
	Ff       string `json:"ff"`
	Name     string `json:"name"`
	Vsize    int    `json:"vsize"`
	Duration int    `json:"duration"`
	Bid      int    `json:"bid"`
	M3u8     string `json:"m3u8"`
	Scrsz    string `json:"scrsz"`
}

type userID struct {
	Uid int `json:"uid"`
}

type Iqiyi struct {
	SpiderObject
}

func (i *Iqiyi) Parse(r string) (body Body, err error) {
	bytes, err := i.DownloadWeb(r, url.Values{}, map[string]string{})
	if err != nil {
		return
	}
	re := regexp.MustCompile(`param\['tvid'\] = \"(?P<tvid>\d+)\";\s+param\['vid'\] = "(?P<vid>[a-zA-Z\d]+)"`)
	match := re.FindAllStringSubmatch(string(bytes), -1)
	if len(match) == 0 || len(match[0]) < 3 {
		err = errors.New("parse web page error")
		return
	}

	tvid := match[0][1]
	vid := match[0][2]

	var title string
	re = regexp.MustCompile(`"tvName":"(?P<title>.*)","isfeizhengpian"`)
	match = re.FindAllStringSubmatch(string(bytes), -1)
	if len(match) == 0 || len(match[0]) < 2 {
		title = vid
	} else {
		title = match[0][1]
	}

	tm := time.Now().Unix() * 1000
	dfp := i.Cookies["__dfp"]
	if dfp != "" {
		dfpSlice := strings.Split(dfp, "@")
		if len(dfpSlice) > 0 {
			dfp = dfpSlice[0]
		}
	}

	p, err := url.ParseQuery(i.Cookies["P00002"])
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

	bytes, err = i.DownloadWeb(iqiyiAuthApi, url.Values{"tvid": []string{tvid}, "tm": []string{strconv.Itoa(int(tm))}}, map[string]string{})
	if err != nil {
		return
	}

	for _, q := range iqiyiQualitys {
		params := url.Values{
			"tvid":          []string{tvid},
			"vid":           []string{vid},
			"bid":           []string{strconv.Itoa(q)},
			"src":           []string{"01010031010000000000"},
			"vt":            []string{"0"},
			"rs":            []string{"1"},
			"uid":           []string{uid},
			"ori":           []string{"pcw"},
			"ps":            []string{"0"},
			"k_uid":         []string{i.Cookies["QC005"]},
			"pt":            []string{"0"},
			"d":             []string{"0"},
			"s":             []string{""},
			"lid":           []string{""},
			"cf":            []string{""},
			"ct":            []string{""},
			"authKey":       []string{string(bytes)},
			"k_tag":         []string{"1"},
			"ost":           []string{"0"},
			"ppt":           []string{"0"},
			"dfp":           []string{dfp},
			"locale":        []string{"zh_cn"},
			"prio":          []string{`{"ff":"f4v","code":2}`},
			"pck":           []string{i.Cookies["P00001"]},
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

		var u []byte
		u, err = i.DownloadWeb(iqiyiVfApi,
			url.Values{
				"url": []string{base64.StdEncoding.EncodeToString([]byte(iqiyiApi + params.Encode() + "&ut=1"))},
			}, map[string]string{})

		if err != nil {
			return
		}
		var data IqiyiBody
		err = i.DownloadJson("GET", string(u), url.Values{}, nil, map[string]string{"Referer": r}, &data)
		if err != nil {
			return
		}

		if data.Code != "A00000" {
			err = errors.New(data.Msg)
			return
		}

		var repeated bool
		for _, v := range data.Data.Program.Video {
			for _, b := range body.VideoList {
				if b.ID == v.Bid && v.M3u8 != "" {
					repeated = true
					break
				}
			}

			if repeated {
				break
			}

			if v.M3u8 != "" {
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
				body.VideoList = append(body.VideoList, &VideoBody{
					ID:       v.Bid,
					Title:    title,
					Part:     v.Name,
					Format:   "mp4",
					Size:     v.Vsize,
					Duration: v.Duration,
					Width:    w,
					Height:   h,
					Quality:  iqiyiQuality[v.Bid],
					Links: []VideoAttr{
						VideoAttr{
							URL:   v.M3u8,
							Order: 0,
							Size:  v.Vsize,
						},
					},
					DownloadProtocol: "hlsText",
				})
				break
			}
		}
	}

	return

}

func (i *Iqiyi) Name() string {
	return "iqiyi"
}

func (i *Iqiyi) WebName() string {
	return "爱奇艺【https://www.iqiyi.com/】"
}

func (i *Iqiyi) Pattern() string {
	return `https?://(?:www\.)iqiyi\.com/v_(?P<id>[\da-z]+)`
}
