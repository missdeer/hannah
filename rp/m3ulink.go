package rp

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
	"github.com/ushis/m3u"

	"github.com/missdeer/hannah/config"
	"github.com/missdeer/hannah/provider"
)

type getterFunc func(provider.IProvider) (provider.Songs, error)
type makeFunc func(*gin.Context, string, string) ([]byte, error)

var (
	playlistPatterns = map[*regexp.Regexp]string{
		regexp.MustCompile(`^https?:\/\/music.163.com\/#\/discover\/toplist\?id=([0-9]+)`):     "netease",
		regexp.MustCompile(`^https?:\/\/music.163.com\/#\/playlist\?id=([0-9]+)`):              "netease",
		regexp.MustCompile(`^https?:\/\/music.163.com\/#/my\/m\/music\/playlist\?id=([0-9]+)`): "netease",
		regexp.MustCompile(`^https?:\/\/www.xiami.com\/collect\/([0-9]+)`):                     "xiami",
		regexp.MustCompile(`^https?:\/\/y.qq.com\/n\/yqq\/playlist\/([0-9]+)\.html`):           "qq",
		regexp.MustCompile(`^https?:\/\/www.kugou.com\/yy\/special\/single\/([0-9]+)\.html`):   "kugou",
		regexp.MustCompile(`^https?:\/\/kuwo.cn\/playlist_detail\/([0-9]+)`):                   "kuwo",
		regexp.MustCompile(`^https?:\/\/music.migu.cn\/v3\/music\/playlist\/([0-9]+)`):         "migu",
	}
	songPatterns = map[*regexp.Regexp]string{
		regexp.MustCompile(`^https?:\/\/music.163.com\/#\/song\?id=([0-9]+)`):        "netease",
		regexp.MustCompile(`^https?:\/\/www.xiami.com\/song\/([0-9a-zA-Z]+)`):        "xiami",
		regexp.MustCompile(`^https?:\/\/y.qq.com/n/yqq\/song\/([0-9a-zA-Z]+)\.html`): "qq",
		regexp.MustCompile(`^https?:\/\/www.kugou.com\/song\/#hash=([0-9A-F]+)`):     "kugou",
		regexp.MustCompile(`^https?:\/\/kuwo.cn\/play_detail\/([0-9]+)`):             "kuwo",
		regexp.MustCompile(`^https?:\/\/music.migu.cn\/v3\/music\/song\/([0-9]+)`):   "migu",
	}
	artistPatterns = map[*regexp.Regexp]string{
		regexp.MustCompile(`^https?:\/\/music.163.com\/weapi\/v1\/artist\/([0-9]+)`):     "netease",
		regexp.MustCompile(`^https?:\/\/y.qq.com\/n\/yqq\/singer\/([0-9a-zA-Z]+)\.html`): "qq",
		regexp.MustCompile(`^https?:\/\/www.xiami.com\/artist\/([0-9a-zA-Z]+)`):          "xiami",
	}
	albumPatterns = map[*regexp.Regexp]string{
		regexp.MustCompile(`^https?:\/\/music.163.com\/weapi\/v1\/album\/([0-9]+)`):     "netease",
		regexp.MustCompile(`^https?:\/\/y.qq.com\/n\/yqq\/album\/([0-9a-zA-Z]+)\.html`): "qq",
		regexp.MustCompile(`^https?:\/\/www.xiami.com\/album\/([0-9a-zA-Z]+)`):          "xiami",
	}
)

func makeSongs(c *gin.Context, providerName string, cacheKey string, getter getterFunc) ([]byte, error) {
	if config.CacheEnabled {
		b, err := redis.GetBytes(cacheKey)
		if err == nil {
			return b, nil
		}
	}

	// resolve playlist
	p := provider.GetProvider(providerName)
	if p == nil {
		return nil, errUnsupportedProvider
	}

	pld, err := getter(p)
	if err != nil {
		return nil, err
	}

	playlist := m3u.Playlist{}
	baseURL := config.BaseURL
	if baseURL == "" {
		scheme := c.Request.Header.Get("X-Forwarded-Proto")
		if scheme == "" {
			originURL := location.Get(c)
			scheme = originURL.Scheme
		}
		baseURL = fmt.Sprintf("%s://%s", scheme, c.Request.Host)
	}
	for _, song := range pld {
		filename := strings.Replace(fmt.Sprintf("%s - %s", song.Title, song.Artist), "/", "-", -1)
		playlist = append(playlist, m3u.Track{
			Path:  fmt.Sprintf("%s/%s/%s/%s", baseURL, song.Provider, song.ID, url.PathEscape(filename)),
			Title: song.Title,
		})
	}
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	_, err = playlist.WriteSimpleTo(w)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return nil, err
	}
	w.Flush()

	if config.CacheEnabled {
		redis.Put(cacheKey, b.Bytes())
	}

	return b.Bytes(), nil
}

