package spiders

import (
	"fmt"
	"strings"
)

func SplitCookies(cookie string) (cookies map[string]string) {
	cookies = make(map[string]string)
	for _, s := range strings.Split(cookie, ";") {
		c := strings.Split(s, "=")
		if len(c) > 1 {
			cookies[strings.Replace(c[0], " ", "", 1)] = c[1]
		} else if len(c) == 0 {
			cookies[strings.Replace(c[0], " ", "", 1)] = ""
		}
	}
	return
}

func JoinCookies(cookies map[string]string) (cookie string) {
	for key, value := range cookies {
		cookie += fmt.Sprintf("%s=%s; ", key, value)
	}
	return
}
