package download

import (
	"github.com/zhangpeihao/goflv"
	"io"
	"net/http"
	"os"
)

type HTTP struct {
	Download
}

func (h *HTTP) Do(link, text, file, fname string, links []string) {
	f, err := os.Create(fname)
	if err != nil {
		h.DownloadMessage = err
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		h.DownloadMessage = err
		return
	}

	for k, v := range h.Headers {
		req.Header.Add(k, v)
	}

	res, err := client.Do(req)
	if err != nil {
		h.DownloadMessage = err
		return
	}

	if h.Len == 0 {
		h.Len = int(res.ContentLength)
	}
	h.ProgressINIT = true

	h.Ch = make(chan int, 1000)
	buf := make([]byte, 4096)
	for {
		n, err := res.Body.Read(buf)
		if n > 0 {
			_, _ = f.Write(buf[:n])
		} else if err != nil && err == io.EOF {
			_ = f.Close()
			break
		} else if err != nil {
			h.DownloadMessage = err
			return
		}

		if !h.DownloadINIT {
			h.DownloadINIT = true
		}

		h.Ch <- n
	}
	_ = res.Body.Close()

	h.DownloadStatus = true

}

type HTTPSegFLV struct {
	Download
}

func (h *HTTPSegFLV) Do(link, text, file, fname string, links []string) {

	h.ProgressINIT = true
	flvFile, err := flv.CreateFile(fname)
	h.Ch = make(chan int, 1000)

	if err != nil {
		h.DownloadMessage = err
		return
	}

	var audioTime uint32
	var videoTime uint32

	for _, l := range links {
		// ==========================================================================================================

		f, err := os.Create(fname + ".tmp")
		if err != nil {
			h.DownloadMessage = err
			return
		}

		client := &http.Client{}
		req, err := http.NewRequest("GET", l, nil)
		if err != nil {
			h.DownloadMessage = err
			return
		}

		for k, v := range h.Headers {
			req.Header.Add(k, v)
		}

		res, err := client.Do(req)
		if err != nil {
			h.DownloadMessage = err
			return
		}

		buf := make([]byte, 4096)
		for {
			n, err := res.Body.Read(buf)
			if n > 0 {
				_, _ = f.Write(buf[:n])
			}
			if err != nil && err == io.EOF {
				_ = f.Close()
				break
			}
			if err != nil {
				h.DownloadMessage = err
				return
			}

			if !h.DownloadINIT {
				h.DownloadINIT = true
			}
			h.Ch <- n
		}
		_ = res.Body.Close()

		// ==========================================================================================================
		fi, err := flv.OpenFile(fname + ".tmp")
		if err != nil {
			h.DownloadMessage = err
			return
		}

		var nowAudioTime uint32
		var nowVideoTime uint32
		for {
			header, data, err := fi.ReadTag()
			if err != nil {
				audioTime += nowAudioTime
				videoTime += nowVideoTime
				break
			}

			if header.TagType == 8 {
				nowAudioTime = header.Timestamp
				err = flvFile.WriteAudioTag(data, header.Timestamp+audioTime)

			} else if header.TagType == 9 {
				nowVideoTime = header.Timestamp
				err = flvFile.WriteVideoTag(data, header.Timestamp+videoTime)
			}

			if err != nil {
				break
			}
		}
		fi.Close()
		err = os.Remove(fname + ".tmp")
	}

	h.DownloadStatus = true

}

type HTTPSegF4V struct {
	Download
}

func (h *HTTPSegF4V) Do(link, text, file, fname string, links []string) {

	h.ProgressINIT = true
	flvFile, err := flv.CreateFile(fname)
	h.Ch = make(chan int, 1000)

	if err != nil {
		h.DownloadMessage = err
		return
	}

	var audioTime uint32
	var videoTime uint32

	for _, l := range links {
		// ==========================================================================================================

		f, err := os.Create(fname + ".tmp")
		if err != nil {
			h.DownloadMessage = err
			return
		}

		client := &http.Client{}
		req, err := http.NewRequest("GET", l, nil)
		if err != nil {
			h.DownloadMessage = err
			return
		}

		for k, v := range h.Headers {
			req.Header.Add(k, v)
		}

		res, err := client.Do(req)
		if err != nil {
			h.DownloadMessage = err
			return
		}

		buf := make([]byte, 4096)
		for {
			n, err := res.Body.Read(buf)
			if n > 0 {
				_, _ = f.Write(buf[:n])
			}
			if err != nil && err == io.EOF {
				_ = f.Close()
				break
			}
			if err != nil {
				h.DownloadMessage = err
				return
			}

			if !h.DownloadINIT {
				h.DownloadINIT = true
			}
			h.Ch <- n
		}
		_ = res.Body.Close()

		// ==========================================================================================================
		fi, err := flv.OpenFile(fname + ".tmp")
		if err != nil {
			h.DownloadMessage = err
			return
		}

		//var nowAudioTime uint32
		//var nowVideoTime uint32
		for {
			header, data, err := fi.ReadTag()
			if err != nil {
				//audioTime += nowAudioTime
				//videoTime += nowVideoTime
				break
			}

			if header.TagType == 8 {
				//nowAudioTime = header.Timestamp
				err = flvFile.WriteAudioTag(data, header.Timestamp+audioTime)

			} else if header.TagType == 9 {
				//nowVideoTime = header.Timestamp
				err = flvFile.WriteVideoTag(data, header.Timestamp+videoTime)
			}

			if err != nil {
				break
			}
		}
		fi.Close()
		err = os.Remove(fname + ".tmp")
	}

	h.DownloadStatus = true

}
