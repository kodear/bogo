package downloader

import (
	"os"
)

type HTTPFileDownloader struct {
	FileDownloader
}

func (cls *HTTPFileDownloader) Meta() *Meta {
	return &Meta{Name: "http"}
}

func (cls *HTTPFileDownloader) start() {
	file, err := os.Create(cls.file)
	if err != nil {
		cls.status.Msg = err
		return
	}
	err = cls.request(cls.urls[0], file)
	if err != nil {
		cls.status.Msg = err
	}
	_ = file.Close()
	close(cls.status.ch)
}
