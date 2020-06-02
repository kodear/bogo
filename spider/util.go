package spider

import (
	"crypto/md5"
	"encoding/hex"
	"regexp"
)

func match(uri, expression string) bool {
	ok, err := regexp.MatchString(expression, uri)
	if err != nil {
		return false
	} else {
		return ok
	}
}

func md5x(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
