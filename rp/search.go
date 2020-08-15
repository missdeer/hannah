package rp

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
	"github.com/ushis/m3u"

	"github.com/missdeer/hannah/config"
	"github.com/missdeer/hannah/provider"
)

func searchSongs(c *gin.Context) {
	providerName := c.Param("provider")
	keyword := c.Param("id")
	refresh := c.Query("refresh")
	page := c.DefaultQuery("page", "1")
	pageNr, err := strconv.Atoi(page)
	if err != nil {
		pageNr = 1
	}
	limit := c.DefaultQuery("limit", "50")
	limitNr, err := strconv.Atoi(limit)
	if err != nil {
		limitNr = 50
	}

	urlKey := fmt.Sprintf("%s:%s:search", providerName, keyword)
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

	sr, err := p.Search(keyword, pageNr, limitNr)
	for i := 0; i < config.ReverseProxyRetries && err != nil; i++ {
		sr, err = p.Search(keyword, pageNr, limitNr)
	}
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
	for _, song := range sr {
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

	if config.CacheEnabled {
		redis.Put(urlKey, b.Bytes())
	}

	c.Writer.Header().Set(`Content-Disposition`, `attachment; filename="playlist.m3u"`)
	c.Data(http.StatusOK, "audio/x-mpegurl", b.Bytes())
}
