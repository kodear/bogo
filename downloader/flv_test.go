package downloader

import (
	"net/http"
	"path/filepath"
	"testing"
)

func TestFLVFileDownloader_Start(t *testing.T) {
	test := FLVFileDownloader{}
	test.Initialize(filepath.Join(getTestPath(), "test_download_flv.flv"), []string{
		"http://samples.mplayerhq.hu/FLV/asian-commercials-are-weird.flv",
		"http://samples.mplayerhq.hu/FLV/asian-commercials-are-weird.flv",
	}, http.Header{})
	test.Start()
	err := test.Wait()
	if err != nil {
		t.Error(err)
	}
}
