package main

import (
	"flag"
	"fmt"
	"github.com/zhxingy/bogo/download"
	"github.com/zhxingy/bogo/spiders"
	"gopkg.in/cheggaaa/pb.v1"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"
)

const Version = "0.0.5"

var Cookies = map[string]string{
	"bilibili": "SESSDATA=8a6276a5%2c1599717086%2c352e5*31",
	"iqiyi":    `QP001=1; QP0017=100; QP0018=100; QC005=a985d5b3de4fa3a8ff19ffa1039b58c9; QC006=ma38kvbmlpyuvd84a3ddzazn; QC008=1562227263.1584082052.1584164957.37; Hm_lvt_53b7374a63c37483e5dd97d78d9bb36e=1583401476,1583405527,1584082053,1584164958; T00404=cc6c887162fde433fadbd4d8537684f8; IMS=IggQARj__rPzBSokCiAxYzdkYWRmNmU2ZjI3MDVjZmUxYmZhZDQwZWUxY2M3NhAAciQKIDFjN2RhZGY2ZTZmMjcwNWNmZTFiZmFkNDBlZTFjYzc2EAA; QC173=0; P00004=.1564972346.62524a52ea; QP008=720; GC_PCA=false; cuuid=f6734c56000846688405c36ddcbd193f; QP007=0; P111114=1574596768; P1111129=1574596798; P00037=A00000; P00039=64r9C0LCP09ZXnScm23NqSpm3x7V0dSm3R8GhSGxnZyyVeXclc6w16YYJEm2bCNnNLF3r0a5; QC021=%5B%7B%22key%22%3A%22%E6%B2%89%E7%9D%A1%E9%AD%94%E5%92%922%22%7D%5D; QC124=1%7C0; _ga=GA1.2.482661439.1574699249; QC160=%7B%22u%22%3A%22%22%2C%22lang%22%3A%22%22%2C%22local%22%3A%7B%22name%22%3A%22%E4%B8%AD%E5%9B%BD%E5%A4%A7%E9%99%86%22%2C%22init%22%3A%22Z%22%2C%22rcode%22%3A48%2C%22acode%22%3A86%7D%2C%22type%22%3A%22s1%22%7D; P00001=8ecTZzDFm1Xm1cZC1N5JTDrys2sB3zuetXClezpxIUHMVQlMxHiXID78QO4SehSQ9F6gc5; P00003=1672781693; P00010=1672781693; P01010=1584115200; P00007=8ecTZzDFm1Xm1cZC1N5JTDrys2sB3zuetXClezpxIUHMVQlMxHiXID78QO4SehSQ9F6gc5; P00PRU=1672781693; P00002=%7B%22uid%22%3A1672781693%2C%22pru%22%3A1672781693%2C%22user_name%22%3A%2218565861644%22%2C%22nickname%22%3A%22%5Cu7231%5Cu6e05%5Cu89c9%5Cu7f57%5Cu2022D%5Cu2022%5Cu5c3c%5Cu53e4%5Cu62c9%5Cu65af%5Cu2022%5Cu8d75%5Cu56db%22%2C%22pnickname%22%3A%22%5Cu7231%5Cu6e05%5Cu89c9%5Cu7f57%5Cu2022D%5Cu2022%5Cu5c3c%5Cu53e4%5Cu62c9%5Cu65af%5Cu2022%5Cu8d75%5Cu56db%22%2C%22type%22%3A11%2C%22email%22%3A%22%22%7D; P000email=""; QC170=1; QP0013=1; QC179=%7B%22userIcon%22%3A%22https%3A//img7.iqiyipic.com/passport/20190703/98/2f/passport_1672781693_156209173635049_130_130.png%22%2C%22vipTypes%22%3A%221%22%7D; QYABEX={"pcw_home_movie":{"value":"new","abtest":"146_B"}}; QC175=%7B%22upd%22%3Atrue%2C%22ct%22%3A1584082054472%7D; QY_PUSHMSG_ID=a985d5b3de4fa3a8ff19ffa1039b58c9; websocket=true; QC163=1; QC159=%7B%22color%22%3A%22FFFFFF%22%2C%22channelConfig%22%3A1%2C%22hadTip%22%3A1%2C%22hideRoleTip%22%3A1%2C%22isOpen%22%3A1%2C%22speed%22%3A10%2C%22density%22%3A30%2C%22opacity%22%3A86%2C%22isFilterColorFont%22%3A1%2C%22proofShield%22%3A0%2C%22forcedFontSize%22%3A24%2C%22isFilterImage%22%3A1%7D; QY00001=1672781693; T00700=EgcI9L-tIRAB; QC007=DIRECT; QC010=206649623; nu=0; Hm_lpvt_53b7374a63c37483e5dd97d78d9bb36e=1584164961; __dfp=a152e24d6d70c6489abb3418a3c9c7dd66e4d92f68d66abc0a8ba79eb55e56cd2a@1584176201865@1582880202865`,
	"youku":    `__ysuid=1562657576484g7V; __arlft=1565336440; cna=0h+eFVk7jUoCAXFaImmt/21r; juid=01di2dk1imtr7; ysestep=2; yseidcount=14; ystep=37; isg=BMLCuNjUKPHZVTeQIXI9k22CEMgkk8at6xiFkgzb3jXuX2PZ9CGDvA3VC9sjFD5F; UM_distinctid=16d48c4ffe3a-00adba159cc9948-4c312373-1fa400-16d48c4ffe4bd9; user_name=%E5%BC%A0%E9%91%AB98222; modalFrequency={"UUID":"2"}; modalBlackListclose={"UUID":"2"}; modalBlackListlogined={"UUID":"2","blackListSrc":"logined"}; P_gck=NA%7Cmiev12TDYqw%2BMFG%2F0wMF7w%3D%3D%7CNA%7C1583416335712; P_pck_rm=4swDIg3Od16f4da0ebcd41ZBcrAoyMq9ceW6Tgi1qSkOv9yDiByryR9CcZoB4pYd7TRz9RowUPfsPdwmKmeTsanRap%2BmbvPsSdZWfMz3aoFHJS9zsYRKguRpRvKztFWm4r%2FRAJpNU2yqs0w9wGYolmOemoKDjpn8dNaXdQ%3D%3D%5FV2; __ayft=1584165175448; __aysid=1584165175448SLN; __arpvid=1584165203432DQz7Le-1584165203451; __ayscnt=1; __aypstp=2; __ayspstp=2; _m_h5_tk=53e53e8aad73467d8df4e20722a88076_1584170215798; _m_h5_tk_enc=931ba391ccf3ec06d63d735b2543ceed; P_ck_ctl=A9B8329412002663DD5FF63853DA362B; _m_h5_c=6fc499fed6f4306c85bf88fdf8c4bb2f_1584175256158%3B8c2b69e9583c062a6607b3e613463fe4; __arycid=dv-3-00; __arcms=dv-3-00; __ayvstp=1; __aysvstp=1`,
	"tencent":  `cm_cookie=V1,10016&G1LIOs21cjIy&AQEB4qeequdWN2FkcyHIGffFKty_f5OkHY8-&190709&190709,10027&1562661248346136&AQEBTroe_Qe5Tx-wGMYT0CCsgFodjs-NS1kN&190709&190709,10012&emjZ9KRn&AQEBYhNbSkQLK-G_UQDuLS5A43Mtf2LhCTzz&190918&190918,10008&0641A9A8A7E04EEBA7010F6216BB2E24-&AQEByUv60O2hvHUpAcyyDviCYn99UABAmcEp&190709&191125,110061&b1766ee24036381&AQEBYhNbSkQLK-FeojWy3VGSpWeHd0UrOVSA&191126&191126,110065&XK3pSt6ZHV&AQEBYhNbSkQLK-Gr0vG-oGirNo1NLMZ-xw9v&190810&200228,110066&csFlh08qkh10&AQEBYhNbSkQLK-E7Xa1WCj33j_f7NZYFUZtk&190810&200314; _ga=GA1.2.1762018758.1563275273; pgv_pvid=3791155328; pgv_pvi=2601128960; RK=sehd5mmMTr; ptcz=776fa43d26f5d5328f9020b1e225bd17aea8c369e729ca21f741932cf37ad92c; sensorsdata2015jssdkcross=%7B%22distinct_id%22%3A%2216c756b73502c2-0308a90f9b69d88-4c312272-2073600-16c756b7351a9a%22%2C%22%24device_id%22%3A%2216c756b73502c2-0308a90f9b69d88-4c312272-2073600-16c756b7351a9a%22%2C%22props%22%3A%7B%22%24latest_traffic_source_type%22%3A%22%E7%9B%B4%E6%8E%A5%E6%B5%81%E9%87%8F%22%2C%22%24latest_referrer%22%3A%22%22%2C%22%24latest_referrer_host%22%3A%22%22%2C%22%24latest_search_keyword%22%3A%22%E6%9C%AA%E5%8F%96%E5%88%B0%E5%80%BC_%E7%9B%B4%E6%8E%A5%E6%89%93%E5%BC%80%22%7D%7D; tvfe_boss_uuid=d236fd0662d67a42; appuser=A72C8DA118EFBE2D; o_minduid=TcWw_L9EytiuntXT7JYBCunrIpN16FdP; ufc=r47_1_1569410303_1569324083; ptui_loginuin=641015302; lv_play_indexl.=89; pgv_info=ssid=s3747265688; uid=99954938; psessionid=716eba25_1584166076_0_28796; psessiontime=1584166081; adid=641015302; exuid=n4HV0XB-waJNz94U87FipA%3D%3D; Lturn=210; LPLFturn=399; LKBturn=39; LPVLturn=711; LVMturn=523; LPSJturn=576; LBSturn=609; LZIturn=919; LZCturn=768; LPCZCturn=170; LCZCturn=821; LVINturn=570`,
	"mgtv":     `locale=CHN; WWW_LOCALE=CN; __random_seed=0.8574044248094317; mba_deviceid=708d7311-baf8-5750-5c0e-94178f757da7; mba_cxid_expiration=1584028800000; mba_cxid=8kpqp21m4cm; sessionid=1583986048561_8kpqp21m4cm; MQGUID=1237952993639518208; __MQGUID=1237952993639518208; pc_v6=v6; id=52531462; rnd=rnd; seqid=bpkrb5hlqhggi8ohjs20; uuid=e6c27eb1894245deaae5ae61e66f9958; vipStatus=3; wei=7e4799e8757e05a9b16bebe9da429d6c; wei2=439bXDGf69cfqDyTUfB3hwROpjABBhNctRRC4a7V2ehAG7lm63xQ1374hK74j9AvZTPEnfYSVsBe2hTBUdjVf9Y5Ue0dGThrwOsbFNrwNOOjpupcEjDaZju0g00MBJD%2F3wVFHHW3rVkKkGwMRKBWccA%2BoPpJ8CbNpncBnMuIfRmVFjIgKkwqcCoTqSM93VdtjVbb94W7RCkqo9A; HDCN=BPKRB5HPAHH7767E04BG-864042134; PM_CHKID=0b9a503289232a45; mba_sessionid=4c07928f-ee89-5745-dea4-5c05e56600d1; mba_last_action_time=1583993031036; beta_timer=1583993032019; lastActionTime=1583994066911`,
}

