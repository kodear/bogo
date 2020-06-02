package spider

type HLSClient struct {
	Client
}

func (cls *HLSClient) Meta() *Meta {
	return &Meta{
		Domain:     "",
		Name:       "hls",
		Expression: `https?://.*\.m3u8(?:\?.*)?$`,
		Cookie:     Cookie{},
	}
}

func (cls *HLSClient) Request() (err error) {
	cls.response = &Response{
		Title: md5x(cls.URL),
		Stream: []Stream{
			{
				ID:     1,
				URLS:   []string{cls.URL},
				Format: "ts",
			},
		},
	}

	return
}

type HTTPClient struct {
	Client
}

func (cls *HTTPClient) Meta() *Meta {
	return &Meta{
		Domain:     "",
		Name:       "http",
		Expression: `https?://.*\.(mp4|flv|f4v|ts)(?:\?.*)?$`,
		Cookie:     Cookie{},
	}
}

func (cls *HTTPClient) Request() (err error) {
	var selector Selector
	var format string
	selector = []byte(cls.URL)
	err = selector.Re(cls.Meta().Expression, &format)
	if err != nil {
		return
	}

	cls.response = &Response{
		Title: md5x(cls.URL),
		Stream: []Stream{
			{
				ID:     1,
				URLS:   []string{cls.URL},
				Format: format,
			},
		},
	}

	return nil
}
