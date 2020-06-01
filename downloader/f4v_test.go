package downloader

import (
	flv "github.com/zhangpeihao/goflv"
	"os"
	"path/filepath"
	"testing"
)

func TestF4VFileDownloader_Start(t *testing.T) {

}

func TestF4VFileDownloader_join(t *testing.T) {
	home := os.Getenv("BOGO")

	testFilePath := filepath.Join(home[1:len(home)-1], "test_data", "video")
	testJoinFile := filepath.Join(testFilePath, "test_join.f4v")
	testFile := filepath.Join(testFilePath, "test.f4v")

	test := F4VFileDownloader{}
	flvFile, err := flv.CreateFile(testJoinFile)
	if err != nil {
		return
	}

	var a, v uint32
	var testFiles = []string{"test_1.f4v", "test_2.f4v", "test_3.f4v"}
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
