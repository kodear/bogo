package spider

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
)

const UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Safari/537.36"

type Stream struct {
	ID               int    //视频ID
	Format           string //视频格式
	Size             int    //视频大小
	Width            int    //视频宽
	Height           int    //视频高
	StreamType       string //视频类型
	Quality          string // 视频质量
	URLS             []string
	DownloadHeaders  http.Header // 下载视频需要的请求头
	DownloadProtocol string      //视频下载协议
}

type Response struct {
	Title    string //视频名称
	Part     string //视频集数
	Duration int    //视频时长
	Site     string
	Stream   []Stream
}

type Client struct {
	URL       string
	Proxy     string
	Header    http.Header
	CookieJar CookiesJar
	response  *Response
}

type Meta struct {
	Domain     string // 网站首页
	Name       string // 项目名
	Expression string
	Cookie     Cookie
}

func (cls *Client) request(uri string, params url.Values) (selector Selector, err error) {
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

func (cls *Client) fromRequest(uri string, params url.Values, data []byte) (selector Selector, err error) {
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

func (cls *Client) Meta() *Meta {
	panic("this method must be implemented by subclasses")
}

func (cls *Client) Request() (err error) {
	panic("this method must be implemented by subclasses")
}

func (cls *Client) Response() *Response {
	return cls.response

}

func (cls *Client) Initialization(jar CookiesJar, header http.Header) {
	if header == nil {
		header = http.Header{}
	}
	cls.Header = header
	cls.CookieJar = jar
}

func (cls *Client) SetURL(url string) {
	cls.URL = url
}

type Cookie struct {
	Name   string
	Enable bool
	Domain []string
}

type CookiesJar []*http.Cookie

func (cls *CookiesJar) Name(key string) (value string) {
	for _, cookie := range *cls {
		if cookie.Name == key {
			value = cookie.Value
			break
		}
	}
	return value
}

func (cls *CookiesJar) SetValue(key, value string) {
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

func (cls *CookiesJar) String() (cookies string) {
	for _, cookie := range *cls {
		cookies += cookie.Name + "=" + cookie.Value + ";"
	}
	return
}
