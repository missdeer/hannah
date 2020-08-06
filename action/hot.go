package action

import (
	"fmt"

	"github.com/missdeer/hannah/config"
	"github.com/missdeer/hannah/provider"
)

func hot(args ...string) error {
	if config.Provider == "" {
		return ErrMissingProvider
	}
	p := provider.GetProvider(config.Provider)
	if p == nil {
		return ErrUnsupportedProvider
	}

	pls, err := p.HotPlaylist(config.Page,config.Limit)
	if err != nil {
		return err
	}

	for i, pl := range pls {
		fmt.Printf("[%d] %s - %s\n", i, pl.Title, pl.ID)
	}
	return nil
}
