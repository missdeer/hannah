package main

import (
	"errors"
	"log"
	"math/rand"
	"os/exec"
	"strings"

	"github.com/missdeer/golib/fsutil"

	"github.com/missdeer/hannah/config"
	"github.com/missdeer/hannah/media"
	"github.com/missdeer/hannah/provider"
)

type actionHandler func([]string) error

var (
	actionHandlerMap = map[string]actionHandler{
		"play":   actionPlay,
		"search": actionSearch,
	}
)

func actionPlay(songs []string) error {
	for played := false; !played || config.Repeat; played = true {
		if config.Shuffle {
			rand.Shuffle(len(songs), func(i, j int) { songs[i], songs[j] = songs[j], songs[i] })
		}
		for i := 0; i < len(songs); i++ {
			song := songs[i]
			err := media.PlayMedia(song, i+1, len(songs), "", "") // TODO: extract from file name or ID3v1/v2 tag
			switch err {
			case media.ShouldQuit:
				return err
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
					log.Println(e, song)
				}
			default:
				log.Println(err)
			}
		}
	}
	return nil
}

func actionSearch(keywords []string) error {
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
			song := songs[i]
			detail, err := p.SongDetail(song)
			if err != nil {
				log.Println(err)
				continue
			}
			err = media.PlayMedia(detail.URL, i+1, len(songs), song.Artist, song.Title)
			switch err {
			case media.ShouldQuit:
				return err
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
	return nil
}
