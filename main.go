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

	"github.com/missdeer/golib/fsutil"
	flag "github.com/spf13/pflag"

	"github.com/jamesnetherton/m3u"

	"github.com/missdeer/hannah/action"
	"github.com/missdeer/hannah/config"
	"github.com/missdeer/hannah/media"
	"github.com/missdeer/hannah/media/decode"
)

var (
	supportedSchema = map[string]struct{}{
		"http://":     {},
		"https://":    {},
		"netease://":  {},
		"qq://":       {},
		"xiami://":    {},
		"bilibili://": {},
		"kugou://":    {},
		"kuwo://":     {},
		"migu://":     {},
	}
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
			if decode.BuiltinSupportedFileType(filepath.Ext(item.Name())) {
				if strings.ToLower(filepath.Ext(item.Name())) == ".m3u" {
					if playlist, err := m3u.Parse(filepath.Join(dir, item.Name())); err == nil {
						for _, track := range playlist.Tracks {
							res = append(res, track.URI)
						}
					}
				} else {
					res = append(res, filepath.Join(dir, item.Name()))
				}
			}
		}
	}
	return
}

func scanSongs(songs []string) (res []string) {
	for _, song := range songs {
		notLocalFile := false
		for k, _ := range supportedSchema {
			if strings.HasPrefix(song, k) {
				notLocalFile = true
				break
			}
		}
		if notLocalFile {
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
			if decode.BuiltinSupportedFileType(filepath.Ext(song)) {
				if strings.ToLower(filepath.Ext(song)) == ".m3u" {
					if playlist, err := m3u.Parse(song); err == nil {
						for _, track := range playlist.Tracks {
							res = append(res, track.URI)
						}
					}
				} else {
					res = append(res, song)
				}
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

	showHelpMessage := false
	flag.IntVarP(&config.Page, "page", "", config.Page, "page number of search result, start from 1")
	flag.IntVarP(&config.Limit, "limit", "l", config.Limit, "max count of search result")
	flag.BoolVarP(&config.Shuffle, "shuffle", "f", config.Shuffle, "shuffle play list order")
	flag.BoolVarP(&config.Repeat, "repeat", "r", config.Repeat, "repeat playing")
	flag.BoolVarP(&config.Mpg123, "mpg123", "m", config.Mpg123, "use mpg123 decoder if it is available")
	flag.StringVarP(&config.AudioDriver, "driver", "d", config.AudioDriver, "set audio deriver, values: "+strings.Join(config.AudioDriverList, ", "))
	flag.StringVarP(&config.Engine, "engine", "e", config.Engine, "specify audio engine, values: builtin, bass, mpv")
	flag.StringVarP(&config.Action, "action", "a", config.Action, "play, search(search and play), m3u(search and save as m3u file), download(search and download media files)")
	flag.StringVarP(&config.Provider, "provider", "p", config.Provider, "netease, xiami, qq, kugou, kuwo, bilibili, migu")
	flag.StringVarP(&config.Socks5Proxy, "socks5", "s", config.Socks5Proxy, "set socks5 proxy, for example: 127.0.0.1:1080")
	flag.StringVarP(&config.HttpProxy, "http-proxy", "t", config.HttpProxy, "set http/https proxy, for example: http://127.0.0.1:1080, https://127.0.0.1:1080 etc.")
	flag.StringVarP(&config.Player, "player", "", config.Player, "specify external player path, use it when the media type is not supported by builtin decoders")
	flag.BoolVarP(&showHelpMessage, "help", "h", false, "show this help message")
	flag.Parse()

	if showHelpMessage {
		flag.PrintDefaults()
		return
	}

	songs := flag.Args()
	if config.Action == "play" {
		songs = scanSongs(flag.Args())
		if len(songs) == 0 {
			log.Fatal("Please input media URL or local path.")
		}
		fmt.Printf("Found %d songs.\n", len(songs))
	}
	rand.Seed(time.Now().UnixNano())

	handler := action.GetActionHandler(config.Action)
	if handler != nil {
		if err := media.Initialize(); err != nil {
			log.Fatal(err)
		}
		defer media.Finalize()
		if err := handler(songs); err != nil {
			log.Println(err)
		}
	} else {
		log.Println("unsupoorted action")
	}
}
