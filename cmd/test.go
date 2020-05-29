package main

import (
	"errors"
	"github.com/zhxingy/bogo/exception"
)

func main()  {
	//x, err := cookie.Export(`C:\Users\A\AppData\Roaming\Mozilla\Firefox\Profiles\y1gkczdt.default-release\cookies.sqlite`)
	//if err != nil{
	//	fmt.Println(err)
	//}
	//cfg := config.Open("")
	//cfg.Config.Proxy = "http://11.11.11.11"
	////cfg.Config.DownloadPath = "./"
	//cfg.Config.Cookies = x
	////fmt.Println(cfg.Config.Cookies)
	//cfg.Write()
	//cfg.Close()

	panic(exception.HTTPHtmlException(errors.New("test error")))
}
