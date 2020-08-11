package provider

import (
	"errors"
	"net/http"
	"sync"

	jsoniter "github.com/json-iterator/go"

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
	ResolveSongURL(song Song) (Song, error)
	ResolveSongLyric(song Song) (Song, error)
	HotPlaylist(page int, limit int) (Playlists, error)
	PlaylistDetail(pl Playlist) (Songs, error)
	Name() string
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
		return i
	}
	return nil
}

var (
	httpClient         *http.Client
	json               = jsoniter.ConfigCompatibleWithStandardLibrary
	ErrStatusNotOK     = errors.New("status != 200")
	providerCreatorMap = map[string]providerGetter{
		"netease":  func() IProvider { return &netease{} },
		"xiami":    func() IProvider { return &xiami{} },
		"qq":       func() IProvider { return &qq{} },
		"kugou":    func() IProvider { return &kugou{} },
		"kuwo":     func() IProvider { return &kuwo{} },
		"bilibili": func() IProvider { return &bilibili{} },
		"migu":     func() IProvider { return &migu{} },
		"ne":       func() IProvider { return &netease{} },
		"xm":       func() IProvider { return &xiami{} },
		"kg":       func() IProvider { return &kugou{} },
		"wu":       func() IProvider { return &kuwo{} },
		"b":        func() IProvider { return &bilibili{} },
		"mg":       func() IProvider { return &migu{} },
		"mt":       func() IProvider { return &musictool{} },
	}
	providers = providerMap{
		m: make(map[string]IProvider),
	}
	once = sync.Once{}
)

// GetProvider return the specified provider
func GetProvider(provider string) IProvider {
	once.Do(func() {
		if httpClient == nil {
			httpClient = util.GetHttpClient()
		}
	})

	if p := providers.get(provider); p != nil {
		return p
	}

	if p := providers.add(provider); p != nil {
		return p
	}
	return nil
}
