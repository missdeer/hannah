package rp

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
	"github.com/ushis/m3u"

	"github.com/missdeer/hannah/config"
	"github.com/missdeer/hannah/provider"
)

func getPlaylist(c *gin.Context) {
	providerName := c.Param("provider")
	id := c.Param("id")
	refresh := c.Query("refresh")

	urlKey := fmt.Sprintf("%s:%s:playlist", providerName, id)
	if config.CacheEnabled && refresh != "1" {
		b, err := redis.GetBytes(urlKey)
		if err == nil {
			c.Data(http.StatusOK, "audio/x-mpegurl", b)
			return
		}
	}

	// resolve playlist
	p := provider.GetProvider(providerName)
	if p == nil {
		c.AbortWithError(http.StatusNotFound, errUnsupportedProvider)
		return
	}

	pld, err := p.PlaylistDetail(provider.Playlist{ID: id})
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
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
	_, err = playlist.WriteTo(w)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	if config.CacheEnabled {
		redis.Put(urlKey, b.Bytes())
	}

	c.Data(http.StatusOK, "audio/x-mpegurl", b.Bytes())
}
