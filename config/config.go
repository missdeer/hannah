package config

import (
	"log"
	"os"

	"gopkg.in/ini.v1"
)

var (
	Shuffle     bool
	Repeat      bool
	Action      = "play"
	Provider    = "netease"
	Socks5Proxy string
	HttpProxy   string
	Limit       = 25
	Page        = 0
	Player      string
	Engine      = "builtin"
	ASIO        bool
	WASAPI      bool
	Mpg123      = true
)

func LoadConfigurationFromFile(fn string) error {
	cfg, err := ini.Load(fn)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	if b, err := cfg.Section("").Key("mpg123").Bool(); err == nil {
		Mpg123 = b
	}
	if b, err := cfg.Section("").Key("shuffle").Bool(); err == nil {
		Shuffle = b
	}
	if b, err := cfg.Section("").Key("repeat").Bool(); err == nil {
		Repeat = b
	}
	if b, err := cfg.Section("").Key("asio").Bool(); err == nil {
		ASIO = b
	}
	if b, err := cfg.Section("").Key("wasapi").Bool(); err == nil {
		WASAPI = b
	}
	if s := cfg.Section("").Key("action").String(); s != "" {
		Action = s
	}
	if s := cfg.Section("").Key("provider").String(); s != "" {
		Provider = s
	}
	if s := cfg.Section("").Key("socks5").String(); s != "" {
		Socks5Proxy = s
	}
	if s := cfg.Section("").Key("http-proxy").String(); s != "" {
		HttpProxy = s
	}
	if s := cfg.Section("").Key("player").String(); s != "" {
		Player = s
	}
	if s := cfg.Section("").Key("engine").String(); s != "" {
		Engine = s
	}
	if i, err := cfg.Section("").Key("limit").Int(); err == nil {
		Limit = i
	}
	if i, err := cfg.Section("").Key("page").Int(); err == nil {
		Page = i
	}

	return nil
}
