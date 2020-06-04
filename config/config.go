package config

import (
	"encoding/json"
	"fmt"
	"github.com/zhxingy/bogo/cookie"
	"io/ioutil"
	"os"
	user2 "os/user"
	"path/filepath"
)

type config struct {
	DownloadPath string            `json:"download_path"`
	Proxy        string            `json:"proxy"`
	Cookies      cookie.CookiesJar `json:"cookies"`
}

type MainConfig struct {
	Config config
	File   *os.File
}

func Open(file string) *MainConfig {
	if file == "" {
		file = DefaultConfig()
	}
	f, err := os.OpenFile(file, os.O_RDWR, 644)
	if err != nil {
		panic(fmt.Sprintf("open config file failure. err msg: %v\n", err))
	}

	body, err := ioutil.ReadAll(f)
	if err != nil {
		panic(fmt.Sprintf("read config file failure. err msg: %v\n", err))
	}

	var cfg config
	err = json.Unmarshal(body, &cfg)
	if err != nil {
		panic(fmt.Sprintf("load config file failure. err msg: %v\n", err))
	}

	return &MainConfig{
		Config: cfg,
		File:   f,
	}
}

func (cfg *MainConfig) Close() {
	_ = cfg.File.Close()
}

func (cfg *MainConfig) Write() {
	body, err := json.MarshalIndent(cfg.Config, "", "\t")

	if err != nil {
		panic(fmt.Sprintf("parse config json failure. err msg: %v\n", err))
	}

	err = cfg.File.Truncate(0)
	_, err = cfg.File.Seek(0, 0)
	_, err = cfg.File.Write(body)
	if err != nil {
		panic(fmt.Sprintf("write config file failure. err msg: %v\n", err))
	}
}

func DefaultConfig() string {
	user, err := user2.Current()
	if err != nil {
		panic(fmt.Sprintf("get config path failure. err msg: %v\n", err))
	}

	return filepath.Join(user.HomeDir, ".config", "bogo.json")
}

func DefaultDownloadPath() string {
	user, err := user2.Current()
	if err != nil {
		panic(fmt.Sprintf("get download path failure. err msg: %v\n", err))
	}

	path := filepath.Join(user.HomeDir, "BogoDownload")
	_ = os.MkdirAll(path, 0666)

	return path
}

func init() {
	configFile := DefaultConfig()
	configDir, _ := filepath.Split(configFile)

	_, err := os.Stat(configDir)
	if err != nil {
		err = os.MkdirAll(configDir, os.ModePerm)
		if err != nil {
			panic(fmt.Sprintf("create config dir failure. err msg: %v\n", err))
		}
	}

	_, err = os.Stat(configFile)
	if err != nil {
		if !os.IsExist(err) {
			f, err := os.Create(configFile)
			if err != nil {
				panic(fmt.Sprintf("create config file failure. err msg: %v\n", err))
			}
			cfg := &MainConfig{
				Config: config{
					DownloadPath: DefaultDownloadPath(),
				},
				File: f,
			}
			cfg.Write()
			cfg.Close()
		}
	}
}
