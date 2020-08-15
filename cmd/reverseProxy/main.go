package main

import (
	"log"
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
	limit := ""
	flag.IntVarP(&config.ReverseProxyRetries, "retry", "", config.ReverseProxyRetries, "reverse proxy retries count")
	flag.StringVarP(&config.BaseURL, "baseurl", "", config.BaseURL, "set base URL for reverse proxy, used in m3u play list items")
	flag.StringVarP(&addr, "bind-addr", "b", addr, "set bind address")
	flag.StringVarP(&limit, "access-limit", "l", limit, "access limit, CDIR list separated by comma, for example: 172.18.0.0/16, 127.0.0.1/32")
	flag.BoolVarP(&config.CacheEnabled, "cache", "c", config.CacheEnabled, "cache song resolving result in Redis")
	flag.StringVarP(&config.CacheAddr, "cache-addr", "", config.CacheAddr, "set cache(Redis) service address")
	flag.BoolVarP(&config.AutoRedirectURL, "auto-redirect", "", config.AutoRedirectURL, "auto detect origin IP, redirect song URL if origin IP is in China")
	flag.BoolVarP(&config.RedirectURL, "redirect", "", config.RedirectURL, "force all request to be redirected song URL, dont' forward stream by reverse proxy, will overwrite auto-redirect option")
	flag.StringVarP(&config.Socks5Proxy, "socks5", "s", config.Socks5Proxy, "set socks5 proxy, for example: 127.0.0.1:1080")
	flag.StringVarP(&config.HttpProxy, "http-proxy", "t", config.HttpProxy, "set http/https proxy, for example: http://127.0.0.1:1080, https://127.0.0.1:1080 etc.")
	flag.BoolVarP(&showHelpMessage, "help", "h", false, "show this help message")
	flag.Parse()

	if showHelpMessage {
		flag.PrintDefaults()
		return
	}

	config.NetworkTimeout = 0 // no timeout, streaming costs much time
	if err := rp.Init(config.CacheAddr); err != nil {
		log.Println(err)
	}
	log.Fatal(rp.Start(addr, limit))
}
