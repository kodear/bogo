package spiders

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bndr/gotabulate"
	"io/ioutil"
	"net/http"
	url2 "net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Safari/537.36"

var spiders = []Spiders{
	&AcfunIE{},
	&AcfunBangumiIE{},
	&BiliBiliIE{},
	&BiliBiliBangumiIE{},
	&IqiyiIE{},
	&TencentIE{},
	&MgtvIE{},
	&YouKuIE{},
	&RiJuTvIE{},
	&YuespIE{},
}

type Spiders interface {
	Parse(string) ([]Response, error)
	SetCookies(string)
	Pattern() string
	CookieName() string
	Name() string
}

type Response struct {
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
	Links            []URLAttr
	DownloadHeaders  map[string]string // 下载视频需要的请求头
	DownloadProtocol string            //视频下载协议
}

type URLAttr struct {
	URL   string //碎片下载地址
	Order int    //碎片视频排序
	Size  int    // 碎片视频大小
}

type Spider struct {
	Cookies      map[string]string
	CookieString string
}

func (tv *Spider) Request(method, url string, params url2.Values, data []byte, headers map[string]string) (bytes []byte, err error) {

	client := &http.Client{}

	var req *http.Request

	if method == "GET" || method == "" {
		if !strings.HasSuffix(url, "?") && len(params) != 0 {
			url += "?"
		}
		req, err = http.NewRequest("GET", url+params.Encode(), nil)
	} else if method == "POST" {
		req, err = http.NewRequest("POST", url, strings.NewReader(string(data)))
	}

	if err != nil {
		return
	}

	// 设置Headers
	req.Header.Set("User-Agent", UserAgent)

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// 设置Cookies
	req.Header.Set("Cookie", tv.CookieString)

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
	if err == nil {
		_ = res.Body.Close()
	}

	return
}

// 下载网页
func (tv *Spider) DownloadWebPage(url string, params url2.Values, headers map[string]string) (parse *Parse, err error) {
	bytes, err := tv.Request("GET", url, params, nil, headers)
	if err != nil {
		return
	}

	parse = NewParse(bytes)
	return
}

func (tv *Spider) SetCookies(cookies string) {
	tv.CookieString = cookies
	tv.Cookies = splitCookies(cookies)
}

type Parse struct {
	Bytes  []byte
	String string
}

func NewParse(b []byte) *Parse {
	return &Parse{
		Bytes:  b,
		String: string(b),
	}
}

func (p *Parse) Json(data interface{}) (err error) {
	err = json.Unmarshal(p.Bytes, &data)
	return
}

func (p *Parse) Search(pattern string) (strings [][]string, err error) {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return
	}

	strings = regex.FindAllStringSubmatch(p.String, -1)
	return
}

func do(r string, cookies map[string]string) (body []Response, err error) {
	for _, p := range spiders {
		if hasMatch(r, p.Pattern()) {
			p.SetCookies(cookies[p.CookieName()])
			body, err = p.Parse(r)
			if err != nil {
				return
			}
			if len(body) == 0 {
				err = errors.New("unauthorized access")
				return
			}
			sort.Sort(Sort{body, func(x, y Response) bool {
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
		for _, q := range body {
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
		fmt.Println(s.Name())
	}
}

func Do(r, q string, id int, cookies map[string]string) (v Response, err error) {
	body, err := do(r, cookies)
	if err != nil {
		return
	}

	if len(body) == 1 {
		return body[0], nil
	}

	// 根据ID取数据
	if id != 0 {
		for _, u := range body {
			if u.ID == id {
				return u, nil
			}
		}
	}

	// 没有匹配到ID, 根据视频质量取数据
	if q != "" {
		for _, u := range body {
			if strings.Contains(u.Quality, q) {
				return u, nil
			}
		}
	}

	// 都没有匹配到, 取默认值, 默认720P
	for _, u := range body {
		if strings.Contains(u.Quality, "720") {
			return u, nil
		}
	}

	// 如果默认值也取不到, 则取最大码率数据
	return body[0], nil
}
