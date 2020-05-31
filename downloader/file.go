package downloader

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

type Status struct {
	Code    int
	Msg     error
	Byte    int
	MaxByte int
	ch      chan int
}

type Meta struct {
	Name string
}

type FileDownloader struct {
	urls       []string
	file       string
	ConfigFile string
	header     http.Header
	status     Status
}

func (cls *FileDownloader) Initialize(filename string, urls []string, size int, header http.Header) {
	cls.file = filename
	cls.ConfigFile = filename + ".json"
	cls.urls = urls
	ch := make(chan int, 1000)
	cls.status = Status{
		Code:    0,
		Msg:     nil,
		Byte:    0,
		MaxByte: size,
		ch:      ch,
	}
}

func (cls *FileDownloader) Meta() *Meta {
	panic("this method must be implemented by subclasses")
}

func (cls *FileDownloader) request(uri string, file *os.File) (err error) {
	client := &http.Client{}
	req, err := http.NewRequest("", uri, nil)
	if err != nil {
		return
	}

	req.Header = cls.header
	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer func() {
		_ = res.Body.Close()
	}()
	if res.StatusCode/100 > 3 {
		body, _ := ioutil.ReadAll(res.Body)
		return fmt.Errorf(string(body))
	}

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

		cls.status.ch <- n
	}

	return
}

func (cls *FileDownloader) start() {
	panic("this method must be implemented by subclasses")
}

func (cls *FileDownloader) Start() {
	go cls.start()
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

type config struct {
	filename string
}
