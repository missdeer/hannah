package action

import (
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
	return shuffleRepeatPlaySongs(provider.Songs(songs), func(song provider.Song) (provider.Songs, error) {
		s, err := p.ResolveSongURL(song)
		if err != nil {
			return nil, err
		}
		return provider.Songs{s}, err
	})
}

func searchSave(keywords ...string) error {
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
	return saveSongsAsM3U(provider.Songs(songs))
}
