package downloader

import (
	"fmt"
	"net/http"
	"path/filepath"
	"testing"
)

func TestSyncDownloader(t *testing.T) {
	downloader, err := NewDownloader("http")
	if err != nil{
		t.Error(err)
	}

	downloader.Initialize(filepath.Join(getTestPath(), "test_sync_download_file.mp4"), []string{
		"http://vfx.mtime.cn/Video/2019/03/18/mp4/190318214226685784.mp4",
	}, http.Header{})
	downloader.Start()
	err = downloader.Wait()
	if err != nil{
		t.Error(err)
	}
}

func TestAsyncDownloader(t *testing.T) {
	downloader, err := NewDownloader("http")
	if err != nil{
		t.Error(err)
	}

	downloader.Initialize(filepath.Join(getTestPath(), "test_async_download_file.mp4"), []string{
		"http://vfx.mtime.cn/Video/2019/03/18/mp4/190318214226685784.mp4",
	}, http.Header{})
	downloader.Start()

	for {
		if _, ok := <-downloader.Status().ch; !ok {
			break
		}else{
			fmt.Printf("%d/%d\n", downloader.Status().Byte, downloader.Status().MaxLength)
		}
	}


	if downloader.Status().Msg != nil{
		t.Error(err)
	}
}