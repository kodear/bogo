package cmd

import (
	"fmt"
	"github.com/zhxingy/bogo/spider"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func formatString(str, seq string) string {
	var newString string
	var n int
	for _, s := range strings.Split(str, "\n") {
		if s == "" || len(strings.Split(s, seq)) == 1 {
			continue
		}
		m := len(strings.Split(s, seq)[0])
		if m > n {
			n = m
		}
	}

	for _, s := range strings.Split(str, "\n") {
		if s == "" {
			continue
		}
		sliceString := strings.Split(s, seq)
		if len(sliceString) == 1 {
			newString += s + "\n"
			continue
		}
		key := sliceString[0]
		value := strings.Join(sliceString[1:], seq)
		x := n
		m := len(key)
		for m < x {
			value = " " + value
			x -= 1
		}
		newString += key + seq + value + "\n"
	}

	return newString
}

func formatTitle(s1, s2 string) string {
	if s1 == s2 || s2 == "" {
		return s1
	} else {
		return s1 + "-" + s2
	}
}

func formatTimeString(timestamp int) string {
	minutes := timestamp / 60
	seconds := timestamp % 60
	hours := minutes / 60
	minutes = minutes % 60

	if hours >= 24 {
		return ">1d"
	}

	return fmt.Sprintf("%dh%dm%ds", hours, minutes, seconds)
}

func formatTimeString2(timestamp int) string {
	minutes := timestamp / 60
	seconds := timestamp % 60
	hours := minutes / 60
	minutes = minutes % 60

	if hours >= 24 {
		return "--:--:--"
	}

	hoursString := strconv.Itoa(hours)
	minutesString := strconv.Itoa(minutes)
	secondsString := strconv.Itoa(seconds)
	if len(hoursString) < 2 {
		hoursString = "0" + hoursString
	}
	if len(minutesString) < 2{
		minutesString = "0" + minutesString
	}
	if len(secondsString) < 2{
		secondsString = "0" + secondsString
	}

	return fmt.Sprintf("%s:%s:%s", hoursString, minutesString, secondsString)
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
	} else {
		return fmt.Sprintf("%.2f EB", float64(fileSize)/float64(1024*1024*1024*1024*1024))
	}
}

func extractQuality(quality string) int {
	regex, _ := regexp.Compile(`(\d+)`)
	str := regex.FindAllStringSubmatch(quality, -1)
	if len(str) == 0 || len(str[0]) == 1 {
		return 0
	}

	result, _ := strconv.Atoi(str[0][1])
	return result
}

type streams []spider.Stream

func (stream streams) Len() int {
	return len(stream)
}

func (stream streams) Less(i, j int) bool {
	return extractQuality(stream[i].Quality) > extractQuality(stream[j].Quality)
}

func (stream streams) Swap(i, j int) {
	stream[i], stream[j] = stream[j], stream[i]
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true, nil
		}
		return false, err
	}
	return true, nil
}