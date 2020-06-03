package main

import (
	"github.com/zhxingy/bogo/cmd"
	"net/http"
)

func main()  {
	err := cmd.Download("https://www.acfun.cn/v/ac15927987", "", "./", nil, http.Header{}, 0, 1,)
	if err != nil{
		panic(err)
	}
}
