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
	"github.com/missdeer/hannah/handler"
	"github.com/missdeer/hannah/provider"
)

func scanMediasInDirectory(dir string) (res []string) {
	items, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Println(err)
		return nil
	}
	for _, item := range items {
		if item.IsDir() {
			if item.Name() != "." && item.Name() != ".." {
				res = append(res, scanMediasInDirectory(path.Join(dir, item.Name()))...)
			}
		} else {
			if config.SupportedFileType(filepath.Ext(item.Name())) {
				res = append(res, filepath.Join(dir, item.Name()))
			}
		}
	}
	return
}

func scanMedias(medias []string) (res []string) {
	for _, media := range medias {
		if strings.HasPrefix(media, "http://") || strings.HasPrefix(media, "https://") {
			res = append(res, media)
			continue
		}
		fi, err := os.Stat(media)
		if err != nil {
			log.Println(err)
			continue
		}
		if fi.IsDir() {
			res = append(res, scanMediasInDirectory(media)...)
		} else {
			if config.SupportedFileType(filepath.Ext(media)) {
				res = append(res, media)
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
	flag.Parse()

	medias := flag.Args()
	if config.Action == "play" {
		medias = scanMedias(flag.Args())
		if len(medias) == 0 {
			fmt.Println("Please input media URL or local path.")
			return
		}
		fmt.Printf("Found %d songs.\n", len(medias))
	}
	rand.Seed(time.Now().UnixNano())
	switch config.Action {
	case "play":
		for played := false; !played || config.Repeat; played = true {
			if config.Shuffle {
				rand.Shuffle(len(medias), func(i, j int) { medias[i], medias[j] = medias[j], medias[i] })
			}
			for i := 0; i < len(medias); i++ {
				media := medias[i]
				err := handler.PlayMedia(media, i+1, len(medias), "", "") // TODO: extract from file name or ID3v1/v2 tag
				switch err {
				case handler.ShouldQuit:
					return
				case handler.PreviousSong:
					i -= 2
				case handler.NextSong:
					// auto next
				default:
				}
			}
		}
	case "search":
		if config.Provider == "" {
			log.Fatal("set the provider parameter to search")
		}
		p := provider.GetProvider(config.Provider)
		if p != nil {
			songs, err := p.Search(strings.Join(medias, " "), config.Page, config.Limit)
			if err != nil {
				log.Fatal(err)
			}

			for played := false; !played || config.Repeat; played = true {
				if config.Shuffle {
					rand.Shuffle(len(songs), func(i, j int) { songs[i], songs[j] = songs[j], songs[i] })
				}
				for i := 0; i < len(songs); i++ {
					song := songs[i]
					songURL, err := p.SongURL(song)
					if err != nil {
						log.Println(err)
						continue
					}
					err = handler.PlayMedia(songURL, i+1, len(songs), song.Artist, song.Title)
					switch err {
					case handler.ShouldQuit:
						return
					case handler.PreviousSong:
						i -= 2
					case handler.NextSong:
						// auto next
					default:
					}
				}
			}
		}
	}
}
