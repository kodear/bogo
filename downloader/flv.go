package downloader

import (
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
	defer close(cls.status.ch)

	// 获取视频大小
	for _, url := range cls.urls {
		res, err := cls.request(url)
		if err != nil {
			cls.status.Msg = err
			return
		}
		cls.status.MaxLength += cls.length(res)
	}

	flvFile, err := flv.CreateFile(cls.file)
	if err != nil {
		cls.status.Msg = err
		return
	}
	defer flvFile.Close()
	var flvVideoTimestamp, flvAudioTimestamp uint32

	temporaryFile := cls.file + ".temporary"
	for _, url := range cls.urls {
		tf, err := os.Create(temporaryFile)
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
		err = cls.download(res, tf)
		if err != nil {
			cls.status.Msg = err
			return
		}

		_ = tf.Close()
		err = cls.join(temporaryFile, flvFile, &flvVideoTimestamp, &flvAudioTimestamp)
		if err != nil {
			cls.status.Msg = err
			return
		}
	}
}

func (cls *FLVFileDownloader) Start() {
	go cls.start()
}

func (cls *FLVFileDownloader) join(file string, flvFile *flv.File, flvVideoTimestamp, flvAudioTimestamp *uint32) (err error) {
	var flvTempVideoTimestamp, flvTempAudioTimestamp uint32
	flvTemporaryFile, err := flv.OpenFile(file)
	if err != nil {
		return
	}

	for {
		header, data, err := flvTemporaryFile.ReadTag()
		if err != nil {
			*flvVideoTimestamp += flvTempVideoTimestamp
			*flvAudioTimestamp += flvTempAudioTimestamp
			break
		}
		if header.TagType == flv.VIDEO_TAG {
			flvTempVideoTimestamp = header.Timestamp
			err = flvFile.WriteVideoTag(data, header.Timestamp+*flvVideoTimestamp)
		} else if header.TagType == flv.AUDIO_TAG {
			flvTempAudioTimestamp = header.Timestamp
			err = flvFile.WriteAudioTag(data, header.Timestamp+*flvAudioTimestamp)
		}
	}

	flvTemporaryFile.Close()
	_ = os.Remove(file)

	return
}
