package downloader

type F4VFileDownloader struct {
	FLVFileDownloader
}

func (cls *F4VFileDownloader) Meta() *Meta {
	return &Meta{Name: "f4v"}
}

// 合并F4V碎片文件
//func (cls *F4VFileDownloader) join(file string, flvFile *flv.File, flvVideoTimestamp, flvAudioTimestamp *uint32) (err error) {
//	flvTemporaryFile, err := flv.OpenFile(file)
//	if err != nil {
//		return
//	}
//
//	for {
//		header, data, err := flvTemporaryFile.ReadTag()
//		if err != nil {
//			break
//		}
//		if header.TagType == flv.VIDEO_TAG {
//			err = flvFile.WriteVideoTag(data, header.Timestamp)
//		} else if header.TagType == flv.AUDIO_TAG {
//			err = flvFile.WriteAudioTag(data, header.Timestamp)
//		}
//	}
//
//	flvTemporaryFile.Close()
//	_ = os.Remove(file)
//
//	return
//}
