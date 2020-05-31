package downloader

import (
	"testing"
)

//func TestFLVFileDownloader_Start(t *testing.T) {
//	x, _ := spider.Do("https://www.bilibili.com/video/BV1jJ411c7s3?p=20", nil)
//	for _, p := range x{
//		if p.DownloadProtocol == "flv"{
//			var urls []string
//			for _, url := range p.Links{
//				urls = append(urls, url.URL)
//			}
//			headers := http.Header{}
//			for k, v := range p.DownloadHeaders{
//				headers[k] = []string{v}
//			}
//			test := FLVFileDownloader{}
//			test.Initialize("./test.flv", urls, 0, headers)
//			test.Start()
//			err := test.Wait()
//			if err != nil{
//				t.Error(err)
//			}
//		}
//	}
//}

func TestFLVFileDownloader_join(t *testing.T) {
	test := FLVFileDownloader{}
	test.Initialize(`H:\code\go\src\github.com\zhxingy\bogo\test_data\video\test.flv`, nil, 0, nil)
	err := test.join(`H:\code\go\src\github.com\zhxingy\bogo\test_data\video\test_1.flv`)
	if err != nil{
		t.Error(err)
	}
}