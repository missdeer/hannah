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
)

func main() {
	if homeDir, err := os.UserHomeDir(); err == nil {
		conf := filepath.Join(homeDir, ".hannah.conf")
		if b, e := fsutil.FileExists(conf); e == nil && b {
			config.LoadConfigurationFromFile(conf)
		}
	}

	showHelpMessage := false
	flag.IntVarP(&config.Page, "page", "", config.Page, "page number of search result, start from 1")
	flag.IntVarP(&config.Limit, "limit", "l", config.Limit, "max count of search result")
	flag.BoolVarP(&config.Shuffle, "shuffle", "f", config.Shuffle, "shuffle play list order")
	flag.BoolVarP(&config.Repeat, "repeat", "r", config.Repeat, "repeat playing")
	flag.BoolVarP(&config.Mpg123, "mpg123", "m", config.Mpg123, "use mpg123 decoder if it is available")
	flag.StringVarP(&config.AudioDriver, "driver", "d", config.AudioDriver, "set audio deriver, values: "+strings.Join(config.AudioDriverList, ", "))
	flag.StringVarP(&config.Engine, "engine", "e", config.Engine, "specify audio engine, values: builtin, bass")
	flag.StringVarP(&config.Action, "action", "a", config.Action, "play, search(search and play), m3u(search and save as m3u file), download(search and download media files), hot(get hot playlists), playlist(play args in the specified playlist)")
	flag.StringVarP(&config.Provider, "provider", "p", config.Provider, "netease, xiami, qq, kugou, kuwo, bilibili, migu")
	flag.StringVarP(&config.Socks5Proxy, "socks5", "s", config.Socks5Proxy, "set socks5 proxy, for example: 127.0.0.1:1080")
	flag.StringVarP(&config.HttpProxy, "http-proxy", "t", config.HttpProxy, "set http/https proxy, for example: http://127.0.0.1:1080, https://127.0.0.1:1080 etc.")
	flag.StringVarP(&config.Player, "player", "", config.Player, "specify external player path, use it when the media type is not supported by builtin decoders")
	flag.BoolVarP(&config.ByExternalPlayer, "by-external-player", "y", config.ByExternalPlayer, "play by external player")
	flag.BoolVarP(&showHelpMessage, "help", "h", false, "show this help message")
	flag.Parse()

	if showHelpMessage {
		flag.PrintDefaults()
		return
	}

	args := flag.Args()
	rand.Seed(time.Now().UnixNano())

	handler := action.GetActionHandler(config.Action)
	if handler != nil {
		if err := media.Initialize(!config.ByExternalPlayer); err != nil {
			log.Fatal(err)
		}
		defer media.Finalize(!config.ByExternalPlayer)
		if err := handler(args...); err != nil {
			log.Println(err)
		}
	} else {
		log.Println("unsupoorted action")
	}
}
