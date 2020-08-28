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

func getPlaylist(c *gin.Context) {
	providerName := c.Param("provider")
	id := c.Param("id")
	refresh := c.Query("refresh")

	urlKey := fmt.Sprintf("%s:%s:playlist", providerName, id)
	if config.CacheEnabled && refresh != "1" {
		b, err := redis.GetBytes(urlKey)
		if err == nil {
			c.Writer.Header().Set(`Content-Disposition`, `attachment; filename="playlist.m3u"`)
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
	for i := 0; i < config.ReverseProxyRetries && err != nil; i++ {
		pld, err = p.PlaylistDetail(provider.Playlist{ID: id})
	}
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
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
	_, err = playlist.WriteSimpleTo(w)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}
	if err = w.Flush(); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if config.CacheEnabled {
		redis.Put(urlKey, b.Bytes())
	}

	c.Writer.Header().Set(`Content-Disposition`, `attachment; filename="playlist.m3u"`)
	c.Data(http.StatusOK, "audio/x-mpegurl", b.Bytes())
}
