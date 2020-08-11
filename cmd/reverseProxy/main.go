package main

import (
	"os"
	"path/filepath"

	"github.com/missdeer/golib/fsutil"
	flag "github.com/spf13/pflag"

	"github.com/missdeer/hannah/config"
	"github.com/missdeer/hannah/rp"
)

func main() {
	if homeDir, err := os.UserHomeDir(); err == nil {
		conf := filepath.Join(homeDir, ".hannah.conf")
		if b, e := fsutil.FileExists(conf); e == nil && b {
			config.LoadConfigurationFromFile(conf)
		}
	}
	addr := `127.0.0.1:8321`
	if b := os.Getenv(`BINDADDR`); b != "" {
		addr = b
	}
	showHelpMessage := false
	flag.StringVarP(&addr, "bind-addr", "b", addr, "set bind address")
	flag.StringVarP(&config.Socks5Proxy, "socks5", "s", config.Socks5Proxy, "set socks5 proxy, for example: 127.0.0.1:1080")
	flag.StringVarP(&config.HttpProxy, "http-proxy", "t", config.HttpProxy, "set http/https proxy, for example: http://127.0.0.1:1080, https://127.0.0.1:1080 etc.")
	flag.BoolVarP(&showHelpMessage, "help", "h", false, "show this help message")
	flag.Parse()

	if showHelpMessage {
		flag.PrintDefaults()
		return
	}

	rp.StartReverseProxy(addr)
}
