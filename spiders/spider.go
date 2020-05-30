package spiders

import (
	"bytes"
	"github.com/zhxingy/bogo/selector"
	"io/ioutil"
	"net/http"
	"net/url"
)

const UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Safari/537.36"

type URLAttr struct {
	URL   string //碎片下载地址
	Order int    //碎片视频排序
	Size  int    // 碎片视频大小
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

type Client struct {
	URL       string
	Proxy     string
	Header    http.Header
	CookieJar CookiesJar
	response  []*Response
}

type Args struct {
	Domain string // 网站首页
	Name   string // 项目名
	Cookie Cookie
}

func (cls *Client) request(uri string, params url.Values) (selector selector.Selector, err error) {
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

func (cls *Client) fromRequest(uri string, params url.Values, data []byte) (selector selector.Selector, err error) {
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

func (cls *Client) Expression() string {
	panic("you have to rewrite the method")
}

func (cls *Client) Request() (err error) {
	panic("you have to rewrite the method")
}

func (cls *Client) Response() []*Response{
	return cls.response

}

func (cls *Client) Args() *Args {
	panic("you have to rewrite the method")
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

func (cls *CookiesJar) String ()(cookies string){
	for _, cookie := range *cls {
		cookies += cookie.Name + "=" + cookie.Value + ";"
	}
	return
}

type Spider interface {
	Request()error
	Response() []*Response
	Expression() string
}
