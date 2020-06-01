package downloader

import (
	flv "github.com/zhangpeihao/goflv"
	"github.com/zhxingy/bogo/spider"
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

func TestFLVFileDownloader_Start(t *testing.T) {
	home := os.Getenv("BOGO")
	testDownloadPath := filepath.Join(home[1:len(home)-1], "test_data", "download")

	body, _ := spider.Do("https://www.bilibili.com/bangumi/play/ep15014", nil)
	for _, p := range body {
		if p.DownloadProtocol == "flv" {
			headers := http.Header{}
			var urls []string
			for _, url := range p.Links {
				urls = append(urls, url.URL)
			}
			for k, v := range p.DownloadHeaders {
				headers[k] = []string{v}
			}

			test := FLVFileDownloader{}
			test.Initialize(filepath.Join(testDownloadPath, "test.flv"), urls, 0, headers)
			test.Start()
			err := test.Wait()
			if err != nil {
				t.Error(err)
			}
		}
	}
}

func TestFLVFileDownloader_join(t *testing.T) {
	home := os.Getenv("BOGO")

	testFilePath := filepath.Join(home[1:len(home)-1], "test_data", "video")
	testJoinFile := filepath.Join(testFilePath, "test_join.flv")
	testFile := filepath.Join(testFilePath, "test.flv")

	test := FLVFileDownloader{}
	flvFile, err := flv.CreateFile(testJoinFile)
	if err != nil {
		return
	}

	var a, v uint32
	var testFiles = []string{"test_1.flv", "test_2.flv", "test_3.flv"}
	for _, testFileName := range testFiles {
		testTempFile := filepath.Join(testFilePath, testFileName)
		_, err := copyFile(testTempFile, testFile)
		if err != nil {
			t.Error(err)
		}

		err = test.join(testTempFile, flvFile, &a, &v)
		if err != nil {
			t.Error(err)
		}

	}

	flvFile.Close()
	_ = os.Remove(testJoinFile)
}
