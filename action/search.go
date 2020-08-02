package action

import (
	"errors"
	"log"
	"math/rand"
	"strings"

	"github.com/missdeer/hannah/config"
	"github.com/missdeer/hannah/media"
	"github.com/missdeer/hannah/provider"
	"github.com/missdeer/hannah/util"
)

func search(keywords ...string) error {
	if config.Provider == "" {
		return errors.New("set the provider parameter to search")
	}
	p := provider.GetProvider(config.Provider)
	if p == nil {
		return errors.New("unsupported provider")
	}
	songs, err := p.Search(strings.Join(keywords, " "), config.Page, config.Limit)
	if err != nil {
		return err
	}

	for played := false; !played || config.Repeat; played = true {
		if config.Shuffle {
			rand.Shuffle(len(songs), func(i, j int) { songs[i], songs[j] = songs[j], songs[i] })
		}
		for i := 0; i < len(songs); i++ {
			song, err := p.ResolveSongURL(songs[i])
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
	}
	return nil
}
