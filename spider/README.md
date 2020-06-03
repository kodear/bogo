## examples

```
package main

import (
	"encoding/json"
	"fmt"
	"github.com/zhxingy/bogo/spider"
)

func main() {
	// get the extractor object
	cls, err := spider.NewSpider("https://www.acfun.cn/v/ac15918856")
	if err != nil {
		panic(err)
	}

	// initialize extractor
	cls.Initialization(nil, nil)
	err = cls.Request()

	// if the msg is not nil, the extraction of failure.
	if err != nil {
		panic(err)
	}

	body, _ := json.MarshalIndent(cls.Response(), "", "  ")
	fmt.Println(string(body))
}

```