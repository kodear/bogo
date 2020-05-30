package spiders

import (
	"bytes"
	"github.com/zhxingy/bogo/selector"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
)

type SpiderResponse struct {
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

type SpiderRequest struct {
	URL       string
	Proxy     string
	Header    http.Header
	CookieJar SpiderCookiesJar
	Response  []*SpiderResponse
}

type SpiderArgs struct {
	Domain string // 网站首页
	Name   string // 项目名
	Cookie struct {
		Name   string   //  cookie key(用于配置文件)
		Enable bool     // 是否启用cookie
		Domain []string // cookie域
	}
}

func (cls *SpiderRequest) request(uri string, params url.Values) (selector selector.Selector, err error) {
	client := &http.Client{}
	req, err := http.NewRequest("", uri+params.Encode(), nil)
	if err != nil {
		return
	}

	// 构造请求头
	req.Header = cls.Header
	req.Header.Set("User-Agent", UserAgent)
	// 构造cookie
	for _, cookie := range cls.CookieJar {
		req.AddCookie(cookie)
	}

	res, err := client.Do(req)
	if err != nil {
		return
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	defer func() { _ = res.Body.Close() }()

	selector = body
	return
}

func (cls *SpiderRequest) fromRequest(uri string, params url.Values, data []byte) (selector selector.Selector, err error) {
	client := &http.Client{}

	req, err := http.NewRequest("POST", uri+params.Encode(), bytes.NewReader(data))
	if err != nil {
		return
	}

	// 构造请求头
	req.Header = cls.Header
	req.Header.Set("User-Agent", UserAgent)
	// 构造cookie
	for _, cookie := range cls.CookieJar {
		req.AddCookie(cookie)
	}
	res, err := client.Do(req)
	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	defer func() { _ = res.Body.Close() }()

	selector = body
	return
}

func (cls *SpiderRequest) Expression() string {
	panic("you have to rewrite the method")
}

func (cls *SpiderRequest) Request() (err error) {
	panic("you have to rewrite the method")
}

func (cls *SpiderRequest) Args() *SpiderArgs {
	panic("you have to rewrite the method")
}

func Match(uri, expression string) bool {
	ok, err := regexp.MatchString(expression, uri)
	if err != nil {
		return false
	} else {
		return ok
	}
}

type SpiderCookiesJar []*http.Cookie

func (cls *SpiderCookiesJar) Name(key string) (value string) {
	for _, cookie := range *cls {
		if cookie.Name == key {
			value = cookie.Value
			break
		}
	}
	return value
}

func (cls *SpiderCookiesJar) SetValue(key, value string) {
	var ok bool
	for _, cookie := range *cls {
		if cookie.Name == key {
			cookie.Value = value
			ok = true
			break
		}
	}

	if !ok {
		*cls = append(*cls, &http.Cookie{
			Name:  key,
			Value: value,
		})
	}

}

func (cls *SpiderCookiesJar) String ()(cookies string){
	for _, cookie := range *cls {
		cookies += cookie.Name + "=" + cookie.Value + ";"
	}
	return
}
