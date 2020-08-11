package main

import (
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/missdeer/golib/fsutil"
	flag "github.com/spf13/pflag"

	"github.com/missdeer/hannah/action"
	"github.com/missdeer/hannah/config"
	"github.com/missdeer/hannah/media"
	"github.com/missdeer/hannah/rp"
)

func main() {
	if homeDir, err := os.UserHomeDir(); err == nil {
		conf := filepath.Join(homeDir, ".hannah.conf")
		if b, e := fsutil.FileExists(conf); e == nil && b {
			config.LoadConfigurationFromFile(conf)
		}
	}

	showHelpMessage := false
	flag.BoolVarP(&config.ReverseProxyEnabled, "reverse-proxy-enabled", "", config.ReverseProxyEnabled, "reverse proxy enabled")
	flag.StringVarP(&config.ReverseProxy, "reverse-proxy", "", config.ReverseProxy, "set reverse proxy address")
	flag.IntVarP(&config.Page, "page", "", config.Page, "page number of search result, start from 1")
	flag.IntVarP(&config.Limit, "limit", "l", config.Limit, "max count of search result")
	flag.BoolVarP(&config.Shuffle, "shuffle", "f", config.Shuffle, "shuffle play list order")
	flag.BoolVarP(&config.Repeat, "repeat", "r", config.Repeat, "repeat playing")
	flag.StringVarP(&config.AudioDriver, "driver", "d", config.AudioDriver, "set audio deriver, values: "+strings.Join(config.AudioDriverList, ", "))
	flag.StringVarP(&config.Engine, "engine", "e", config.Engine, "specify audio engine, values: builtin, bass")
	flag.StringVarP(&config.Action, "action", "a", config.Action, "play(play songs in file/playlist), search(search songs and play), search-save(search songs and append to m3u file), hot(get hot playlists), playlist(play songs in the specified playlist), playlist-save(parse playlist and append to m3u file)")
	flag.StringVarP(&config.Provider, "provider", "p", config.Provider, "netease, xiami, qq, kugou, kuwo, bilibili, migu")
	flag.StringVarP(&config.Socks5Proxy, "socks5", "s", config.Socks5Proxy, "set socks5 proxy, for example: 127.0.0.1:1080")
	flag.StringVarP(&config.HttpProxy, "http-proxy", "t", config.HttpProxy, "set http/https proxy, for example: http://127.0.0.1:1080, https://127.0.0.1:1080 etc.")
	flag.StringVarP(&config.Player, "player", "", config.Player, "specify external player path, use it when the media type is not supported by builtin decoders")
	flag.BoolVarP(&config.ByExternalPlayer, "by-external-player", "y", config.ByExternalPlayer, "play by external player")
	flag.StringVarP(&config.DownloadDir, "dir", "", config.DownloadDir, "set directory to save download files")
	flag.StringVarP(&config.M3UFileName, "m3u", "", config.M3UFileName, "set m3u file name to save play list")
	flag.BoolVarP(&showHelpMessage, "help", "h", false, "show this help message")
	flag.Parse()

	if showHelpMessage {
		flag.PrintDefaults()
		return
	}

	if config.ReverseProxyEnabled {
		config.NetworkTimeout = 0 // no timeout, streaming costs much time
		go rp.StartReverseProxy(config.ReverseProxy)
	}

	args := flag.Args()
	rand.Seed(time.Now().UnixNano())

	handler, needScreenPanel := action.GetActionHandler(config.Action)
	if handler != nil {
		if err := media.Initialize(!config.ByExternalPlayer && needScreenPanel); err != nil {
			log.Fatal(err)
		}
		defer media.Finalize(!config.ByExternalPlayer && needScreenPanel)
		if err := handler(args...); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal("unsupoorted action")
	}
}
