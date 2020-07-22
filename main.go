package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/missdeer/golib/fsutil"
	flag "github.com/spf13/pflag"

	"github.com/missdeer/hannah/config"
	"github.com/missdeer/hannah/media"
	"github.com/missdeer/hannah/provider"
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

	medias := flag.Args()
	if config.Action == "play" {
		medias = scanSongs(flag.Args())
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
				song := medias[i]
				err := media.PlayMedia(song, i+1, len(medias), "", "") // TODO: extract from file name or ID3v1/v2 tag
				switch err {
				case media.ShouldQuit:
					return
				case media.PreviousSong:
					i -= 2
				case media.NextSong:
				// auto next
				case media.UnsupportedMediaType:
					if b, e := fsutil.FileExists(config.Player); e == nil && b {
						log.Println(err, song, ", try to use external player", config.Player)
						cmd := exec.Command(config.Player, song)
						cmd.Run()
					} else {
						log.Println(err, song)
					}
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
					detail, err := p.SongDetail(song)
					if err != nil {
						log.Println(err)
						continue
					}
					err = media.PlayMedia(detail.URL, i+1, len(songs), song.Artist, song.Title)
					switch err {
					case media.ShouldQuit:
						return
					case media.PreviousSong:
						i -= 2
					case media.NextSong:
						// auto next
					case media.UnsupportedMediaType:
						if b, e := fsutil.FileExists(config.Player); e == nil && b {
							log.Println(err, detail.URL, ", try to use external player", config.Player)
							cmd := exec.Command(config.Player, detail.URL)
							cmd.Run()
						} else {
							log.Println(err, detail.URL)
						}
					default:
					}
				}
			}
		}
	}
}
