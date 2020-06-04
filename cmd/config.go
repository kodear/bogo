package cmd

import (
	"github.com/zhxingy/bogo/config"
	"github.com/zhxingy/bogo/cookie"
	"net/http"
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

	newJar := cookie.CookiesJar{}
	for key, cookies := range jar {
		var newCookies []*http.Cookie
		for _, v := range cookies {
			var ok bool
			for i := 0; i < len(v.Value); i++ {
				if validCookieValueByte(v.Value[i]) {
					continue
				} else {
					ok = true
					break
				}
			}
			if !ok {
				newCookies = append(newCookies, &http.Cookie{Name: v.Name, Value: v.Value})
			}
		}
		newJar[key] = newCookies
	}

	cfg := config.Open("")
	cfg.Config.Cookies = newJar
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
