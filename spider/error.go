package spider

import "fmt"

const (
	statusDownloadHtml  = 100
	statusDownloadJson  = 101
	statusDownloadXml   = 102
	statusDownloadText  = 103
	statusParseHtml     = 200
	statusParseJson     = 201
	statusParseXml      = 202
	statusParseText     = 203
	statusServerAuth    = 300
	statusServerAuthKey = 301
	statusUnknown       = 500
)

func DownloadHtmlErr(err error) error {
	return fmt.Errorf("request html error. errcode=%d, errmsg=%s", statusDownloadHtml, err)
}

func DownloadJsonErr(err error) error {
	return fmt.Errorf("request json error. errcode=%d, errmsg=%s", statusDownloadJson, err)
}

func DownloadXmlErr(err error) error {
	return fmt.Errorf("request xml error. errcode=%d, errmsg=%s", statusDownloadXml, err)
}

func DownloadTextErr(err error) error {
	return fmt.Errorf("request text error. errcode=%d, errmsg=%s", statusDownloadText, err)
}

func ParseHtmlErr(err error) error {
	return fmt.Errorf("parse html error. errcode=%d, errmsg=%s", statusParseHtml, err)
}

func ParseJsonErr(err error) error {
	return fmt.Errorf("parse json error. errcode=%d, errmsg=%s", statusParseJson, err)
}

func ParseXmlErr(err error) error {
	return fmt.Errorf("parse xml error. errcode=%d, errmsg=%s", statusParseXml, err)
}

func ParseTextErr(err error) error {
	return fmt.Errorf("parse text error. errcode=%d, errmsg=%s", statusParseText, err)
}

func ServerAuthErr(err error) error {
	return fmt.Errorf("server auth error. errcode=%d, errmsg=%s", statusServerAuth, err)
}

func ServerAuthKeyErr(err error) error {
	return fmt.Errorf("get auth key error. errcode=%d, errmsg=%s", statusServerAuthKey, err)
}

func UnknownErr(err error) error {
	return fmt.Errorf("other error. errcode=%d, errmsg=%s", statusUnknown, err)
}
