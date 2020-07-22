package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

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
	flag.IntVarP(&config.Page, "page", "", 0, "page number of search result")
	flag.IntVarP(&config.Limit, "limit", "", 25, "max count of search result")
	flag.BoolVarP(&config.Shuffle, "shuffle", "", false, "shuffle play list order")
	flag.BoolVarP(&config.Repeat, "repeat", "", false, "repeat playing")
	flag.StringVarP(&config.Action, "action", "a", "play", "play, search(search and play), m3u(search and save as m3u file), download(search and download media files)")
	flag.StringVarP(&config.Provider, "provider", "p", "netease", "netease, xiami, qq, kugou, kuwo, bilibili, migu")
	flag.StringVarP(&config.Socks5Proxy, "socks5", "", "", "set socks5 proxy, for example: 127.0.0.1:1080")
	flag.StringVarP(&config.HttpProxy, "http-proxy", "", "", "set http/https proxy, for example: http://127.0.0.1:1080, https://127.0.0.1:1080 etc.")
	flag.StringVarP(&config.Player, "player", "", "", "specify external player path, use it when the media type is not supported by builtin decoders")
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
