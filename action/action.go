package action

import (
	"errors"
	"log"

	"github.com/missdeer/hannah/config"
	"github.com/missdeer/hannah/media"
	"github.com/missdeer/hannah/provider"
	"github.com/missdeer/hannah/util"
)

type actionHandler func(...string) error
type songResolver func(provider.Song) (provider.Song, error)

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

func playSongs(songs provider.Songs, r songResolver) error {
	for i := 0; i < len(songs); i++ {
		song, err := r(songs[i])
		if err != nil {
			log.Println(err)
			continue
		}
		if song.URL == "" {
			continue
		}
		if config.ByExternalPlayer {
			util.ExternalPlay(song.URL)
			continue
		}
		err = media.PlayMedia(song.URL, i+1, len(songs), song.Artist, song.Title)
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
	return nil
}
