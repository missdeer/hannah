package rp

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ushis/m3u"

	"github.com/missdeer/hannah/config"
	"github.com/missdeer/hannah/provider"
	"github.com/missdeer/hannah/util"
)

type getterFunc func(provider.IProvider) (provider.Songs, error)
type makeFunc func(*gin.Context, string, string) ([]byte, error)

var (
	playlistPatterns = map[*regexp.Regexp]string{
		regexp.MustCompile(`^https?:\/\/music\.163\.com\/#\/discover\/toplist\?id=(\d+)`):     "netease",
		regexp.MustCompile(`^https?:\/\/music\.163\.com\/#\/playlist\?id=(\d+)`):              "netease",
		regexp.MustCompile(`^https?:\/\/music\.163\.com\/#/my\/m\/music\/playlist\?id=(\d+)`): "netease",
		regexp.MustCompile(`^https?:\/\/www\.xiami\.com\/collect\/(\d+)`):                     "xiami",
		regexp.MustCompile(`^https?:\/\/y\.qq\.com\/n\/yqq\/playlist\/(\d+)\.html`):           "qq",
		regexp.MustCompile(`^https?:\/\/www\.kugou\.com\/yy\/special\/single\/(\d+)\.html`):   "kugou",
		regexp.MustCompile(`^https?:\/\/(www\.)?kuwo\.cn\/playlist_detail\/(\d+)`):            "kuwo",
		regexp.MustCompile(`^https?:\/\/music\.migu\.cn\/v3\/music\/playlist\/(\d+)`):         "migu",
	}
	songPatterns = map[*regexp.Regexp]string{
		regexp.MustCompile(`^https?:\/\/music\.163\.com\/#\/song\?id=(\d+)`):       "netease",
		regexp.MustCompile(`^https?:\/\/www\.xiami\.com\/song\/(\w+)`):             "xiami",
		regexp.MustCompile(`^https?:\/\/y\.qq\.com/n/yqq\/song\/(\w+)\.html`):      "qq",
		regexp.MustCompile(`^https?:\/\/www\.kugou\.com\/song\/#hash=([0-9A-F]+)`): "kugou",
		regexp.MustCompile(`^https?:\/\/(www\.)kuwo.cn\/play_detail\/(\d+)`):       "kuwo",
		regexp.MustCompile(`^https?:\/\/music\.migu\.cn\/v3\/music\/song\/(\d+)`):  "migu",
	}
	artistPatterns = map[*regexp.Regexp]string{
		regexp.MustCompile(`^https?:\/\/music\.163\.com\/weapi\/v1\/artist\/(\d+)`):                                       "netease",
		regexp.MustCompile(`^https?:\/\/music\.163\.com\/#\/artist\?id=(\d+)`):                                            "netease",
		regexp.MustCompile(`^https?:\/\/y\.qq\.com\/n\/yqq\/singer\/(\w+)\.html`):                                         "qq",
		regexp.MustCompile(`^https?:\/\/www\.xiami\.com\/artist\/(\w+)`):                                                  "xiami",
		regexp.MustCompile(`^https?:\/\/www\.xiami\.com\/list\?scene=artist&type=\w+&query={%22artistId%22:%22(\d+)%22}`): "xiami",
		regexp.MustCompile(`^https?:\/\/www\.xiami\.com\/list\?scene=artist&type=\w+&query={"artistId":"(\d+)"}`):         "xiami",
		regexp.MustCompile(`^https?:\/\/(www\.)?kuwo\.cn\/singer_detail\/(\d+)`):                                          "kuwo",
		regexp.MustCompile(`^https?:\/\/music\.migu\.cn\/v3\/music\/artist\/(\d+)`):                                       "migu",
	}
	albumPatterns = map[*regexp.Regexp]string{
		regexp.MustCompile(`^https?:\/\/music\.163\.com\/weapi\/v1\/album\/(\d+)`): "netease",
		regexp.MustCompile(`^https?:\/\/music\.163\.com\/#\/album\?id=(\d+)`):      "netease",
		regexp.MustCompile(`^https?:\/\/y\.qq\.com\/n\/yqq\/album\/(\w+)\.html`):   "qq",
		regexp.MustCompile(`^https?:\/\/www\.xiami\.com\/album\/(\w+)`):            "xiami",
		regexp.MustCompile(`^https?:\/\/(www\.)?kuwo\.cn\/album_detail\/(\d+)`):    "kuwo",
		regexp.MustCompile(`^https?:\/\/music\.migu\.cn\/v3\/music\/album\/(\d+)`): "migu",
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
	baseURL := util.GetBaseURL(c)
	for _, song := range pld {
		filename := strings.Replace(fmt.Sprintf("%s - %s", song.Title, song.Artist), "/", "-", -1)
		playlist = append(playlist, m3u.Track{
			Path:  fmt.Sprintf("%s/%s/%s/%s", baseURL, song.Provider, song.ID, url.PathEscape(filename)),
			Title: song.Title,
		})
	}
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	if _, err = playlist.WriteSimpleTo(w); err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return nil, err
	}

	if err = w.Flush(); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return nil, err
	}

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
	if err = w.Flush(); err != nil {
		return nil, err
	}
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

	baseURL := util.GetBaseURL(c)
	filename := strings.Replace(fmt.Sprintf("%s - %s", song.Title, song.Artist), "/", "-", -1)

	return makeSongInM3U(fmt.Sprintf("%s/%s/%s/%s", baseURL, song.Provider, song.ID, url.PathEscape(filename)), song.Title)
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

	c.Data(http.StatusNotFound, "text/html; charset=UTF-8", notFoundPage)
}

func makeM3U(c *gin.Context, u string, patterns map[*regexp.Regexp]string, make makeFunc) bool {
	for pattern, providerName := range patterns {
		if pattern.MatchString(u) {
			ss := pattern.FindAllStringSubmatch(u, -1)
			if len(ss) == 1 && len(ss[0]) >= 2 {
				b, err := make(c, ss[0][len(ss[0])-1], providerName)
				if err != nil {
					c.Data(http.StatusNotFound, "text/html; charset=UTF-8", notFoundPage)
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
