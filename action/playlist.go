package action

import (
	"log"
	"math/rand"

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

			if err = shuffleRepeatPlaySongs(songs, p.ResolveSongURL); err != nil {
				return err
			}
		}
	}
	return nil
}