var (
	videoShow         bool
	videoURL          string
	videoQuality      string
	videoID           int
	videoDownloadFile string
	downloadHost      string
	downloadText      string
	downloadUrl       string
	downloadUrls      []string
	configPath        = "bogo"
	configName        = "bogo.ini"
	downloadPath      = "BogoDownloads"
	showWeb           bool
	version           bool
)

func init() {
	flag.BoolVar(&videoShow, "l", false, "output all video information and exit")
	flag.StringVar(&videoURL, "i", "", "url that needs to be parsed")
	flag.StringVar(&videoQuality, "q", "", "select the quality of the downloaded video")
	flag.IntVar(&videoID, "f", 0, "download the video by number")
	flag.StringVar(&videoDownloadFile, "o", "", "set download file save name")
	flag.StringVar(&downloadHost, "set-download-path", "", "set download file save path")
	flag.BoolVar(&showWeb, "s", false, "list supported parsing sites and log out")
	flag.BoolVar(&version, "v", false, "print the software version and exit")
}

func ConfigName() string {
	username, err := user.Current()
	if err != nil {
		panic(err)
	}

	configPath := filepath.Join(username.HomeDir, ".config", configPath)
	_, err = os.Stat(configPath)
	if err != nil && os.IsNotExist(err) {
		err := os.MkdirAll(configPath, 0644)
		if err != nil {
			panic(err)
		}
	}
	configFile := filepath.Join(configPath, configName)
	return configFile
}

