## bogo bug灰常多的媒体下载器
### 安装
```
go get github.com/zhxingy/bogo
cd cli
go build -o bogo
```
### 依赖
[FFMPEG](https://www.ffmpeg.org/)
`备注: 非必须, 但会导致解析mp4碎片的网址下载失败`

### 使用
* 用法
`Usage: bogo -i <url> [options...]`

* 导入浏览器cookie
```
bogo config import-cookie <cookie_file>
```
`备注: 仅支持firefox`
* 设置文件存储路径
```
bogo config set-download-path <path>
```

### API
#### [提取器](spider/README.md)
#### [下载器](downloader/README.md)

### 支持网站
* [acfun](https://www.acfun.cn/)
* [acfun番剧](https://www.acfun.cn/v/list155/index.htm)
* [哔哩哔哩](https://www.bilibili.com/)
* [哔哩哔哩番剧](https://www.bilibili.com/anime/)
* [爱奇艺](https://www.iqiyi.com/)
* [优酷](https://www.youku.com/)
* [腾讯视频](https://v.qq.com/)
* [芒果TV](https://www.mgtv.com/tv/)
* [日剧TV](https://www.rijutv.com/)
* [粤视频](http://www.yuesp.com/)
* [西瓜视频](https://www.ixigua.com/)

### 参考项目
* [youtube-dl](https://github.com/ytdl-org/youtube-dl)
* [annie](https://github.com/iawia002/annie)

### 声明
#### 侵权删