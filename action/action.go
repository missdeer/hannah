package action

import (
	"errors"
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
		h               actionHandler
		needScreenPanel bool
	}{
		"play":     {play, true},
		"search":   {search, true},
		"m3u":      {save, false},
		"download": {download, false},
		"hot":      {hot, false},
		"playlist": {playlist, true},
	}
	ErrMissingProvider     = errors.New("set the provider parameter to search")
	ErrUnsupportedProvider = errors.New("unsupported provider")
)

func GetActionHandler(action string) (actionHandler, bool) {
	s, ok := actionHandlerMap[config.Action]
	if !ok {
		return nil, false
	}
	return s.h, s.needScreenPanel
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
	for _, inputSong := range songs {
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
				song.URL = s.URL
			}
			err = media.PlayMedia(song.URL, i+1, len(ss), song.Artist, song.Title)
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
