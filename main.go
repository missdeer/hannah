package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/missdeer/golib/fsutil"
	flag "github.com/spf13/pflag"

	"github.com/missdeer/hannah/config"
	"github.com/missdeer/hannah/media"
)

func scanSongsInDirectory(dir string) (res []string) {
	items, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Println(err)
		return nil
	}
	for _, item := range items {
		if item.IsDir() {
			if item.Name() != "." && item.Name() != ".." {
				res = append(res, scanSongsInDirectory(path.Join(dir, item.Name()))...)
			}
		} else {
			if media.BuiltinSupportedFileType(filepath.Ext(item.Name())) {
				res = append(res, filepath.Join(dir, item.Name()))
			}
		}
	}
	return
}

func scanSongs(songs []string) (res []string) {
	for _, song := range songs {
		if strings.HasPrefix(song, "http://") || strings.HasPrefix(song, "https://") {
			res = append(res, song)
			continue
		}
		fi, err := os.Stat(song)
		if err != nil {
			log.Println(err)
			continue
		}
		if fi.IsDir() {
			res = append(res, scanSongsInDirectory(song)...)
		} else {
			if media.BuiltinSupportedFileType(filepath.Ext(song)) {
				res = append(res, song)
			}
		}
	}
	return
}

func main() {
	if homeDir, err := os.UserHomeDir(); err == nil {
		conf := filepath.Join(homeDir, ".hannah.conf")
		if b, e := fsutil.FileExists(conf); e == nil && b {
			config.LoadConfigurationFromFile(conf)
		}
	}

	flag.IntVarP(&config.Page, "page", "", config.Page, "page number of search result")
	flag.IntVarP(&config.Limit, "limit", "l", config.Limit, "max count of search result")
	flag.BoolVarP(&config.Shuffle, "shuffle", "s", config.Shuffle, "shuffle play list order")
	flag.BoolVarP(&config.Repeat, "repeat", "r", config.Repeat, "repeat playing")
	flag.BoolVarP(&config.Mpg123, "mpg123", "m", config.Mpg123, "use mpg123 decoder if it is available")
	if runtime.GOOS == "windows" {
		flag.BoolVarP(&config.ASIO, "asio", "", config.ASIO, "use ASIO driver, Windows only")
		flag.BoolVarP(&config.WASAPI, "wasapi", "w", config.WASAPI, "use WASAPI driver, Windows only")
	}
	flag.StringVarP(&config.Engine, "engine", "e", config.Engine, "specify audio engine, values: builtin, bass, mpv")
	flag.StringVarP(&config.Action, "action", "a", config.Action, "play, search(search and play), m3u(search and save as m3u file), download(search and download media files)")
	flag.StringVarP(&config.Provider, "provider", "p", config.Provider, "netease, xiami, qq, kugou, kuwo, bilibili, migu")
	flag.StringVarP(&config.Socks5Proxy, "socks5", "", config.Socks5Proxy, "set socks5 proxy, for example: 127.0.0.1:1080")
	flag.StringVarP(&config.HttpProxy, "http-proxy", "", config.HttpProxy, "set http/https proxy, for example: http://127.0.0.1:1080, https://127.0.0.1:1080 etc.")
	flag.StringVarP(&config.Player, "player", "", config.Player, "specify external player path, use it when the media type is not supported by builtin decoders")
	flag.Parse()

	songs := flag.Args()
	if config.Action == "play" {
		songs = scanSongs(flag.Args())
		if len(songs) == 0 {
			log.Fatal("Please input media URL or local path.")
		}
		fmt.Printf("Found %d songs.\n", len(songs))
	}
	rand.Seed(time.Now().UnixNano())
	defer media.Finalize()

	handler, ok := actionHandlerMap[config.Action]
	if ok {
		handler(songs)
	} else {
		log.Fatal("unsupoorted action")
	}
}
