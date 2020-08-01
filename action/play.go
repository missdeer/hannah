package action

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/bogem/id3v2"
	"github.com/jamesnetherton/m3u"

	"github.com/missdeer/hannah/config"
	"github.com/missdeer/hannah/media"
	"github.com/missdeer/hannah/provider"
	"github.com/missdeer/hannah/util"
)

var (
	ErrEmptyArgs    = errors.New("empty arguments")
	supportedRemote = map[string]struct{}{
		"http://":  {},
		"https://": {},
	}
	supportedService = map[string]struct{}{
		"netease":  {},
		"qq":       {},
		"xiami":    {},
		"bilibili": {},
		"kugou":    {},
		"kuwo":     {},
		"migu":     {},
	}

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
			if media.IsSupportedFileType(filepath.Ext(item.Name())) {
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
			if media.IsSupportedFileType(filepath.Ext(song)) {
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

func resolve(song string) provider.Song {
	// local filesystem
	if _, err := os.Stat(song); !os.IsNotExist(err) {
		tag, err := id3v2.Open(song, id3v2.Options{Parse: true})
		if err == nil {
			defer tag.Close()
			return provider.Song{
				URL:      song,
				Artist:   tag.Artist(),
				Title:    tag.Title(),
				Provider: "local filesystem",
			}
		}

		return provider.Song{
			URL:      song,
			Provider: "local filesystem",
		}
	}
	// http/https
	for k, _ := range supportedRemote {
		if strings.HasPrefix(song, k) {
			return provider.Song{URL: song}
		}
	}
	// services
	ss := strings.Split(song, "://")
	if len(ss) == 2 {
		schema := ss[0]
		if _, ok := supportedService[schema]; ok {
			p := provider.GetProvider(schema)
			if s, err := p.ResolveSongURL(provider.Song{ID: ss[1]}); err == nil {
				return s
			}
		}
	}
	return provider.Song{URL: song}
}

func play(args ...string) error {
	songs := scanSongs(args)
	if len(songs) == 0 {
		return ErrEmptyArgs
	}
	fmt.Printf("Found %d songs.\n", len(songs))

	for played := false; !played || config.Repeat; played = true {
		if config.Shuffle {
			rand.Shuffle(len(songs), func(i, j int) { songs[i], songs[j] = songs[j], songs[i] })
		}
		for i := 0; i < len(songs); i++ {
			song := resolve(songs[i])
			if config.ByExternalPlayer {
				util.ExternalPlay(song.URL)
				continue
			}
			err := media.PlayMedia(song.URL, i+1, len(songs), song.Artist, song.Title) // TODO: extract from file name or ID3v1/v2 tag
			switch err {
			case media.ShouldQuit:
				return err
			case media.PreviousSong:
				i -= 2
			case media.NextSong:
			// auto next
			case media.UnsupportedMediaType:
				log.Println(err, song.URL, ", try to use external player", config.Player)
				if e := util.ExternalPlay(song.URL); e != nil {
					log.Println(err, song.URL)
				}
			default:
				log.Println(err)
			}
		}
	}
	return nil
}
