package exception

import "fmt"

const (
	HTTPHtmlExceptionCode   = 100
	HTTPJsonExceptionCode   = 101
	HTTPXmlExceptionCode    = 102
	HTTPTextExceptionCode   = 103
	HTMLParseExceptionCode  = 200
	JSONParseExceptionCode  = 201
	XMLParseExceptionCode   = 202
	TextParseExceptionCode  = 203
	ServerAuthExceptionCode = 300
	AuthKeyExceptionCode    = 301
	OtherExceptionCode      = 500
)

func HTTPHtmlException(err error) error {
	return fmt.Errorf("request html error. errcode=%d, errmsg=%s", HTTPHtmlExceptionCode, err)
}

func HTTPJsonException(err error) error {
	return fmt.Errorf("request json error. errcode=%d, errmsg=%s", HTTPJsonExceptionCode, err)
}

func HTTPXmlException(err error) error {
	return fmt.Errorf("request xml error. errcode=%d, errmsg=%s", HTTPXmlExceptionCode, err)
}

func HTTPTextException(err error) error {
	return fmt.Errorf("request text error. errcode=%d, errmsg=%s", HTTPTextExceptionCode, err)
}

func HTMLParseException(err error) error {
	return fmt.Errorf("parse html error. errcode=%d, errmsg=%s", HTMLParseExceptionCode, err)
}

func JSONParseException(err error) error {
	return fmt.Errorf("parse json error. errcode=%d, errmsg=%s", JSONParseExceptionCode, err)
}

func XMLParseException(err error) error {
	return fmt.Errorf("parse xml error. errcode=%d, errmsg=%s", XMLParseExceptionCode, err)
}

func TextParseException(err error) error {
	return fmt.Errorf("parse text error. errcode=%d, errmsg=%s", TextParseExceptionCode, err)
}

func ServerAuthException(err error) error {
	return fmt.Errorf("server auth error. errcode=%d, errmsg=%s", ServerAuthExceptionCode, err)
}

func AuthKeyException(err error) error {
	return fmt.Errorf("get auth key error. errcode=%d, errmsg=%s", AuthKeyExceptionCode, err)
}

func OtherException(err error) error {
	return fmt.Errorf("other error. errcode=%d, errmsg=%s", OtherExceptionCode, err)
}
