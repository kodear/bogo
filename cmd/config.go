package cmd

import (
	"github.com/zhxingy/bogo/config"
	"github.com/zhxingy/bogo/cookie"
)

func ImportCookie(file string) (err error) {
	ok, err := pathExists(file)
	if !ok {
		return
	}

	jar, err := cookie.Export(file)
	if err != nil {
		return
	}

	cfg := config.Open("")
	cfg.Config.Cookies = jar
	cfg.Write()
	cfg.Close()
	return nil
}

func SetDownloadPath(path string) {
	cfg := config.Open("")
	cfg.Config.DownloadPath = path
	cfg.Write()
	cfg.Close()
}
