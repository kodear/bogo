package download

type Download struct {
	Headers map[string]string
	Ch      chan int
	Len     int
	Ok      bool
	Err     error
	Init    bool
	C       bool
}

type Downloader interface {
	Do(link, text, file, fname string, links []string)
}
