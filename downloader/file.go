package downloader

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

type DownloadStatus struct {
	Code      int
	Msg       error
	Length    int
	MaxLength int
	Byte      int
	CH        chan int
	OK        bool
}

type Meta struct {
	Name string
}

type FileDownloader struct {
	URL            string
	URLS           []string
	File           string
	Header         http.Header
	DownloadStatus DownloadStatus
}

func (cls *FileDownloader) Initialize(filename string, urls []string, header http.Header) {
	cls.File = filename
	cls.URLS = urls
	cls.URL = urls[0]
	cls.Header = header
	cls.DownloadStatus = DownloadStatus{
		Code:   0,
		Msg:    nil,
		Length: 0,
		CH:     make(chan int, 1000),
	}
}

func (cls *FileDownloader) Meta() *Meta {
	panic("this method must be implemented by subclasses")
}

func (cls *FileDownloader) Status() DownloadStatus {
	return cls.DownloadStatus
}

func (cls *FileDownloader) Wait() (err error) {
	for {
		if _, ok := <-cls.DownloadStatus.CH; !ok {
			break
		}
	}

	return cls.DownloadStatus.Msg
}

func (cls *FileDownloader) request(uri string) (res *http.Response, err error) {
	client := &http.Client{}
	req, err := http.NewRequest("", uri, nil)
	if err != nil {
		return
	}

	req.Header = cls.Header
	res, err = client.Do(req)
	if err != nil {
		return
	}

	if res.StatusCode/100 > 3 {
		body, _ := ioutil.ReadAll(res.Body)
		_ = res.Body.Close()
		return nil, fmt.Errorf(string(body))
	}

	return
}

func (cls *FileDownloader) length(res *http.Response) int {
	l, _ := strconv.Atoi(res.Header["Content-Length"][0])
	return l
}

func (cls *FileDownloader) download(res *http.Response, file *os.File) (err error) {
	defer func() { _ = res.Body.Close() }()
	buf := make([]byte, 4096)
	for {
		n, err := res.Body.Read(buf)
		if n > 0 {
			_, _ = file.Write(buf[:n])
		} else if err != nil && err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		cls.DownloadStatus.Byte += n
		cls.DownloadStatus.CH <- n
	}
	return
}
