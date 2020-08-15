package action

import (
	"fmt"
	"log"
	"math/rand"
	"net/url"

	"github.com/missdeer/hannah/config"
	"github.com/missdeer/hannah/provider"
)

func playlist(ids ...string) error {
	if config.Provider == "" {
		return ErrMissingProvider
	}
	p := provider.GetProvider(config.Provider)
	if p == nil {
		return ErrUnsupportedProvider
	}

	for playedPlaylist := false; !playedPlaylist || config.Repeat; playedPlaylist = true {
		if config.Shuffle {
			rand.Shuffle(len(ids), func(i, j int) { ids[i], ids[j] = ids[j], ids[i] })
		}
		for i := 0; i < len(ids); i++ {
			songs, err := p.PlaylistDetail(provider.Playlist{ID: ids[i]})
			if err != nil {
				log.Println(err)
				continue
			}
			err = shuffleRepeatPlaySongs(songs, func(song provider.Song) (provider.Songs, error) {
				s, err := p.ResolveSongURL(song)
				if err != nil {
					return nil, err
				}
				if config.ReverseProxyEnabled {
					scheme := `http`
					host := config.ReverseProxy
					if u, err := url.Parse(config.ReverseProxy); err == nil {
						scheme = u.Scheme
						host = u.Host
					}
					s.URL = fmt.Sprintf("%s://%s/%s/%s", scheme, host, s.Provider, s.ID)
				}
				return provider.Songs{s}, err
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func playlistSave(ids ...string) error {
	if config.Provider == "" {
		return ErrMissingProvider
	}
	p := provider.GetProvider(config.Provider)
	if p == nil {
		return ErrUnsupportedProvider
	}

	for i := 0; i < len(ids); i++ {
		songs, err := p.PlaylistDetail(provider.Playlist{ID: ids[i]})
		if err != nil {
			log.Println(err)
			return err
		}

		if err = saveSongsAsM3U(songs); err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}
