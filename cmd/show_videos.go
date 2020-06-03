package cmd

import (
	"fmt"
	"github.com/zhxingy/bogo/spider"
	"strings"
)

func ShowVideoStreamInfo(response *spider.Response){
	var out string
	out = "Site:  " + response.Site +  "\n"
	out += "Title:  " + formatTitle(response.Title, response.Part) +  "\n"
	if response.Duration != 0{
		out += "Duration:  " + formatTimeString(response.Duration)  + "\n"
	}
	out += "Streams:  # All available quality"
	out = formatString(out, ":")
	for _, stream := range response.Stream{
		out += formatString(streamInfo(stream), ":") + "\n"
	}
	fmt.Printf(out)

}

func videoStreamByIDInfo(response *spider.Response, id int) string{
	var out string
	out = "Site:  " + response.Site +  "\n"
	out += "Title:  " + formatTitle(response.Title, response.Part) +  "\n"
	if response.Duration != 0{
		out += "Duration:  " + formatTimeString(response.Duration)  + "\n"
	}
	out += "Streams:  "
	out = formatString(out, ":")
	for _, stream := range response.Stream{
		if stream.ID == id{
			out += formatString(streamInfo(stream), ":") + "\n"
		}
	}

	return out
}

func streamInfo(stream spider.Stream)(info string){
	var out string
	out = fmt.Sprintf("    [%d]  -------------------\n", stream.ID)
	out += "    Quality:  " + stream.Quality + "\n"
	if stream.StreamType != ""{
		out += "    Type:  " + stream.StreamType + "\n"
	}
	if stream.Size != 0{
		out += fmt.Sprintf("    Size:  %s (%d Bytes)\n", formatSize(int64(stream.Size)), stream.Size)
	}
	if stream.Height != 0 && stream.Width != 0{
		out += fmt.Sprintf("    Window:  %dx%d\n", stream.Width, stream.Height)
	}
	out += fmt.Sprintf("    # download with: -f %d ...\n", stream.ID)
	return out
}

func formatString(str, seq string) string{
	var newString string
	var n int
	for _, s := range strings.Split(str, "\n"){
		if s == ""|| len(strings.Split(s, seq)) == 1{
			continue
		}
		m := len(strings.Split(s, seq)[0])
		if m > n{
			n = m
		}
	}


	for _, s := range strings.Split(str, "\n"){
		if s == ""{
			continue
		}
		sliceString := strings.Split(s, seq)
		if len(sliceString) == 1{
			newString += s + "\n"
			continue
		}
		key := sliceString[0]
		value  := strings.Join(sliceString[1:], seq)
		x := n
		m := len(key)
		for m < x{
			value = " " + value
			x -= 1
		}
		newString += key + seq + value + "\n"
	}

	return newString
}

func formatTitle(title, part string) string{
	if title == part || part == ""{
		return title
	}else{
		return  title + "-" + part
	}
}

func formatTimeString(timestamp int) string{
	minutes := timestamp / 60
	seconds := timestamp % 60
	hours := minutes / 60
	minutes = minutes % 60

	if hours >= 24 {
		return  ">1d"
	}

	return fmt.Sprintf("%dh%dm%ds", hours, minutes, seconds)
}

func formatSize(fileSize int64) (size string) {
	if fileSize < 1024 {
		return fmt.Sprintf("%.2f B", float64(fileSize)/float64(1))
	} else if fileSize < (1024 * 1024) {
		return fmt.Sprintf("%.2f KB", float64(fileSize)/float64(1024))
	} else if fileSize < (1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2f MB", float64(fileSize)/float64(1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2f GB", float64(fileSize)/float64(1024*1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2f TB", float64(fileSize)/float64(1024*1024*1024*1024))
	} else { //if fileSize < (1024 * 1024 * 1024 * 1024 * 1024 * 1024)
		return fmt.Sprintf("%.2f EB", float64(fileSize)/float64(1024*1024*1024*1024*1024))
	}
}