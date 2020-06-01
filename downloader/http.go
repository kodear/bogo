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
	defer close(cls.status.ch)

	file, err := os.Create(cls.file)
	if err != nil {
		cls.status.Msg = err
		return
	}
	defer func() { _ = file.Close() }()

	res, err := cls.request(cls.urls[0])
	if err != nil {
		cls.status.Msg = err
		return
	}

	// 获取视频大小
	cls.status.MaxLength = cls.length(res)

	// 开始下载
	err = cls.download(res, file)
	if err != nil {
		cls.status.Msg = err
		return
	}

}

func (cls *HTTPFileDownloader) Start() {
	go cls.start()
}
