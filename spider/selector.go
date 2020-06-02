package spider

import (
	"encoding/json"
	"errors"
	"regexp"
)

type Selector []byte

func (cls *Selector) String() string {
	return string(*cls)
}

func (cls *Selector) Json(v interface{}) error {
	return json.Unmarshal(*cls, &v)
}

func (cls *Selector) Re(pattern string, args ...*string) (err error) {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return
	}

	strings := regex.FindAllStringSubmatch(string(*cls), -1)
	if len(strings) == 0 {
		return errors.New("no match was found")
	}

	if len(strings[0])-1 != len(args) {
		return errors.New("the matching result is inconsistent with the length of the output parameter")
	}

	for index, arg := range strings[0] {
		if index > 0 {
			*args[index-1] = arg
		}
	}

	return nil
}

func (cls *Selector) ReByJson(pattern string, v interface{}) (err error) {
	var arg string
	err = cls.Re(pattern, &arg)
	if err != nil {
		return ParseTextErr(err)
	} else {
		err = json.Unmarshal([]byte(arg), &v)
	}

	if err != nil {
		return ParseJsonErr(err)
	}

	return
}
