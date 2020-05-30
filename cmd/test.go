package cmd


import (
	"fmt"
	"github.com/zhxingy/bogo/selector"
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

	//panic(exception.HTTPHtmlException(errors.New("test error")))
	//s := selector.Create([]byte("abcdefg"))
	var s selector.Selector
	s = []byte("abcdefg")

	fmt.Println(s.String())
	var x, y, z string
	err := s.Re(`(a).*(e).*(g)`, &x, &y, &z)
	if err != nil{
		fmt.Println(err)
	}
	fmt.Println(x, y, z)
}
