package provider

import (
	"errors"

	"github.com/missdeer/hannah/util"
)

type Song struct {
	ID       string
	Title    string
	URL      string
	Image    string
	Artist   string
	Provider string
	Lyric    string
}

type Songs []Song
type SearchResult []Song

type Playlist struct {
	ID       string
	URL      string
	Image    string
	Title    string
	Provider string
}

type Playlists []Playlist

type IProvider interface {
	Search(keyword string, page int, limit int) (SearchResult, error)
	SongDetail(song Song) (Song, error)
	HotPlaylist(page int) (Playlists, error)
	PlaylistDetail(pl Playlist) (Songs, error)
	Name() string
}

type providerGetter func() IProvider

var (
	ErrStatusNotOK = errors.New("status != 200")

	providerCreatorMap = map[string]providerGetter{
		"netease":   func() IProvider { return &netease{} },
		"xiami":     func() IProvider { return &xiami{client: util.GetHttpClient()} },
		"qq":        func() IProvider { return &qq{} },
		"kugou":     func() IProvider { return &kugou{} },
		"kuwo":      func() IProvider { return &kuwo{} },
		"bilibili":  func() IProvider { return &bilibili{} },
		"migu":      func() IProvider { return &migu{} },
		"musictool": func() IProvider { return &musictool{} },
	}
	providers = make(map[string]IProvider)
)

// GetProvider return the specified provider
func GetProvider(provider string) IProvider {
	if p, ok := providers[provider]; ok {
		return p
	}

	if c, ok := providerCreatorMap[provider]; ok {
		p := c()
		providers[provider] = p
		return p
	}
	return nil
}
