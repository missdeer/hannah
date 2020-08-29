package provider

import (
	"errors"
	"regexp"
	"sync"

	jsoniter "github.com/json-iterator/go"
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
	SearchSongs(keyword string, page int, limit int) (SearchResult, error)
	ResolveSongURL(song Song) (Song, error)
	ResolveSongLyric(song Song, format string) (Song, error)
	HotPlaylist(page int, limit int) (Playlists, error)
	PlaylistDetail(pl Playlist) (Songs, error)
	ArtistSongs(id string) (Songs, error)
	AlbumSongs(id string) (Songs, error)
	Name() string
	Login() error
}

type providerGetter func() IProvider

type providerMap struct {
	sync.RWMutex
	m map[string]IProvider
}

func (p *providerMap) get(provider string) IProvider {
	p.RLock()
	defer p.RUnlock()
	if res, ok := p.m[provider]; ok {
		return res
	}
	return nil
}

func (p *providerMap) add(provider string) IProvider {
	p.Lock()
	defer p.Unlock()
	if c, ok := providerCreatorMap[provider]; ok {
		i := c()
		p.m[provider] = i
		i.Login()
		return i
	}
	return nil
}

var (
	json               = jsoniter.ConfigCompatibleWithStandardLibrary
	ErrStatusNotOK     = errors.New("status != 200")
	ErrNotImplemented  = errors.New("not implemented yet")
	ErrNoAuthorizeInfo = errors.New("no authorize info")
	providers          = providerMap{m: make(map[string]IProvider)}
	providerCreatorMap = map[string]providerGetter{
		"netease":  func() IProvider { return &netease{} },
		"xiami":    func() IProvider { return &xiami{} },
		"qq":       func() IProvider { return &qq{} },
		"kugou":    func() IProvider { return &kugou{} },
		"kuwo":     func() IProvider { return &kuwo{} },
		"bilibili": func() IProvider { return &bilibili{} },
		"migu":     func() IProvider { return &migu{} },
	}
	providerIDPatternMap = map[string]*regexp.Regexp{
		"netease":  regexp.MustCompile(`^[0-9]+$`),
		"xiami":    regexp.MustCompile(`^[0-9]+$`),
		"qq":       regexp.MustCompile(`^[0-9a-zA-Z]+$`),
		"kugou":    regexp.MustCompile(`^[0-9A-F]+$`),
		"kuwo":     regexp.MustCompile(`^[0-9]+$`),
		"bilibili": regexp.MustCompile(`^[0-9]+$`),
		"migu":     regexp.MustCompile(`^[0-9a-zA-Z]+$`),
	}
)

// GetProvider return the specified provider
func GetProvider(provider string) IProvider {
	if p := providers.get(provider); p != nil {
		return p
	}

	if p := providers.add(provider); p != nil {
		return p
	}
	return nil
}

func GetSongIDPattern(provider string) *regexp.Regexp {
	r, ok := providerIDPatternMap[provider]
	if ok {
		return r
	}
	return nil
}
