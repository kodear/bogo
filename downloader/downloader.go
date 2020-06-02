package downloader

import (
	"errors"
	"net/http"
)

type Downloader interface {
	Start()
	Status() DownloadStatus
	Wait() error
	Meta() *Meta
	Initialize(string, []string, http.Header)
}

var downloader = []Downloader{
	&HTTPFileDownloader{},
	&FLVFileDownloader{},
	&F4VFileDownloader{},
	&ISMFileDownloader{},
	&HLSFileDownloader{},
	&HLSNativeFileDownloader{},
}

func NewDownloader(protocol string) (cls Downloader, err error){
	for _, ie := range downloader{
		if ie.Meta().Name == protocol{
			cls = ie
			break
		}
	}

	if cls == nil{
		err = errors.New("did not match to the downloader: "  + protocol)
	}

	return cls, err
}

