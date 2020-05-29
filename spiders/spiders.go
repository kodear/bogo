package spiders

import (
	"bytes"
	"encoding/json"
	"github.com/zhxingy/bogo/exception"
	"io/ioutil"
	"net/http"
	"net/url"
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
	URL     string
	Proxy  	string
	Header http.Header
	CookieJar []*http.Cookie
	Response  *SpiderResponse
}

func (cls *SpiderRequest)request(uri string, params url.Values)(body []byte, err error){
	client := &http.Client{}
	req, err := http.NewRequest("", uri + params.Encode(), nil)
	if err != nil{
		return
	}
	req = cls.setRequest(req)

	res, err := client.Do(req)
	if res != nil{
		return
	}

	body, err = ioutil.ReadAll(req.Body)
	if err != nil{
		return
	}

	return
}

func (cls *SpiderRequest)fromRequest(uri string, params url.Values, data []byte)(body []byte, err error){
	client := &http.Client{}
	req, err := http.NewRequest("POST", uri + params.Encode(), bytes.NewBuffer(data))
	if err != nil{
		return
	}
	req = cls.setRequest(req)

	res, err := client.Do(req)
	if res != nil{
		return
	}

	body, err = ioutil.ReadAll(req.Body)
	if err != nil{
		return
	}

	return
}

func (cls *SpiderRequest)setRequest(req *http.Request) *http.Request{
	// 构造请求头
	req.Header = cls.Header
	req.Header.Set("User-Agent", UserAgent)
	// 构造cookie
	for _, cookie := range cls.CookieJar{
		req.AddCookie(cookie)
	}
	return req
}

func (cls *SpiderRequest)downloadWeb(uri string, params url.Values)error{
	_, err := cls.request(uri, params)
	if err != nil{
		return exception.HTTPHtmlException(err)
	}
	return nil
}

func (cls *SpiderRequest)downloadJson(uri string, params url.Values, v interface{})error{
	body, err := cls.request(uri, params)
	if err != nil{
		return exception.HTTPJsonException(err)
	}

	err = json.Unmarshal(body, v)
	if err != nil{
		return exception.JSONParseException(err)
	}

	return nil
}

func (cls *SpiderRequest)fromDownloadJson(uri string, params url.Values, data []byte, v interface{})error{
	body, err := cls.fromRequest(uri, params, data)
	if err != nil{
		return exception.HTTPJsonException(err)
	}

	err = json.Unmarshal(body, v)
	if err != nil{
		return exception.JSONParseException(err)
	}

	return nil
}