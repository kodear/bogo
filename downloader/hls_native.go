package downloader

import (
	"io"
	"strings"
)

type HLSNativeFileDownloader struct {
	HLSFileDownloader
}

func (cls *HLSNativeFileDownloader) Meta() *Meta {
	return &Meta{Name: "hls_native"}
}

func (cls *HLSNativeFileDownloader) start() {
	reader, err := cls.open()
	if err != nil {
		cls.DownloadStatus.Msg = err
	}else{
		cls.run(reader)
	}
}

func (cls *HLSNativeFileDownloader) Start() {
	go cls.start()
}

func (cls *HLSNativeFileDownloader) open()(io.Reader, error)  {
	return strings.NewReader(cls.URL), nil
}
