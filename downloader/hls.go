package downloader

import (
	"github.com/grafov/m3u8"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

type HLSFileDownloader struct {
	FileDownloader
}

func (cls *HLSFileDownloader) Meta() *Meta {
	return &Meta{Name: "hls"}
}

func (cls *HLSFileDownloader) start() {
	reader, err := cls.open()
	if err != nil {
		cls.DownloadStatus.Msg = err
	} else {
		cls.run(reader)
	}
}

func (cls *HLSFileDownloader) Start() {
	go cls.start()
}

func (cls *HLSFileDownloader) download(res *http.Response, file *os.File, key []byte) (err error) {
	defer func() { _ = res.Body.Close() }()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	var decryptBody []byte
	if key == nil {
		decryptBody = body
	} else {
		decryptBody = AESDecrypt(body, key)
	}

	n, err := file.Write(decryptBody)
	if err != nil {
		return
	}
	cls.DownloadStatus.Byte += n
	cls.DownloadStatus.CH <- 1

	return
}

func (cls *HLSFileDownloader) parse(reader io.Reader) (playlist *m3u8.MediaPlaylist, err error) {
	hls, listType, err := m3u8.DecodeFrom(reader, true)
	if err != nil {
		return
	}

	switch listType {
	case m3u8.MASTER:
		masterPlaylist := hls.(*m3u8.MasterPlaylist)
		variants := masterPlaylist.Variants[len(masterPlaylist.Variants)-1]
		masterURL := urlJoin(cls.URL, variants.URI)
		res, err := cls.request(masterURL)
		if err != nil {
			return nil, err
		}

		return cls.parse(res.Body)

	case m3u8.MEDIA:
		playlist = hls.(*m3u8.MediaPlaylist)
	}

	return
}

func (cls *HLSFileDownloader) open() (io.Reader, error) {
	res, err := cls.request(cls.URL)
	if err != nil {
		return nil, err
	}

	return res.Body, nil
}

func (cls *HLSFileDownloader) run(reader io.Reader) {
	defer close(cls.DownloadStatus.CH)

	file, err := os.Create(cls.File)
	if err != nil {
		cls.DownloadStatus.Msg = err
		return
	}
	defer func() { _ = file.Close() }()

	playlist, err := cls.parse(reader)
	if err != nil {
		cls.DownloadStatus.Msg = err
		return
	}

	for _, segment := range playlist.Segments {
		if segment == nil {
			break
		}
		cls.DownloadStatus.MaxLength += 1
	}
	cls.DownloadStatus.OK = true

	var key []byte
	if playlist.Key != nil {
		res, err := cls.request(playlist.Key.URI)
		if err != nil {
			cls.DownloadStatus.Msg = err
			return
		}
		key, err = ioutil.ReadAll(res.Body)
		if err != nil {
			cls.DownloadStatus.Msg = err
			return
		}
		_ = res.Body.Close()
	}

	for _, segment := range playlist.Segments {
		if segment == nil {
			break
		}

		res, err := cls.request(urlJoin(cls.URL, segment.URI))
		if err != nil {
			cls.DownloadStatus.Msg = err
			return
		}

		err = cls.download(res, file, key)
		if err != nil {
			cls.DownloadStatus.Msg = err
			return
		}
	}

	return

}
