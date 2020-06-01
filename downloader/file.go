package downloader

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

type Status struct {
	Code      int
	Msg       error
	Length    int
	MaxLength int
	Byte      int
	ch        chan int
}

type Meta struct {
	Name string
}

type FileDownloader struct {
	urls   []string
	file   string
	header http.Header
	status Status
}

func (cls *FileDownloader) Initialize(filename string, urls []string, size int, header http.Header) {
	cls.file = filename
	cls.urls = urls
	ch := make(chan int, 1000)
	cls.header = header
	cls.status = Status{
		Code:      0,
		Msg:       nil,
		Length:    0,
		MaxLength: size,
		ch:        ch,
	}
}

func (cls *FileDownloader) Meta() *Meta {
	panic("this method must be implemented by subclasses")
}

func (cls *FileDownloader) request(uri string) (res *http.Response, err error) {
	client := &http.Client{}
	req, err := http.NewRequest("", uri, nil)
	if err != nil {
		return
	}

	req.Header = cls.header
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
		cls.status.Byte += n
		cls.status.ch <- n
	}
	return
}

func (cls *FileDownloader) start() {
	panic("this method must be implemented by subclasses")
}

func (cls *FileDownloader) Status() Status {
	return cls.status
}

func (cls *FileDownloader) Wait() (err error) {
	for {
		if _, ok := <-cls.status.ch; !ok {
			break
		}
	}

	if cls.status.Msg != nil {
		return cls.status.Msg
	}

	return
}
