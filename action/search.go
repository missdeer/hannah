package action

import (
	"math/rand"
	"strings"

	"github.com/missdeer/hannah/config"
	"github.com/missdeer/hannah/provider"
)

func search(keywords ...string) error {
	if config.Provider == "" {
		return ErrMissingProvider
	}
	p := provider.GetProvider(config.Provider)
	if p == nil {
		return ErrUnsupportedProvider
	}
	songs, err := p.Search(strings.Join(keywords, " "), config.Page, config.Limit)
	if err != nil {
		return err
	}

	for played := false; !played || config.Repeat; played = true {
		if config.Shuffle {
			rand.Shuffle(len(songs), func(i, j int) { songs[i], songs[j] = songs[j], songs[i] })
		}
		if err = playSongs(provider.Songs(songs), p.ResolveSongURL); err != nil {
			return err
		}
	}
	return nil
}
