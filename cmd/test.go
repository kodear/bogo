package main

import (
	"encoding/json"
	"fmt"
	"github.com/zhxingy/bogo/spider"
)

func main() {

	response, err := spider.Do("http://www.bilibili.com/index.mp4/1", nil)
	if err != nil {
		panic(err)
	}

	body, _ := json.MarshalIndent(response, "", "\t")
	fmt.Println(string(body))

}
