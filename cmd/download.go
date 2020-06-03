package cmd

import (
	"fmt"
	//"github.com/cheggaaa/pb/v3"
	//"github.com/zhxingy/bogo/downloader"
	"github.com/zhxingy/bogo/spider"
	"net/http"
	//"path/filepath"
	"regexp"
	"sort"
	"strconv"
	//"time"
)

func extractQuality(quality string) int{
	regex, _ := regexp.Compile(`(\d+)`)
	strings := regex.FindAllStringSubmatch(quality, -1)
	if len(strings) == 0 || len(strings[0]) == 1{
		return 0
	}

	result, _ := strconv.Atoi(strings[0][1])
	return result
}

func Download(url, filename, path string, jar spider.CookiesJar, header http.Header, id, quality int)(err error)  {
	response, err := extract(url, header, jar)
	if err != nil{
		return
	}

	streams := Streams{}
	for _, stream := range(*response).Stream{
		streams = append(streams, stream)
	}
	sort.Sort(streams)

	var stream spider.Stream
	if len(response.Stream) < 2 {
		stream = response.Stream[0]
	}else if id != 0 {
		for _, x := range response.Stream{
			if x.ID == id{
				stream = x
				break
			}
		}
	}else if quality != 0{
		tempDict := map[int]spider.Stream{}
		for _, x := range streams{
			fmt.Println(x)
			if extractQuality(x.Quality) == quality{
				stream = x
				break
			}else{
				tempDict[extractQuality(x.Quality)] = x
			}
		}
		//if len(stream.URLS) == 0 {
		//	var y spider.Stream

		//}
	}else{
		stream = response.Stream[0]
	}
	fmt.Println(stream)

	//ie, err := downloader.NewDownloader(stream.DownloadProtocol)
	//if err != nil{
	//	return
	//}
	//
	//if filename == ""{
	//	filename = formatTitle(response.Title, response.Part) + "." + stream.Format
	//}
	//
	//ie.Initialize(filepath.Join(path, filename), stream.URLS, stream.DownloadHeaders)
	//ie.Start()
	//
	//for {
	//	if ie.Status().MaxLength > 0 {
	//		break
	//	}else if ie.Status().Msg != nil{
	//		return ie.Status().Msg
	//	}
	//}
	//
	//template := "{{string . \"length\"}} {{ bar . \"\" \"=\" (cycle . \">\" ) \"-\" \"\"}} {{percent .}} {{string . \"net_speed\"}} {{string . \"time\"}}"
	//bar := pb.ProgressBarTemplate(template).Start(ie.Status().MaxLength)
	//size := formatSize(int64(ie.Status().MaxLength))
	//startTime := time.Now().Unix()
	//chunkNum := 0
	//fmt.Printf(videoStreamByIDInfo(response, id))
	//for {
	//	if n, ok := <-ie.Status().CH; !ok {
	//		break
	//	}else {
	//		bar.Add(n)
	//
	//		chunkNum += 1
	//		diffTime := time.Now().Unix() - startTime
	//
	//		var speed, speed2 int
	//		if diffTime == 0{
	//			speed = ie.Status().Byte
	//			speed2 = chunkNum
	//		}else{
	//			speed = ie.Status().Byte / int(diffTime)
	//			speed2 = chunkNum / int(diffTime)
	//		}
	//		bar.Set("net_speed", formatSize(int64(speed)) + "/s")
	//
	//		if stream.DownloadProtocol == "hls" || stream.DownloadProtocol == "hls_native"  {
	//			bar.Set("time", formatTimeString((ie.Status().MaxLength - chunkNum)/speed2))
	//			bar.Set("length", formatSize(int64(ie.Status().Byte)))
	//		}else{
	//			bar.Set("time", formatTimeString((ie.Status().MaxLength - ie.Status().Byte)/speed))
	//			bar.Set("length", formatSize(int64(ie.Status().Byte)) + "/" + size)
	//		}
	//	}
	//}
	//bar.Finish()
	//
	//if ie.Status().Msg != nil{
	//	return ie.Status().Msg
	//}

	return
}

type Streams []spider.Stream

func (stream Streams)Len() int{
	return len(stream)
}

func (stream Streams)Less(i, j int) bool{
	return extractQuality(stream[i].Quality) < extractQuality(stream[j].Quality)
}

func  (stream Streams)Swap(i, j int){
	stream[i] = stream[j]
}