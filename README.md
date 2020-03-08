# Bogo 主流网站视频下载器

### 支持网站
* [acfun](https://www.acfun.cn/)
* [acfun番剧](https://www.acfun.cn/v/list155/index.htm)
* [哔哩哔哩](https://www.bilibili.com/)
* [哔哩哔哩番剧](https://www.bilibili.com/anime/)
* [爱奇艺](https://www.iqiyi.com/)
* [优酷](https://www.youku.com/)
* [腾讯视频](https://v.qq.com/)

### 下载地址
* [bogo-v0.0.2-windows-x64.zip](http://pkg.maoq.pw/pkg/bogo/bogo-v0.0.2-windows-x64.zip)
* [bogo-v0.0.2-linux-amd64.zip](http://pkg.maoq.pw/pkg/bogo/bogo-v0.0.2-linux-amd64.zip)

### 安装
* windows操作系统
    1. 下载软件包
    2. 打开文件所在目录
    3. 解压软件至用户家目录
* linux操作系统
```
$ cd /usr/src
$ wget "http://pkg.maoq.pw/pkg/bogo/bogo-v0.0.2-linux-amd64.zip"
$ unzip bogo-v0.0.2-linux-amd64.zip
$ chmox +x bogo-v0.0.2-linux-amd64/bogo
$ ln -s /usr/src/bogo-v0.0.2-linux-amd64/bogo /usr/bin/
$ bogo -v
Bogo Version: 0.0.2
```

### 使用

#### 帮助信息
```
$ bogo
Usage of bogo:
  -f int
    	download the video by number                                    // 指定视频编号
  -i string
    	url that needs to be parsed                                     // 视频网址
  -l	output all video information and exit                           // 列出所有可下载视频并退出
  -o string
    	set download file save name                                     // 设置文件名
  -q string
    	select the quality of the downloaded video                      // 指定视频清晰度
  -s	list supported parsing sites and log out                        // 查看支持解析网站
  -set-download-path string                                             //设置下载文件保存目录
    	set download file save path
  -v	print the software version and exit                              //查看软件版本号
```

#### 设置下载文件保存目录
```
// 默认保存目录为 $HOME/BogoDownloads
// bogo --set-download-path $HOME/BogoDownloads
$ bogo --set-download-path /home/data/video
```

#### 查看支持解析网站
```
$ bogo -s

Acfun【https://www.acfun.cn/】
Acfun番剧【https://www.acfun.cn/v/list155/index.htm】
哔哩哔哩【https://www.bilibili.com/】
哔哩哔哩番剧【https://www.bilibili.com/anime/】
爱奇艺【https://www.iqiyi.com/】
腾讯视频【https://v.qq.com/】
优酷【https://www.youku.com/】
```

#### 查看软件版本号
```
$ bogo -v
Bogo Version: 0.0.2
```

#### 列出所有可下载视频
```
$ bogo -i "https://v.qq.com/x/cover/m5zzglrbt5zdv6d.html" -l
-----------  --------------------  ---------  -----------  ---------------  ------------  ----------  -----------  -------------  ---------------  ---------------------
        ID                 Title       Part       Format       StreamType       Quality       Width       Height       Duration             Size       DownloadProtocol
-----------  --------------------  ---------  -----------  ---------------  ------------  ----------  -----------  -------------  ---------------  ---------------------
    321004       叶问4(普通话版)          -          mp4              fhd         1080P        1920          868           6419       1766372424                    hls

    321003       叶问4(普通话版)          -          mp4              shd          720P        1280          578           6419        979091028                    hls

    321002       叶问4(普通话版)          -          mp4               hd          480P         864          390           6419        574805300                    hls

    321001       叶问4(普通话版)          -          mp4               sd          270P         480          218           6419        259170784                    hls
-----------  --------------------  ---------  -----------  ---------------  ------------  ----------  -----------  -------------  ---------------  ---------------------
```

#### 按照默认设置下载视频
```
// 默认为720P, 当720P不存在时会尽可能向上匹配
// 当地址匹配到多个720P视频时,默认取第一个
// 精密控制请与 "-f" 结合使用
$ bogo -i "https://v.qq.com/x/cover/m5zzglrbt5zdv6d.html"
Download: [叶问4(普通话版).mp4] 619.37 MB / 979.09 MB [===================================================================>---------------------------------------]  63.26% 24.98 MB/s 00m13s

```

#### 指定视频编号下载
```
// 如果视频编号不存在则忽略
$ bogo -i "https://v.qq.com/x/cover/m5zzglrbt5zdv6d.html" -f 321004
Download: [叶问4(普通话版).mp4] 90.30 MB / 1.77 GB [=====>--------------------------------------------------------------------------------------------------------]   5.11% 28.17 MB/s 00m58s
```

#### 指定视频清晰度下载
```
// 如果存在两个相同的清晰度, 默认取第一个. 精密控制请用"-f"参数
// 如果匹配不到任何清晰度则会尽可能向上匹配
$ bogo -i "https://v.qq.com/x/cover/m5zzglrbt5zdv6d.html" -q 480
Download: [叶问4(普通话版).mp4] 185.31 MB / 574.81 MB [==================================>------------------------------------------------------------------------]  32.24% 23.57 MB/s 00m16s
```

#### 指定文件名
```
$ bogo -i "https://v.qq.com/x/cover/m5zzglrbt5zdv6d.html"  -o test.mp4
Download: [test.mp4] 73.16 MB / 979.09 MB [========>--------------------------------------------------------------------------------------------------------------]   7.47% 23.57 MB/s 00m38s
```