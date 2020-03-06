package download

import (
	flv "github.com/zhangpeihao/goflv"
	"io"
	"net/http"
	"os"
)

type HTTP struct {
	Download
}

func (h *HTTP) Do(link, text, file, fname string, links []string) {
	h.Init = true

	f, err := os.Create(fname)
	if err != nil {
		h.Err = err
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		h.Err = err
		return
	}

	for k, v := range h.Headers {
		req.Header.Add(k, v)
	}

	res, err := client.Do(req)
	if err != nil {
		h.Err = err
		return
	}

	h.Ch = make(chan int, 1000)
	buf := make([]byte, 4096)
	for {
		n, err := res.Body.Read(buf)
		if err != nil {
			f.Close()
			break
		}
		f.Write(buf[:n])
		if !h.C {
			h.C = true
		}
		h.Ch <- n
	}
	res.Body.Close()

	h.Ok = true

}

type HTTPSegFLV struct {
	Download
}

func (h *HTTPSegFLV) Do(link, text, file, fname string, links []string) {

	h.Init = true
	flvFile, err := flv.CreateFile(fname)
	h.Ch = make(chan int, 1000)

	if err != nil {
		h.Err = err
		return
	}

	var audioTime uint32
	var videoTime uint32

	for _, l := range links {
		// ==========================================================================================================

		f, err := os.Create(fname + ".tmp")
		if err != nil {
			h.Err = err
			return
		}

		client := &http.Client{}
		req, err := http.NewRequest("GET", l, nil)
		if err != nil {
			h.Err = err
			return
		}

		for k, v := range h.Headers {
			req.Header.Add(k, v)
		}

		res, err := client.Do(req)
		if err != nil {
			h.Err = err
			return
		}

		buf := make([]byte, 4096)
		for {
			n, err := res.Body.Read(buf)
			if err != nil && err == io.EOF {
				f.Close()
				break
			}
			if err != nil {
				h.Err = err
				return
			}

			f.Write(buf[:n])
			if !h.C {
				h.C = true
			}
			h.Ch <- n
		}
		res.Body.Close()

		// ==========================================================================================================
		fi, err := flv.OpenFile(fname + ".tmp")
		if err != nil {
			h.Err = err
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

	h.Ok = true

}
