package downloader

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

func TestHLSFileDownloader_Start(t *testing.T) {
	home := os.Getenv("BOGO")
	testDownloadFile := filepath.Join(home[1:len(home)-1], "test_data", "download", "testHlsDownload.ts")
	test := HLSFileDownloader{}
	test.Initialize(testDownloadFile, []string{"https://meiju11.qfxmj.com/20200514/F1kJTMiM/2000kb/hls/index.m3u8?wsSecret=755f8f3be85bc620ae1c510f87c3fe20&wsTime=1591013964&watch=ae6f4caa08e521511949081bd4dd75f9"}, 0, http.Header{})
	test.Start()
	err := test.Wait()
	if err != nil {
		t.Error(err)
	}
}
