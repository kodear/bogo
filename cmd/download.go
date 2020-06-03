package cmd

import (
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"github.com/zhxingy/bogo/downloader"
	"github.com/zhxingy/bogo/spider"
	"path/filepath"
	"sort"
	"time"
)

func Download(filename, path string, id, quality int, response *spider.Response) (err error) {
	streams := streams{}
	for _, stream := range (*response).Stream {
		streams = append(streams, stream)
	}
	sort.Sort(streams)

	// 获取下载ID
	var stream, defaultStream spider.Stream
	if len(response.Stream) < 2 {
		stream = response.Stream[0]
	}
	if id != 0 && stream.DownloadProtocol == "" {
		for _, x := range response.Stream {
			if x.ID == id {
				stream = x
				break
			}
		}
	}
	if quality != 0 && stream.DownloadProtocol == "" {
		for _, x := range response.Stream {
			if extractQuality(x.Quality) == quality {
				stream = x
				break
			} else if extractQuality(x.Quality) > quality {
				defaultStream = x
			}
		}
	}

	if stream.DownloadProtocol == "" {
		if defaultStream.DownloadProtocol == "" {
			stream = response.Stream[0]
		} else {
			stream = defaultStream
		}
	}

	ie, err := downloader.NewDownloader(stream.DownloadProtocol)
	if err != nil {
		return
	}

	if filename == "" {
		filename = formatTitle(response.Title, response.Part) + "." + stream.Format
	}

	ie.Initialize(filepath.Join(path, filename), stream.URLS, stream.DownloadHeaders)
	ie.Start()

	downloadStartTime := time.Now().Unix()
	out := "Site:  " + response.Site + "\n"
	out += "Title:  " + formatTitle(response.Title, response.Part) + "\n"
	if response.Duration != 0 {
		out += "Duration:  " + formatTimeString(response.Duration) + "\n"
	}
	out += "Streams:  "
	out = formatString(out, ":")
	out += formatString(sprintMediaStream(stream), ":") + "\n"
	for {
		if ie.Status().OK {
			fmt.Printf(out)
			break
		} else if ie.Status().Msg != nil {
			return ie.Status().Msg
		}
	}

	template := `{{string . "length"}} {{ bar . "<" "-" (cycle . "↖" "↗" "↘" "↙" ) "." ">"}} {{percent .}} {{string . "net_speed"}} {{string . "time"}}`
	bar := pb.ProgressBarTemplate(template).Start(ie.Status().MaxLength)
	downloadMaxBytes := formatSize(int64(ie.Status().MaxLength))
	hlsDownloadChunk := 1
	go func() {
		for {
			diffTime := time.Now().Unix() - downloadStartTime

			var speed, speed2 int
			if diffTime == 0 {
				speed = ie.Status().Byte
				speed2 = hlsDownloadChunk
			} else {
				speed = ie.Status().Byte / int(diffTime)
				speed2 = hlsDownloadChunk / int(diffTime)
			}
			bar.Set("net_speed", formatSize(int64(speed))+"/s")

			if stream.DownloadProtocol == "hls" || stream.DownloadProtocol == "hls_native" {
				bar.Set("time", formatTimeString((ie.Status().MaxLength-hlsDownloadChunk)/speed2))
				bar.Set("length", formatSize(int64(ie.Status().Byte)))
			} else {
				bar.Set("time", formatTimeString((ie.Status().MaxLength-ie.Status().Byte)/speed))
				bar.Set("length", formatSize(int64(ie.Status().Byte))+"/"+downloadMaxBytes)
			}
		}
	}()

	for {
		n, ok := <-ie.Status().CH
		if !ok {
			break
		}else{
			hlsDownloadChunk += 1
			bar.Add(n)
		}
	}
	bar.Finish()

	if ie.Status().Msg != nil {
		return ie.Status().Msg
	}

	return
}