func DownloadPath() string {
	username, err := user.Current()
	if err != nil {
		panic(err)
	}

	downloadRoot := filepath.Join(username.HomeDir, downloadPath)
	_, err = os.Stat(downloadRoot)
	if err != nil && os.IsNotExist(err) {
		err := os.MkdirAll(downloadRoot, 0644)
		if err != nil {
			panic(err)
		}
	}

	return downloadRoot

}
func main() {
	flag.Parse()

	// 加载配置文件
	configFile := ConfigName()
	downloadRoot := DownloadPath()
	cfg := NewConfig(configFile, downloadRoot, Cookies)
	cfg.Read()

	if downloadHost != cfg.root && downloadHost != "" {
		cfg.root = downloadHost
		cfg.Write()
		os.Exit(0)
	}

	if version {
		fmt.Printf("Bogo Version: %v\n", Version)
		os.Exit(0)
	}

	if showWeb {
		spiders.ShowWeb()
		os.Exit(0)
	}

	if videoURL == "" {
		flag.Usage()
		os.Exit(1)
	}

	cookies := make(map[string]string)
	for k, v := range cfg.cookies {
		cookies[k] = v
	}

	if videoShow {
		spiders.ShowVideo(videoURL, cookies)
		os.Exit(0)
	}

	video, err := spiders.Do(videoURL, videoQuality, videoID, cookies)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	// 初始化下载器错误
	downloader, err := download.LoadDownloader(video.DownloadProtocol)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	if video.DownloadProtocol == "hls" || video.DownloadProtocol == "http" {
		downloadUrl = video.Links[0].URL
	} else if video.DownloadProtocol == "hlsText" {
		downloadText = video.Links[0].URL
	} else if video.DownloadProtocol == "httpSegFlv" {
		for _, v := range video.Links {
			downloadUrls = append(downloadUrls, v.URL)
		}
	} else {
		fmt.Println("did not match to downloader")
		os.Exit(3)
	}

	if videoDownloadFile == "" {
		if video.Title == video.Part {
			video.Part = ""
		} else if video.Part != "" {
			video.Part = "-" + video.Part
		}
		videoDownloadFile = video.Title + video.Part + "." + video.Format

		// The system cannot find the path specified.
		// C:\Users\Administrator\Desktop\千年女子最强音《华夏巾帼志》【茶理理/小缘/肥皂菌/三畿道】.flv
		// windows 文件名中不能包含 \  /  :  *  ?  "  <  >  |
		if runtime.GOOS == "windows" {
			re := regexp.MustCompile(`\/|\\|\:|\*|\?|\"|\<|\>|\|`)
			videoDownloadFile = re.ReplaceAllString(videoDownloadFile, "、")
		} else {
			videoDownloadFile = strings.Replace(videoDownloadFile, `/`, `\/`, -1)
		}
	}

	DownloadFile := filepath.Join(cfg.root, videoDownloadFile)

	downloader.SetHeaders(video.DownloadHeaders)
	downloader.SetMax(video.Size)
	go downloader.Do(downloadUrl, downloadText, "", DownloadFile, downloadUrls)

	go func() {
		for {
			if downloader.Status() || downloader.Error() != nil {
				if downloader.Error() != nil {
					fmt.Println(downloader.Error())
					os.Exit(5)
				}
				if downloader.Chan() != nil {
					close(downloader.Chan())
				} else {
					fmt.Println("close of nil channel")
					os.Exit(5)
				}

				break
			}
			time.Sleep(1000)
		}
	}()

	// 等待获取进度条最大值
	for !downloader.Progress() {
		time.Sleep(1000)
	}

	bar := pb.New(downloader.Max()).SetRefreshRate(time.Millisecond * 10)
	if video.Size != 0 || video.DownloadProtocol == "http" {
		bar.SetUnits(pb.U_BYTES_DEC)
	}
	bar.ShowSpeed = true
	bar.ShowTimeLeft = true
	bar.ShowFinalTime = true
	bar.SetMaxWidth(200)
	bar.Prefix("Download: [" + path.Base(videoDownloadFile) + "]")
	bar.Start()

	//
	for !downloader.INIT() {
		time.Sleep(1000)
	}

	for p := range downloader.Chan() {
		bar.Add(p)
	}
	bar.Finish()
}
