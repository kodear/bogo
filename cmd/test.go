package main

import (
	"encoding/json"
	"fmt"
	"github.com/zhxingy/bogo/spider"
)

func main() {
	response, err := spider.Do("https://www.bilibili.com/bangumi/play/ep15014", nil)
	if err != nil{
		panic(err)
	}

	body, _ := json.MarshalIndent(response, "", "\t")
	fmt.Println(string(body))

}
