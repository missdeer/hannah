package rp

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ushis/m3u"

	"github.com/missdeer/hannah/config"
	"github.com/missdeer/hannah/provider"
	"github.com/missdeer/hannah/util"
)

type getterFunc func(provider.IProvider) (provider.Songs, error)
type makeFunc func(*gin.Context, string, string) ([]byte, error)
type patternMatchFunc func(string) (string, string, bool)

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

func makeSingleSongInM3U(songURL string, songTitle string) ([]byte, error) {
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

func makeSingleSong(c *gin.Context, id string, providerName string) ([]byte, error) {
	// check cache first
	urlKey := fmt.Sprintf("%s:%s:url", providerName, id)
	if config.CacheEnabled {
		if songURL, err := redis.GetString(urlKey); err == nil {
			return makeSingleSongInM3U(songURL, "")
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

	return makeSingleSongInM3U(fmt.Sprintf("%s/%s/%s/%s", baseURL, song.Provider, song.ID, url.PathEscape(filename)), song.Title)
}

func generateM3ULink(c *gin.Context) {
	u := c.Query("u")
	if u == "" {
		c.AbortWithError(http.StatusNotFound, errInvalidURL)
		return
	}
	if makeM3U(c, u, util.PlaylistMatch, makePlaylist) ||
		makeM3U(c, u, util.AlbumMatch, makeAlbumSongs) ||
		makeM3U(c, u, util.ArtistMatch, makeArtistSongs) ||
		makeM3U(c, u, util.SingleSongMatch, makeSingleSong) {
		return
	}

	c.Data(http.StatusNotFound, "text/html; charset=UTF-8", notFoundPage)
}

// makeM3U make M3U playlist
// return value true - pattern matched, false - pattern not match
func makeM3U(c *gin.Context, u string, patternMatch patternMatchFunc, make makeFunc) bool {
	id, providerName, matched := patternMatch(u)
	if matched {
		b, err := make(c, id, providerName)
		if err != nil {
			c.Data(http.StatusNotFound, "text/html; charset=UTF-8", notFoundPage)
		} else {
			c.Writer.Header().Set(`Content-Disposition`, `attachment; filename="playlist.m3u"`)
			c.Data(http.StatusOK, "audio/x-mpegurl", b)
		}
		return true
	}
	return false
}
