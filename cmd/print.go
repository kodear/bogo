package cmd

import (
	"fmt"
	"github.com/zhxingy/bogo/spider"
)

func PrintMedia(response *spider.Response) {
	var out string
	out = "Site:  " + response.Site + "\n"
	out += "Title:  " + formatTitle(response.Title, response.Part) + "\n"
	if response.Duration != 0 {
		out += "Duration:  " + formatTimeString(response.Duration) + "\n"
	}
	out += "Streams:  # All available quality"
	out = formatString(out, ":")
	for _, stream := range response.Stream {
		out += formatString(sprintMediaStream(stream), ":") + "\n"
	}
	fmt.Printf(out)

}

func sprintMediaStream(stream spider.Stream) (info string) {
	var out string
	out = fmt.Sprintf("    [%d]  -------------------\n", stream.ID)
	out += "    Quality:  " + stream.Quality + "\n"
	if stream.StreamType != "" {
		out += "    Type:  " + stream.StreamType + "\n"
	}
	if stream.Size != 0 {
		out += fmt.Sprintf("    Size:  %s (%d Bytes)\n", formatSize(int64(stream.Size)), stream.Size)
	}
	if stream.Height != 0 && stream.Width != 0 {
		out += fmt.Sprintf("    Window:  %dx%d\n", stream.Width, stream.Height)
	}
	out += fmt.Sprintf("    # download with: -f %d ...\n", stream.ID)
	return out
}
