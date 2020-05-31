package cookie

import (
	"github.com/zhxingy/bogo/spider"
	"net/http"
	"strings"
)

type Dumper interface {
	Run()([]*http.Cookie, error)
}

type Dump struct {
	file string
}

type CookiesJar map[string][]*http.Cookie

func Export(file string)(CookiesJar, error){
	dump := firefoxDump{Dump{file:file}}
	cookies, err := dump.Run()
	if err != nil{
		return nil, err
	}

	hostNames := make(map[string][]string)
	for _, x := range spider.Spiders {
		if x.Meta().Cookie.Enable{
			hostNames[x.Meta().Cookie.Name] = x.Meta().Cookie.Domain
		}
	}

	cookiesJar := make(CookiesJar)
	for id, hosts := range hostNames{
		for _, host := range hosts{
			for _, cookie := range cookies{
				if strings.Contains(cookie.Domain, host){
					if len(cookiesJar[id]) == 0{
						cookiesJar[id] = []*http.Cookie{cookie}
					}else{
						cookiesJar[id] = append(cookiesJar[id], cookie)
					}
				}
			}
		}
	}

	return cookiesJar, nil
}