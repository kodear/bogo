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
	defer close(cls.status.ch)

	file, err := os.Create(cls.file)
	if err != nil {
		cls.status.Msg = err
		return
	}
	defer func() { _ = file.Close() }()

	res, err := cls.request(cls.urls[0])
	if err != nil {
		cls.status.Msg = err
		return
	}

	playlist, err := cls.parse(res.Body)
	if err != nil {
		return
	}

	cls.status.MaxLength = len(playlist.Segments)

	var key []byte
	if playlist.Key != nil {
		res, err := cls.request(playlist.Key.URI)
		if err != nil {
			return
		}
		key, err = ioutil.ReadAll(res.Body)
		if err != nil {
			return
		}
		_ = res.Body.Close()
	}

	for _, segment := range playlist.Segments {
		if segment == nil {
			break
		}

		res, err := cls.request(urlJoin(cls.urls[0], segment.URI))
		if err != nil {
			return
		}

		err = cls.download(res, file, key)
		if err != nil {
			return
		}
	}

	return

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

	cls.status.Byte += n
	cls.status.ch <- 1

	return
}

func (cls *HLSFileDownloader) Start() {
	go cls.start()
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
		masterURL := urlJoin(cls.urls[0], variants.URI)
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
