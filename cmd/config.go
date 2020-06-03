package cmd

import (
	"github.com/zhxingy/bogo/config"
)

//func ImportCookie (file string)(err error){
//	jar, err := cookie.Export(file)
//	if err != nil{
//		return
//	}
//
//	cfg := config.Open("")
//	cfg.Config.Cookies = jar
//	cfg.Write()
//	cfg.Close()
//	return nil
//}

func SetDownloadPath(path string){
	cfg := config.Open("")
	cfg.Config.DownloadPath = path
	cfg.Write()
	cfg.Close()
}