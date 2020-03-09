package spiders

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bndr/gotabulate"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Safari/537.36"

type Spider interface {
	Parse(string) (Body, error)
	Pattern() string
	SetCookies(string)
	Name() string
	WebName() string
}

var spiders = []Spider{
	&Acfun{}, &AcfunBangumi{}, &BiliBili{}, &BiliBiliBangumi{}, &Iqiyi{}, &Tencent{}, &YouKu{}, &RiJuTv{},
}

type SpiderObject struct {
	Cookies      map[string]string
	CookieString string
}

func HasMatch(r, pattern string) bool {
	match, err := regexp.MatchString(pattern, r)
	if err != nil {
		return false
	} else {
		return match
	}
}

func (s *SpiderObject) Request(method, r string, params url.Values, data []byte, headers map[string]string) (bytes []byte, err error) {
	client := &http.Client{}
	var req *http.Request
	if method == "GET" || method == "" {
		if !strings.HasSuffix(r, "?") && len(params) != 0 {
			r += "?"
		}
		req, err = http.NewRequest("GET", r+params.Encode(), nil)
	} else if method == "POST" {
		req, err = http.NewRequest("POST", r, strings.NewReader(string(data)))
	}

	if err != nil {
		return
	}

	// 设置Headers
	req.Header.Set("User-Agent", UserAgent)
	//if method == "POST" {
	//	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// 设置Cookies
	//for key, value := range s.Cookies {
	//	req.AddCookie(&http.Cookie{Name: key, Value: value})
	//}
	req.Header.Set("Cookie", s.CookieString)

	// 发送请求
	res, err := client.Do(req)
	if err != nil {
		return
	}

	if res.StatusCode != 200 {
		err = errors.New("response code " + strconv.Itoa(res.StatusCode))
		return
	}

	bytes, err = ioutil.ReadAll(res.Body)

	return
}

func (s *SpiderObject) DownloadWeb(r string, params url.Values, headers map[string]string) (bytes []byte, err error) {
	return s.Request("GET", r, params, nil, headers)
}

func (s *SpiderObject) DownloadJson(method, r string, params url.Values, data []byte, headers map[string]string, d interface{}) (err error) {
	bytes, err := s.Request(method, r, params, data, headers)
	if err != nil {
		return
	}
	err = json.Unmarshal(bytes, d)
	return
}

func (s *SpiderObject) SetCookies(cookies string) {
	s.CookieString = cookies
	s.Cookies = SplitCookies(cookies)
}

func (s *SpiderObject) AddCookies(key, value string) {
	s.CookieString += key + "=" + value + ";"
	s.Cookies["key"] = value
}

type Body struct {
	Code      int   //返回值
	Msg       error //返回错误信息
	VideoList []*VideoBody
}

type VideoBody struct {
	ID               int    //视频ID
	Title            string //视频名称
	Part             string //视频集数
	Format           string //视频格式
	Size             int    //视频大小
	Duration         int    //视频长度
	Width            int    //视频宽
	Height           int    //视频高
	StreamType       string //视频类型
	Quality          string // 视频质量
	Links            []VideoAttr
	DownloadHeaders  map[string]string // 下载视频需要的请求头
	DownloadProtocol string            //视频下载协议
}

type VideoAttr struct {
	URL   string //碎片下载地址
	Order int    //碎片视频排序
	Size  int    // 碎片视频大小
}

func do(r string, cookies map[string]string) (body Body, err error) {
	for _, p := range spiders {
		if HasMatch(r, p.Pattern()) {
			p.SetCookies(cookies[p.Name()])
			body, err = p.Parse(r)
			if err != nil {
				return
			}
			if len(body.VideoList) == 0 {
				err = errors.New("unauthorized access")
				return
			}
			sort.Sort(sortVideo{body.VideoList, func(x, y *VideoBody) bool {
				if x.Size != y.Size {
					return x.Size > y.Size
				}
				if x.Height != y.Height {
					return x.Height > y.Height
				}
				if x.Width != y.Width {
					return x.Width > y.Width
				}
				return false
			}})
			return
		}
	}

	err = errors.New("the url is not matched")
	return
}

func ShowVideo(r string, cookies map[string]string) {
	body, err := do(r, cookies)
	if err != nil {
		fmt.Println(err)
	} else {
		var p [][]string
		for _, q := range body.VideoList {
			var part string
			if q.Part == "" {
				part = "-"
			} else {
				part = q.Part
			}

			var streamType string
			if q.StreamType == "" {
				streamType = "-"
			} else {
				streamType = q.StreamType
			}

			var size string
			if q.Size == 0 {
				size = "-"
			} else {
				size = strconv.Itoa(q.Size)
			}

			var duration string
			if q.Duration == 0 {
				duration = "-"
			} else {
				duration = strconv.Itoa(q.Duration)
			}

			var width string
			if q.Width == 0 {
				width = "-"
			} else {
				width = strconv.Itoa(q.Width)
			}

			var height string
			if q.Height == 0 {
				height = "-"
			} else {
				height = strconv.Itoa(q.Height)
			}

			p = append(p, []string{strconv.Itoa(q.ID), q.Title, part, q.Format, streamType, q.Quality, width, height, duration, size, q.DownloadProtocol})
		}

		tabulate := gotabulate.Create(p)
		tabulate.SetHeaders([]string{"ID", "Title", "Part", "Format", "StreamType", "Quality", "Width", "Height", "Duration", "Size", "DownloadProtocol"})
		fmt.Println(tabulate.Render("simple"))
	}
}

func ShowWeb() {
	fmt.Println()
	for _, s := range spiders {
		fmt.Println(s.WebName())
	}
}

func Do(r, q string, id int, cookies map[string]string) (v *VideoBody, err error) {
	body, err := do(r, cookies)
	if err != nil {
		return
	}

	if len(body.VideoList) == 1 {
		return body.VideoList[0], nil
	}

	// 根据ID取数据
	if id != 0 {
		for _, u := range body.VideoList {
			if u.ID == id {
				return u, nil
			}
		}
	}

	// 没有匹配到ID, 根据视频质量取数据
	if q != "" {
		for _, u := range body.VideoList {
			if strings.Contains(u.Quality, q) {
				return u, nil
			}
		}
	}

	// 都没有匹配到, 取默认值, 默认720P
	for _, u := range body.VideoList {
		if strings.Contains(u.Quality, "720") {
			return u, nil
		}
	}

	// 如果默认值也取不到, 则取最大码率数据
	return body.VideoList[0], nil
}

type sortVideo struct {
	v    []*VideoBody
	less func(x, y *VideoBody) bool
}

func (s sortVideo) Len() int {
	return len(s.v)
}

func (s sortVideo) Less(i, j int) bool {
	return s.less(s.v[i], s.v[j])
}

func (s sortVideo) Swap(i, j int) {
	s.v[i], s.v[j] = s.v[j], s.v[i]
}
