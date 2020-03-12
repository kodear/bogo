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

const Version = "0.0.4"

var Cookies = map[string]string{
	"bilibili": "SESSDATA=9371411d%2C1585468536%2C0e682721",
	"iqiyi":    `QP001=1; QP0017=100; QP0018=100; QC005=a985d5b3de4fa3a8ff19ffa1039b58c9; QC006=ma38kvbmlpyuvd84a3ddzazn; QC008=1562227263.1583401476.1583405527.35; Hm_lvt_53b7374a63c37483e5dd97d78d9bb36e=1583317628,1583382067,1583401476,1583405527; T00404=cc6c887162fde433fadbd4d8537684f8; IMS=IggQBRj_w4TzBSokCiBjNTRkYWU1ZTg0MGIzNGI2YWZhOWRjMmZiMDBlZjg4NBAAMAAwAA; QC173=0; P00004=.1564972346.62524a52ea; QP008=2040; GC_PCA=false; cuuid=f6734c56000846688405c36ddcbd193f; QP007=0; P111114=1574596768; P1111129=1574596798; P00037=A00000; P00039=64r9C0LCP09ZXnScm23NqSpm3x7V0dSm3R8GhSGxnZyyVeXclc6w16YYJEm2bCNnNLF3r0a5; QC021=%5B%7B%22key%22%3A%22%E6%B2%89%E7%9D%A1%E9%AD%94%E5%92%922%22%7D%5D; QC124=1%7C0; _ga=GA1.2.482661439.1574699249; QC160=%7B%22u%22%3A%22%22%2C%22lang%22%3A%22%22%2C%22local%22%3A%7B%22name%22%3A%22%E4%B8%AD%E5%9B%BD%E5%A4%A7%E9%99%86%22%2C%22init%22%3A%22Z%22%2C%22rcode%22%3A48%2C%22acode%22%3A86%7D%2C%22type%22%3A%22s1%22%7D; P00001=24lqd9pm3oIWG30lf57PWwFfSBd9VRDXIjJw1TzWm29IDQJI0Q75a6zOK4LJIspSaDqM4c; P00003=1672781693; P00010=1672781693; P01010=1583424000; P00007=24lqd9pm3oIWG30lf57PWwFfSBd9VRDXIjJw1TzWm29IDQJI0Q75a6zOK4LJIspSaDqM4c; P00PRU=1672781693; P00002=%7B%22uid%22%3A1672781693%2C%22pru%22%3A1672781693%2C%22user_name%22%3A%2218565861644%22%2C%22nickname%22%3A%22%5Cu7231%5Cu6e05%5Cu89c9%5Cu7f57%5Cu2022D%5Cu2022%5Cu5c3c%5Cu53e4%5Cu62c9%5Cu65af%5Cu2022%5Cu8d75%5Cu56db%22%2C%22pnickname%22%3A%22%5Cu7231%5Cu6e05%5Cu89c9%5Cu7f57%5Cu2022D%5Cu2022%5Cu5c3c%5Cu53e4%5Cu62c9%5Cu65af%5Cu2022%5Cu8d75%5Cu56db%22%2C%22type%22%3A11%2C%22email%22%3A%22%22%7D; P000email=""; QC170=1; QP0013=1; QC179=%7B%22userIcon%22%3A%22https%3A//img7.iqiyipic.com/passport/20190703/98/2f/passport_1672781693_156209173635049_130_130.png%22%2C%22vipTypes%22%3A%221%22%7D; QYABEX={"pcw_home_movie":{"value":"new","abtest":"146_B"}}; QC175={%22upd%22:true%2C%22ct%22:1583405529307}; QY_PUSHMSG_ID=a985d5b3de4fa3a8ff19ffa1039b58c9; websocket=true; QC176=%7B%22state%22%3A0%2C%22ct%22%3A1583382179138%7D; QP0010=1; T00700=EgcI9L-tIRAB; PCAU=0; QP009=1; Hm_lpvt_53b7374a63c37483e5dd97d78d9bb36e=1583405531; QC007=DIRECT; QC010=177035425; nu=0; QC163=1; CM0001=1; QY00001=1672781693; QILINPUSH=1; TQC002=type%3Djspfmc140109%26pla%3D11%26uid%3Da985d5b3de4fa3a8ff19ffa1039b58c9%26ppuid%3D1672781693%26brs%3Dff%26pgtype%3Dplay%26purl%3Dhttps%3A%252F%252Fwww.iqiyi.com%252Fv_19rwfbn2j4.html%3Fvfrm%253Dpcw_home%2526vfrmblk%253DCZ%2526vfrmrst%253D712211_cainizaizhui_image4%26cid%3D1%26tmplt%3D%26tm1%3D3543%2C0%26tm13%3D4565%26tm6%3D7371%2C0; __dfp=a152e24d6d70c6489abb3418a3c9c7dd66e4d92f68d66abc0a8ba79eb55e56cd2a@1584176201865@1582880202865`,
	"youku":    `__ysuid=1562657576484g7V; __arlft=1565336440; cna=0h+eFVk7jUoCAXFaImmt/21r; juid=01di2dk1imtr7; ysestep=2; yseidcount=14; ystep=37; isg=BOXl1k6m14ZDIzBpSkMyeubj96EfIpm06IiGq-fK6pw0_gRwr3GUhIWIiuKIZbFs; UM_distinctid=16d48c4ffe3a-00adba159cc9948-4c312373-1fa400-16d48c4ffe4bd9; user_name=%E5%BC%A0%E9%91%AB98222; _m_h5_tk=3bb8a272ae83f1b151ac298e86982275_1583489770436; _m_h5_tk_enc=8d7c0969ddc4f89efa5c53c497d62ee2; __aysid=1583415790629H7o; __ayspstp=25; modalFrequency={"UUID":"2"}; modalBlackListclose={"UUID":"2"}; modalBlackListlogined={"UUID":"2","blackListSrc":"logined"}; __aysvstp=17; P_gck=NA%7Cmiev12TDYqw%2BMFG%2F0wMF7w%3D%3D%7CNA%7C1583416335712; P_pck_rm=4swDIg3Od16f4da0ebcd41ZBcrAoyMq9ceW6Tgi1qSkOv9yDiByryR9CcZoB4pYd7TRz9RowUPfsPdwmKmeTsanRap%2BmbvPsSdZWfMz3aoFHJS9zsYRKguRpRvKztFWm4r%2FRAJpNU2yqs0w9wGYolmOemoKDjpn8dNaXdQ%3D%3D%5FV2; P_ck_ctl=42E8B4993141D71F32B96637585B5AAF; __ayft=1583467796112; __arpvid=15834850912405L2txz-1583485091288; __ayscnt=1; __aypstp=15; _m_h5_c=cf17683b25d3c4adf4d1aa98b5c65fba_1583494091015%3Bdfe9144052332002ab8ebdd1001c3e3c; __arycid=dc-3-00; __arcms=dc-3-00; __ayvstp=12`,
	"tencent":  `_ga=GA1.2.1762018758.1563275273; pgv_pvid=3791155328; pgv_pvi=2601128960; RK=sehd5mmMTr; ptcz=776fa43d26f5d5328f9020b1e225bd17aea8c369e729ca21f741932cf37ad92c; sensorsdata2015jssdkcross=%7B%22distinct_id%22%3A%2216c756b73502c2-0308a90f9b69d88-4c312272-2073600-16c756b7351a9a%22%2C%22%24device_id%22%3A%2216c756b73502c2-0308a90f9b69d88-4c312272-2073600-16c756b7351a9a%22%2C%22props%22%3A%7B%22%24latest_traffic_source_type%22%3A%22%E7%9B%B4%E6%8E%A5%E6%B5%81%E9%87%8F%22%2C%22%24latest_referrer%22%3A%22%22%2C%22%24latest_referrer_host%22%3A%22%22%2C%22%24latest_search_keyword%22%3A%22%E6%9C%AA%E5%8F%96%E5%88%B0%E5%80%BC_%E7%9B%B4%E6%8E%A5%E6%89%93%E5%BC%80%22%7D%7D; tvfe_boss_uuid=d236fd0662d67a42; video_guid=f95f19357edb465c; video_platform=2; ptui_loginuin=641015302; main_login=qq; vqq_access_token=D5AA7CEB9FDB0A73BF999EF063F6752C; vqq_appid=101483052; vqq_openid=ED87B84A15549ED3E60D6D2F927314D6; vqq_vuserid=119772433; vqq_vusession=-KVfOFOwyDcffaDKqgFceQ..; vqq_refresh_token=8EC653149DAF690299935F78CEBBFEFC; vqq_next_refresh_time=1399; vqq_login_time_init=1583483345; pgv_info=ssid=s2253475206; login_time_last=2020-3-6 16:29:7; uid=99954938`,
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
	if video.Size != 0 {
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
