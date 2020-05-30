package spiders

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"
)

func hasMatch(r, pattern string) bool {
	match, err := regexp.MatchString(pattern, r)
	if err != nil {
		return false
	} else {
		return match
	}
}

func splitCookies(cookie string) (cookies map[string]string) {
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

type Sort struct {
	v    []Response
	less func(x, y Response) bool
}

func (s Sort) Len() int {
	return len(s.v)
}

func (s Sort) Less(i, j int) bool {
	return s.less(s.v[i], s.v[j])
}

func (s Sort) Swap(i, j int) {
	s.v[i], s.v[j] = s.v[j], s.v[i]
}

type Parse struct {
	Bytes  []byte
	String string
}

func NewParse(b []byte) *Parse {
	return &Parse{
		Bytes:  b,
		String: string(b),
	}
}

func (p *Parse) Json(data interface{}) (err error) {
	err = json.Unmarshal(p.Bytes, &data)
	return
}

func (p *Parse) Search(pattern string) (strings [][]string, err error) {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return
	}

	strings = regex.FindAllStringSubmatch(p.String, -1)
	return
}

func cookieName(key string, cookieJar []*http.Cookie) (value string) {
	for _, cookie := range cookieJar {
		if cookie.Name == key {
			value = cookie.Value
			break
		}
	}
	return
}
