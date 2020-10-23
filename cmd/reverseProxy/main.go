package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/missdeer/golib/fsutil"
	flag "github.com/spf13/pflag"

	"github.com/missdeer/hannah/config"
	"github.com/missdeer/hannah/rp"
)

var (
	// Gitcommit contains the commit where we built reverseProxy from.
	GitCommit string
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
	showVersion := false
	limit := ""
	flag.IntVarP(&config.FixedStreamCacheSize, "fixed-stream-cache-size", "", config.FixedStreamCacheSize, "set fixed stream cache size, if it's 0, then the cache size is decided by HTTP header Content-Length")
	flag.StringVarP(&config.NetworkInterface, "network-interface", "i", config.NetworkInterface, "set local network interface name, for example: en1, will overwirte socks5/http-proxy option")
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
	flag.BoolVarP(&showVersion, "version", "v", false, "show version number")
	flag.Parse()

	if showHelpMessage {
		flag.PrintDefaults()
		return
	}

	if showVersion {
		fmt.Println("Hannah Reverse Proxy version:", GitCommit)
		return
	}

	rand.Seed(time.Now().UnixNano())

	config.NetworkTimeout = 0 // no timeout, streaming costs much time
	if err := rp.Init(config.CacheAddr); err != nil {
		log.Println(err)
	}
	rp.StartDaemon(addr, limit)
}
