## examples

```
package main

import (
	"encoding/json"
	"fmt"
	"github.com/zhxingy/bogo/spider"
)

func main() {
	response, err := spider.Do("https://www.acfun.cn/v/ac15533372", nil, nil)
	if err != nil{
		panic(err)
	}

	body, _ := json.MarshalIndent(response, "", "  ")
	fmt.Println(string(body))
}

```