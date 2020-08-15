package action

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/url"

	"github.com/missdeer/hannah/config"
	"github.com/missdeer/hannah/media"
	"github.com/missdeer/hannah/provider"
	"github.com/missdeer/hannah/util"
)

type actionHandler func(...string) error
type songResolver func(provider.Song) (provider.Songs, error)

var (
	actionHandlerMap = map[string]struct {
		h      actionHandler
		holdOn bool // false - will exit application in a short time, true - player will play songs, application is hold on
	}{
		"play":          {play, true},
		"search":        {search, true},
		"search-save":   {searchSave, false},
		"hot":           {hot, false},
		"playlist":      {playlist, true},
		"playlist-save": {playlistSave, false},
	}
	ErrMissingProvider     = errors.New("set the provider parameter to search")
	ErrUnsupportedProvider = errors.New("unsupported provider")
)

func GetActionHandler(action string) (h actionHandler, holdOn bool) {
	s, ok := actionHandlerMap[config.Action]
	if !ok {
		return nil, false
	}
	return s.h, s.holdOn
}

func shuffleRepeatPlaySongs(songs provider.Songs, r songResolver) error {
	for played := false; !played || config.Repeat; played = true {
		if config.Shuffle {
			rand.Shuffle(len(songs), func(i, j int) { songs[i], songs[j] = songs[j], songs[i] })
		}
		if err := playSongs(songs, r); err != nil {
			return err
		}
	}
	return nil
}

func playSongs(songs provider.Songs, r songResolver) error {
	for j, inputSong := range songs {
		ss, err := r(inputSong)
		if err != nil {
			log.Println(err)
			continue
		}
		if len(ss) == 0 {
			continue
		}
		if config.Shuffle {
			rand.Shuffle(len(ss), func(i, j int) { ss[i], ss[j] = ss[j], ss[i] })
		}
		u, err := url.Parse(inputSong.URL)
		var p provider.IProvider
		if err == nil {
			p = provider.GetProvider(u.Scheme)
		}
		index := j + 1
		count := len(songs)
		for i, song := range ss {
			if config.ByExternalPlayer {
				util.ExternalPlay(song.URL)
				continue
			}
			if song.URL == "" && p != nil {
				// from playlist, only song ID exists, get the song URL now
				s, err := p.ResolveSongURL(song)
				if err != nil {
					log.Println(err)
					continue
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
				song.URL = s.URL
				count = len(ss)
				index = i + 1
			}
			if song.URL == "" {
				continue
			}
			err = media.PlayMedia(song, index, count)
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
			}
		}
	}
	return nil
}

func saveSongsAsM3U(songs provider.Songs) error {
	return media.AppendSongsToM3U(songs, true)
}
