package download

import (
	"errors"
)

type Downloader interface {
	Do(link, text, file, fname string, links []string)
	Max() int
	SetMax(len int)
	Error() error
	Status() bool
	Chan() chan int
	INIT() bool
	Progress() bool
	SetHeaders(headers map[string]string)
}

type Download struct {
	Headers         map[string]string
	Ch              chan int
	Len             int
	DownloadINIT    bool
	DownloadStatus  bool
	DownloadMessage error
	ProgressINIT    bool
}

func (d *Download) Max() int {
	return d.Len
}

func (d *Download) SetMax(len int) {
	d.Len = len
}

func (d *Download) Chan() chan int {
	return d.Ch
}

func (d *Download) Error() error {
	return d.DownloadMessage
}

func (d *Download) Status() bool {
	return d.DownloadStatus
}

func (d *Download) SetHeaders(headers map[string]string) {
	d.Headers = headers
}

func (d *Download) Progress() bool {
	return d.ProgressINIT
}

func (d *Download) INIT() bool {
	return d.DownloadINIT
}

func LoadDownloader(protocol string) (downloader Downloader, err error) {
	if protocol == "hls" || protocol == "hlsText" || protocol == "hlsFile" {
		downloader = &HLS{}
	} else if protocol == "http" {
		downloader = &HTTP{}
	} else if protocol == "httpSegFlv" {
		downloader = &HTTPSegFLV{}
	} else if protocol == "httpSegF4v" {
		downloader = &HTTPSegF4V{}
	} else {
		err = errors.New("the agreement is not currently supported: " + protocol)
	}

	return
}
