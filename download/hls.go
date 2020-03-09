package download

import (
	"bufio"
	"errors"
	"github.com/grafov/m3u8"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type HLS struct {
	Download
}

func (h *HLS) Parse(link, text, file string) (urls []string, err error) {
	var reader io.Reader

	if link != "" {
		client := http.Client{}
		res, err := client.Get(link)
		if err != nil {
			return nil, err
		}
		reader = res.Body

	} else if text != "" {
		reader = strings.NewReader(text)
	} else if file != "" {
		f, err := os.Open(file)
		if err != nil {
			return nil, err
		}
		reader = bufio.NewReader(f)
	} else {
		err = errors.New("input not found")
		return nil, err
	}

	p, listType, err := m3u8.DecodeFrom(bufio.NewReader(reader), true)
	if err != nil {
		return nil, err
	}

	switch listType {
	case m3u8.MASTER:
		x := p.(*m3u8.MasterPlaylist)
		newLink := x.Variants[0].URI
		if strings.HasPrefix(newLink, "http") {
			return h.Parse(newLink, "", "")
		}
		if link == "" {
			err = errors.New("a secondary address appears, but the source address is not found")
			return
		}

		if strings.HasPrefix(newLink, "/") {
			u, err := url.Parse(link)
			if err != nil {
				return nil, err
			}
			newLink = u.Scheme + "://" + u.Host + "/" + newLink
		} else {
			s := strings.Split(link, "/")
			newLink = strings.Join(s[0:len(s)-1], "/") + "/" + newLink
		}

		return h.Parse(newLink, "", "")

	case m3u8.MEDIA:
		x := p.(*m3u8.MediaPlaylist)
		for _, body := range x.Segments {
			if body != nil {
				if strings.HasPrefix(body.URI, "http") {
					urls = append(urls, body.URI)
				} else if strings.HasPrefix(body.URI, "/") && link != "" {
					u, err := url.Parse(link)
					if err != nil {
						return nil, err
					}
					urls = append(urls, u.Scheme+"://"+u.Host+"/"+body.URI)
				} else if link != "" {
					s := strings.Split(link, "/")
					urls = append(urls, strings.Join(s[0:len(s)-1], "/")+"/"+body.URI)
				} else {
					err = errors.New("playlists have no parent links")
				}
			}

		}
	}

	return
}

func (h *HLS) Do(link, text, file, fname string, links []string) {
	urls, err := h.Parse(link, text, file)
	if err != nil {
		h.DownloadMessage = err
		return
	}

	out, err := os.Create(fname)
	if err != nil {
		h.DownloadMessage = err
		return
	}

	defer out.Close()

	var ok bool
	if h.Len == 0 {
		h.Len = len(urls)
		ok = true
	}

	h.ProgressINIT = true

	h.Ch = make(chan int, 1000)
	for _, link := range urls {
		client := &http.Client{}
		res, err := client.Get(link)
		if err != nil {
			h.DownloadMessage = err
			return
		}

		_, err = io.Copy(out, res.Body)
		if err != nil {
			h.DownloadMessage = err
			return
		}

		if !h.DownloadINIT {
			h.DownloadINIT = true
		}

		if ok {
			h.Ch <- 1
		} else {
			h.Ch <- int(res.ContentLength)
		}

		_ = res.Body.Close()
	}

	h.DownloadStatus = true
}
