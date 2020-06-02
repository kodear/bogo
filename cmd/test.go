package main

import (
	"github.com/zhxingy/bogo/downloader"
	"net/http"
)

func main() {
	// get the downloader object
	ie, err := downloader.NewDownloader("http")
	if err != nil {
		panic(err)
	}

	// initialize downloader
	ie.Initialize("test_sync_download_file.mp4", []string{"http://vfx.mtime.cn/Video/2019/03/18/mp4/190318214226685784.mp4"}, http.Header{})

	// start downloader
	ie.Start()

	// wait for download to complete
	err = ie.Wait()

	// if the msg is not nil, the download failed.
	if err != nil {
		panic(err)
	}
}