func makeArtistSongs(c *gin.Context, id string, providerName string) ([]byte, error) {
	urlKey := fmt.Sprintf("%s:%s:artist", providerName, id)
	return makeSongs(c, providerName, urlKey, func(p provider.IProvider) (provider.Songs, error) {
		pld, err := p.ArtistSongs(id)
		if err != nil {
			return nil, err
		}
		return pld, nil
	})
}

func makeAlbumSongs(c *gin.Context, id string, providerName string) ([]byte, error) {
	urlKey := fmt.Sprintf("%s:%s:album", providerName, id)
	return makeSongs(c, providerName, urlKey, func(p provider.IProvider) (provider.Songs, error) {
		pld, err := p.AlbumSongs(id)
		if err != nil {
			return nil, err
		}
		return pld, nil
	})
}

func makePlaylist(c *gin.Context, id string, providerName string) ([]byte, error) {
	urlKey := fmt.Sprintf("%s:%s:playlist", providerName, id)
	return makeSongs(c, providerName, urlKey, func(p provider.IProvider) (provider.Songs, error) {
		pld, err := p.PlaylistDetail(provider.Playlist{ID: id})
		if err != nil {
			return nil, err
		}
		return pld, nil
	})
}

func makeSongInM3U(songURL string, songTitle string) ([]byte, error) {
	playlist := m3u.Playlist{m3u.Track{
		Path:  songURL,
		Title: songTitle,
	}}
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	_, err := playlist.WriteSimpleTo(w)
	if err != nil {
		return nil, err
	}
	w.Flush()
	return b.Bytes(), nil
}

func makeSong(c *gin.Context, id string, providerName string) ([]byte, error) {
	// check cache first
	urlKey := fmt.Sprintf("%s:%s:url", providerName, id)
	if config.CacheEnabled {
		if songURL, err := redis.GetString(urlKey); err == nil {
			return makeSongInM3U(songURL, "")
		}
	}

	// resolve URL now
	p := provider.GetProvider(providerName)
	if p == nil {
		return nil, errUnsupportedProvider
	}
	song, err := p.ResolveSongURL(provider.Song{ID: id})
	if err != nil {
		return nil, err
	}

	// cache the resolved result
	if redis != nil {
		redis.PutWithTimeout(urlKey, song.URL, cacheTimeout)
	}

	return makeSongInM3U(song.URL, song.Title)
}

func generateM3ULink(c *gin.Context) {
	u := c.Query("u")
	if u == "" {
		c.AbortWithError(http.StatusNotFound, errInvalidURL)
		return
	}
	if makeM3U(c, u, playlistPatterns, makePlaylist) ||
		makeM3U(c, u, albumPatterns, makeAlbumSongs) ||
		makeM3U(c, u, artistPatterns, makeArtistSongs) ||
		makeM3U(c, u, songPatterns, makeSong) {
		return
	}

	c.Data(http.StatusNotFound, "text/html; charset=UTF-8",
		[]byte(`<html><script type="text/javascript" src="//qzonestyle.gtimg.cn/qzone/hybrid/app/404/search_children.js" charset="utf-8"></script><body></body></html>`))
}

func makeM3U(c *gin.Context, u string, patterns map[*regexp.Regexp]string, make makeFunc) bool {
	for pattern, providerName := range patterns {
		if pattern.MatchString(u) {
			ss := pattern.FindAllStringSubmatch(u, -1)
			if len(ss) == 1 && len(ss[0]) == 2 {
				b, err := make(c, ss[0][1], providerName)
				if err != nil {
					c.Data(http.StatusNotFound, "text/html; charset=UTF-8",
						[]byte(`<html><script type="text/javascript" src="//qzonestyle.gtimg.cn/qzone/hybrid/app/404/search_children.js" charset="utf-8"></script><body></body></html>`))
				} else {
					c.Writer.Header().Set(`Content-Disposition`, `attachment; filename="playlist.m3u"`)
					c.Data(http.StatusOK, "audio/x-mpegurl", b)
				}
				return true
			}
		}
	}
	return false
}
