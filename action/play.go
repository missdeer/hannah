package action

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/bogem/id3v2"
	"github.com/ushis/m3u"

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
					fn := filepath.Join(dir, item.Name())
					if f, err := os.OpenFile(fn, os.O_RDONLY, 0644); err == nil {
						if playlist, err := m3u.Parse(f); err == nil {
							for _, track := range playlist {
								res = append(res, track.Path)
							}
						}
						f.Close()
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
					if f, err := os.OpenFile(song, os.O_RDONLY, 0644); err == nil {
						if playlist, err := m3u.Parse(f); err == nil {
							for _, track := range playlist {
								res = append(res, track.Path)
							}
						}
						f.Close()
					}
				} else {
					res = append(res, song)
				}
			}
		}
	}
	return
}

func resolve(song provider.Song) (provider.Songs, error) {
	// local filesystem
	if stat, err := os.Stat(song.URL); os.IsExist(err) && !stat.IsDir() {
		tag, err := id3v2.Open(song.URL, id3v2.Options{Parse: true})
		song.Provider = "local filesystem"
		if err == nil {
			defer tag.Close()
			song.Artist = tag.Artist()
			song.Title = tag.Title()
			return provider.Songs{song}, nil
		}

		return provider.Songs{song}, err
	}
	// http/https
	for k, _ := range supportedRemote {
		if strings.HasPrefix(song.URL, k) {
			song.Provider = "http(s)"
			return provider.Songs{song}, nil
		}
	}
	// services
	u, err := url.Parse(song.URL)
	if err != nil {
		return provider.Songs{song}, err
	}
	values, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return provider.Songs{song}, err
	}

	linkType := `song`
	if t, ok := values[`type`]; ok && len(t) > 0 {
		linkType = t[0]
	}
	providerName := u.Scheme
	if _, ok := supportedService[providerName]; ok {
		p := provider.GetProvider(providerName)
		switch linkType {
		case "song":
			if config.ReverseProxyEnabled {
				scheme := `http`
				host := config.ReverseProxy
				if u, err := url.Parse(config.ReverseProxy); err == nil {
					scheme = u.Scheme
					host = u.Host
				}
				return provider.Songs{
					provider.Song{
						ID:       u.Host,
						Provider: providerName,
						URL:      fmt.Sprintf("%s://%s/%s/%s", scheme, host, providerName, u.Host),
					},
				}, nil
			}
			// TODO extract song title & artist
			if s, err := p.ResolveSongURL(provider.Song{ID: u.Host}); err == nil {
				return provider.Songs{s}, nil
			}
		case "playlist":
			if songs, err := p.PlaylistDetail(provider.Playlist{ID: u.Host}); err == nil {
				return songs, nil
			}
		default:
		}
	}

	return provider.Songs{song}, nil
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
	return shuffleRepeatPlaySongs(songs, resolve)
}
