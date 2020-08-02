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

func resolve(song provider.Song) (provider.Song, error) {
	// local filesystem
	if _, err := os.Stat(song.URL); !os.IsNotExist(err) {
		tag, err := id3v2.Open(song.URL, id3v2.Options{Parse: true})
		song.Provider = "local filesystem"
		if err == nil {
			defer tag.Close()
			song.Artist = tag.Artist()
			song.Title = tag.Title()
			return song, nil
		}

		return song, err
	}
	// http/https
	for k, _ := range supportedRemote {
		if strings.HasPrefix(song.URL, k) {
			return song, nil
		}
	}
	// services
	ss := strings.Split(song.URL, "://")
	if len(ss) == 2 {
		schema := ss[0]
		if _, ok := supportedService[schema]; ok {
			p := provider.GetProvider(schema)
			if s, err := p.ResolveSongURL(provider.Song{ID: ss[1]}); err == nil {
				return s, nil
			}
		}
	}
	return song, nil
}

func play(args ...string) error {
	medias := scanSongs(args)
	if len(medias) == 0 {
		return ErrEmptyArgs
	}
	fmt.Printf("Found %d songs.\n", len(medias))

	var songs provider.Songs
	for _, media := range medias {
		songs = append(songs, provider.Song{URL: media})
	}

	for played := false; !played || config.Repeat; played = true {
		if config.Shuffle {
			rand.Shuffle(len(songs), func(i, j int) { songs[i], songs[j] = songs[j], songs[i] })
		}
		if err := playSongs(songs, resolve); err != nil {
			return err
		}
	}
	return nil
}
