package downloader

type Downloader interface {
	start()
	Wait()
}

func Start(downloader Downloader) {
	go downloader.start()
}
