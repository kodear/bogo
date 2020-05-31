package downloader

import (
	"net/http"
	"testing"
)

func TestFileDownloader_Start(t *testing.T) {
	test := HTTPFileDownloader{}
	test.Initialize("./test.mp4", []string{"http://v3-tt.ixigua.com/27c2e899a7202c4921d766e324638093/5ed40b3f/video/tos/cn/tos-cn-ve-67/b28d5f00877b484895609c9693458800/?a=1768&br=2124&bt=708&cr=0&cs=0&dr=0&ds=1&er=0&l=202006010114480100290550381BF1F37F&lr=default&mime_type=video%2Fmp4&qs=0&rc=anJpZThnZzZzcDMzPDYzM0ApZjo3M2U7PGU1N2c6ZmloPGdvbzZjby4xXzFfLS0tMS9zc2BeL2I2X19iYy1eL2NgL2E6Yw%3D%3D&vl=&vr="}, 0, http.Header{})
	test.Start()
	err := test.Wait()
	if err != nil {
		t.Error(err)
	}
}
