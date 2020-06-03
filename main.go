package main

import (
	"fmt"
	"github.com/urfave/cli"
	"os"
)

func main(){
	//ie, err := spider.NewSpider("https://v.qq.com/x/cover/mzc00200q3ctjw0.html")
	//if err != nil{
	//	panic(err)
	//}
	//
	//ie.Initialization(nil, nil)
	//err = ie.Request()
	//if err != nil{
	//	panic(err)
	//}
	//
	//err = cmd.Download("", "./", 0, 720, ie.Response())
	//err := cmd.ImportCookie(`C:\Users\Administrator\AppData\Roaming\Mozilla\Firefox\Profiles\za7glfdv.default-release\cookies.sqlite`)
	//if err != nil{
	//	panic(err)
	//}

	//c := cli.NewCLI("bogo", cmd.Version)
	//c.Args = os.Args[1:]
	//c.Commands = map[string]cli.CommandFactory{
	//	"foo": fooCommandFactory,
	//	"bar": barCommandFactory,
	//}
	//
	//exitStatus, err := c.Run()
	//if err != nil {
	//	log.Println(err)
	//}
	//
	//os.Exit(exitStatus)
	//var language string
	var m_port int
	app := cli.NewApp()
	app.Name = "greet"           // 指定程序名称
	app.Usage = "say a greeting" //  程序功能描述

	app.Action = func(c *cli.Context) error {
		println("Greetings")
		fmt.Println(c.Int("port"))
		fmt.Println(m_port)
		return nil
	}

	app.Run(os.Args)

}
