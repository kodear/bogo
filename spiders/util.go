package spiders

import "regexp"

func match(uri, expression string) bool {
	ok, err := regexp.MatchString(expression, uri)
	if err != nil {
		return false
	} else {
		return ok
	}
}


