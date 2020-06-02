package downloader

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

type ISMFileDownloader struct {
	FileDownloader
}

func (cls *ISMFileDownloader) Meta() *Meta {
	return &Meta{Name: "ism"}
}

func (cls *ISMFileDownloader) start() {
	defer close(cls.DownloadStatus.ch)

	// 获取视频大小
	for _, url := range cls.URLS {
		res, err := cls.request(url)
		if err != nil {
			cls.DownloadStatus.Msg = err
			return
		}
		cls.DownloadStatus.MaxLength += cls.length(res)
	}

	temporaryFile := cls.File + ".temporary"
	for _, url := range cls.URLS {
		tf, err := os.Create(temporaryFile)
		if err != nil {
			cls.DownloadStatus.Msg = err
			return
		}

		res, err := cls.request(url)
		if err != nil {
			cls.DownloadStatus.Msg = err
			return
		}

		// 开始下载
		err = cls.download(res, tf)
		if err != nil {
			cls.DownloadStatus.Msg = err
			return
		}

		_ = tf.Close()
		err = cls.join(temporaryFile)
		if err != nil {
			cls.DownloadStatus.Msg = err
			return
		}
	}
}

func (cls *ISMFileDownloader) Start() {
	go cls.start()
}

func (cls *ISMFileDownloader) join(file string) (err error) {
	path, err := exec.LookPath("ffmpeg")
	if err != nil{
		return err
	}
	if !pathExists(cls.File){
		err = os.Rename(file, cls.File)
		return
	}

	outFile := cls.File + ".mp4"
	mergeFile := cls.File + ".txt"
	f, _ := os.Create(mergeFile)
	_, _ = f.Write([]byte(fmt.Sprintf("file '%s'\nfile '%s'\n", cls.File, file)))
	cmd := exec.Command(
		path, "-y", "-f", "concat", "-safe", "-1",
		"-i", mergeFile, "-c", "copy", "-bsf:a", "aac_adtstoasc", outFile,
	)

	var stderr bytes.Buffer
	err = cmd.Run()
	cmd.Stderr = &stderr
	if err != nil{
		return fmt.Errorf("%s %s", err, stderr.String())
	}

	_ = f.Close()
	_ = os.Remove(mergeFile)
	_ = os.Remove(file)
	_ = os.Remove(cls.File)
	err = os.Rename(outFile, cls.File)
	if err != nil{
		return
	}

	return
}