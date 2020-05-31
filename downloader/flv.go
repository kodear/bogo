package downloader

import (
	"fmt"
	"github.com/zhangpeihao/goflv"
	"os"
)

type FLVFileDownloader struct {
	FileDownloader
}

func (cls *FLVFileDownloader) Meta() *Meta {
	return &Meta{Name: "flv"}
}

func (cls *FLVFileDownloader) start() {
	defer 	close(cls.status.ch)

	// 获取视频大小
	if cls.status.MaxLength == 0 {
		for _, url := range cls.urls{
			res, err := cls.request(url)
			if err != nil{
				cls.status.Msg = err
				return
			}
			cls.status.MaxLength += cls.length(res)
		}
	}

	flvFile := cls.file + ".temporary"
	for _, url := range cls.urls{
		temporaryFile, err := os.Create(flvFile)
		if err != nil {
			cls.status.Msg = err
			return
		}

		res, err := cls.request(url)
		if err != nil {
			cls.status.Msg = err
			return
		}

		// 开始下载
		err = cls.download(res, temporaryFile)
		if err != nil{
			cls.status.Msg = err
			return
		}

		_ = temporaryFile.Close()
		err = cls.join(flvFile)
		if err != nil{
			cls.status.Msg = err
			return
		}
	}
}

func (cls *FLVFileDownloader)  Start(){
	go cls.start()
}

// 合并FLV碎片文件
func (cls *FLVFileDownloader) join(file string)(err error){
	if !pathExists(cls.file){
		_ = os.Rename(file, cls.file)
		return
	}

	flvFile, err := flv.OpenFile(cls.file)
	defer flvFile.Close()
	if err != nil{
		return
	}

	flvTemporaryFile, err := flv.OpenFile(file)
	if err != nil{
		return
	}

	var flvVideoTimestamp, flvAudioTimestamp uint32
	for {
		header, _, err := flvFile.ReadTag()
		if header.TagType == flv.VIDEO_TAG{
			flvVideoTimestamp = header.Timestamp
		}else if header.TagType == flv.AUDIO_TAG{
			flvAudioTimestamp = header.Timestamp
		}
		if err != nil{
			break
		}
	}

	for {
		header, data, err := flvTemporaryFile.ReadTag()
		if header.TagType == flv.VIDEO_TAG{
			err = flvFile.WriteVideoTag(data, header.Timestamp + flvVideoTimestamp)
		}else if header.TagType == flv.AUDIO_TAG{
			err = flvFile.WriteAudioTag(data, header.Timestamp + flvAudioTimestamp)
		}
		if err != nil{
			fmt.Println(err)
			break
		}
	}

	flvTemporaryFile.Close()
	//_ = os.Remove(file)

	return
}