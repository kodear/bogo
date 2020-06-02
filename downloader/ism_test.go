package downloader

import (
	"net/http"
	"path/filepath"
	"testing"
)

func TestISMFileDownloader_Start(t *testing.T) {
	test := ISMFileDownloader{}
	test.Initialize(filepath.Join(getTestPath(), "test_download_ism.mp4"), []string{
		"http://vfx.mtime.cn/Video/2019/03/18/mp4/190318214226685784.mp4",
		"http://vfx.mtime.cn/Video/2019/03/12/mp4/190312083533415853.mp4",
	}, http.Header{})
	test.Start()
	err := test.Wait()
	if err != nil {
		t.Error(err)
	}
}