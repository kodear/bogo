package main

import (
	"fmt"
	"github.com/urfave/cli"
	"github.com/zhxingy/bogo/cmd"
	"github.com/zhxingy/bogo/config"
	"github.com/zhxingy/bogo/spider"
	"net/http"
	"os"
	"strings"
)

func main() {

	app := cli.NewApp()
	app.Name = "bogo"
	app.Usage = "bug灰常多的媒体下载器~"
	app.Version = cmd.Version
	app.HelpName = "bogo"
	app.HideHelp = false
	cli.HelpFlag = &cli.BoolFlag{
		Name:        "help",
		Aliases:     []string{"h"},
		Usage:       "显示帮助信息",
		DefaultText: "关闭",
	}
	cli.VersionFlag = &cli.BoolFlag{
		Name:        "version",
		Aliases:     []string{"v"},
		Usage:       "显示版本信息",
		DefaultText: "关闭",
	}

	//app.HideHelp = false
	//app.UsageText = "bogo [-h|-v|-i <url> [-H <header>|-q <quality>|-f <fid>|-o <out>]|config {import-cookie <cooke_file>|set-download-path <download_path>}]"
	app.Commands = []*cli.Command{
		{
			Name:     "config",
			Usage:    "运行环境配置",
			Category: "config",
			Subcommands: []*cli.Command{
				{
					Name:  "import-cookie",
					Usage: "导入浏览器cookie到本地 (仅支持firefox)",
					Action: func(context *cli.Context) error {
						if context.Args().First() == "" {
							return nil
						}

						err := cmd.ImportCookie(context.Args().First())
						if err != nil {
							//fmt.Printf("import cookie file failed. err msg: %v\n", err )
							return err
						}

						return nil
					},
				},
				{
					Name:  "set-download-path",
					Usage: fmt.Sprintf("设置下载文件保存路径 (默认: %s)", config.DefaultDownloadPath()),
					Action: func(context *cli.Context) error {
						if context.Args().First() == "" {
							return nil
						}
						cmd.SetDownloadPath(context.Args().First())
						return nil
					},
				},
				{
					Name:  "help",
					Usage: "显示一个命令的命令或帮助列表",
				},
			},
		},
		//{
		//	Name: "help, h",
		//	Usage:    "显示帮助信息",
		//},
	}
	app.Flags = []cli.Flag{
		&cli.StringFlag{Name: "input", Usage: "指定解析媒体地址 (支持m3u8/http媒体直链)", Aliases: []string{"i"}},
		&cli.StringFlag{Name: "out", Usage: "指定媒体文件名", Aliases: []string{"o"}, DefaultText: ""},
		&cli.StringFlag{Name: "header", Usage: "添加请求请求头 (多个http-headers设置是采用CRLF来进行分割)", Aliases: []string{"H"}, DefaultText: ""},
		&cli.IntFlag{Name: "quality", Usage: "指定媒体质量下载", Aliases: []string{"q"}, DefaultText: "720"},
		&cli.IntFlag{Name: "fid", Usage: "指定媒体编号下载", Aliases: []string{"f"}, DefaultText: "nil"},
		&cli.BoolFlag{Name: "print", Usage: "列出所有可下载的媒体并退出", Aliases: []string{"p"}, DefaultText: "关闭"},
	}
	app.Action = func(context *cli.Context) error {
		url := context.String("input")
		sprint := context.Bool("print")
		quality := context.Int("quality")
		fid := context.Int("fid")
		out := context.String("out")
		headers := context.String("header")

		if url == "" {
			//fmt.Println(app.UsageText)
			return nil
		}

		cfg := config.Open("")
		ie, err := spider.NewSpider(url)
		if err != nil {
			//fmt.Println(err)
			return err
		}

		header := http.Header{}
		if headers != "" {
			for _, h := range strings.Split(headers, ";") {
				x := strings.Split(h, "=")
				if len(x) < 2 {
					continue
				}
				if len(header[strings.TrimSpace(x[0])]) > 0 {
					header[strings.TrimSpace(x[0])] = append(header[strings.TrimSpace(x[0])], strings.TrimSpace(x[1]))
				} else {
					header[strings.TrimSpace(x[0])] = []string{strings.TrimSpace(x[1])}
				}
			}
		}

		ie.Initialization(cfg.Config.Cookies[ie.Meta().Cookie.Name], header)
		err = ie.Request()
		if err != nil {
			//fmt.Println(err)
			return err
		}

		if sprint {
			cmd.PrintMedia(ie.Response())
			return nil
		}

		if quality == 0 {
			quality = 720
		}

		err = cmd.Download(out, cfg.Config.DownloadPath, fid, quality, ie.Response())
		if err != nil {
			//fmt.Println(err)
			return err
		}

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		//os.Exit(1)
	}
	if len(os.Args) < 2 {
		fmt.Println("usage: bogo -i <url> [-p|-H <header>|...]")
		//os.Exit(2)
	}
}
