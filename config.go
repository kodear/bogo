package main

import (
	"github.com/larspensjo/config"
	"os"
)

type Config struct {
	fname   string
	root    string
	cookies map[string]string
}

func NewConfig(fname, root string, cookies map[string]string) *Config {
	cfg := &Config{
		fname:   fname,
		root:    root,
		cookies: cookies,
	}
	cfg.init()

	return cfg
}

func (c *Config) init() {
	if _, err := os.Stat(c.fname); err != nil {
		if os.IsNotExist(err) {
			c.Write()
		} else {
			panic(err)
		}
	}
	return
}

func (c *Config) Read() {
	cfg, err := config.ReadDefault(c.fname)
	if err != nil {
		panic(err)
	}

	if cfg.HasSection("bogo") {
		options, err := cfg.SectionOptions("bogo")
		if err != nil {
			panic(err)
		}
		for _, key := range options {
			value, err := cfg.String("bogo", key)
			if err != nil {
				panic(err)
			}
			switch key {
			case "root":
				c.root = value
			}

		}
	}

	if cfg.HasSection("cookies") {
		options, err := cfg.SectionOptions("cookies")
		if err != nil {
			panic(err)
		}
		for _, key := range options {
			value, err := cfg.String("cookies", key)
			if err != nil {
				panic(err)
			}
			c.cookies[key] = value
		}
	}

}

func (c *Config) Write() {
	_ = os.Remove(c.fname)
	cfg := config.New(config.ALTERNATIVE_COMMENT, config.ALTERNATIVE_SEPARATOR, false, false)
	cfg.AddSection("bogo")
	cfg.AddOption("bogo", "root", c.root)
	cfg.AddSection("cookies")
	for k, v := range c.cookies {
		cfg.AddOption("cookies", k, v)
	}

	err := cfg.WriteFile(c.fname, 600, "")
	if err != nil {
		panic(err)
	}

}
