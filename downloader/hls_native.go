package downloader

import (
	"io/ioutil"
	"os"
	"strings"
)

type HLSNativeFileDownloader struct {
	HLSFileDownloader
}

func (cls *HLSNativeFileDownloader) Meta() *Meta {
	return &Meta{Name: "hls_native"}
}

func (cls *HLSNativeFileDownloader) start() {
	defer close(cls.status.ch)

	file, err := os.Create(cls.file)
	if err != nil {
		cls.status.Msg = err
		return
	}
	defer func() { _ = file.Close() }()

	playlist, err := cls.parse(strings.NewReader(cls.urls[0]))
	if err != nil {
		return
	}

	cls.status.MaxLength = len(playlist.Segments)

	var key []byte
	if playlist.Key != nil {
		res, err := cls.request(playlist.Key.URI)
		if err != nil {
			return
		}
		key, err = ioutil.ReadAll(res.Body)
		if err != nil {
			return
		}
		_ = res.Body.Close()
	}

	for _, segment := range playlist.Segments {
		if segment == nil {
			break
		}

		res, err := cls.request(urlJoin(cls.urls[0], segment.URI))
		if err != nil {
			return
		}

		err = cls.download(res, file, key)
		if err != nil {
			return
		}
	}

	return

}

func (cls *HLSNativeFileDownloader) Start() {
	go cls.start()
}
