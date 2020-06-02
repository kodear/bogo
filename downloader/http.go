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
	defer close(cls.DownloadStatus.ch)

	file, err := os.Create(cls.File)
	if err != nil {
		cls.DownloadStatus.Msg = err
		return
	}
	defer func() { _ = file.Close() }()

	res, err := cls.request(cls.URL)
	if err != nil {
		cls.DownloadStatus.Msg = err
		return
	}

	// 获取视频大小
	cls.DownloadStatus.MaxLength = cls.length(res)

	// 开始下载
	err = cls.download(res, file)
	if err != nil {
		cls.DownloadStatus.Msg = err
		return
	}

}

func (cls *HTTPFileDownloader) Start() {
	go cls.start()
}

