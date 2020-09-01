package rp

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/missdeer/hannah/config"
	"github.com/missdeer/hannah/provider"
	"github.com/missdeer/hannah/util"
)

var (
	playersNotSupportRedirectURL = []string{
		"libmpv",                               // mpv
		"BASS",                                 // bass
		"mpg123",                               // mpg123
		"wbx 1.0.0; wbxapp 1.0.0; zhumu 4.0.0", // TTPlayer
		"TTPlayer",                             // TTPlayer
	}
	providersNotSupportRedirectURL = map[string]struct{}{
		"netease": {},
	}
)

func supportRedirectURL(providerName string, userAgent string) bool {
	if _, ok := providersNotSupportRedirectURL[providerName]; ok {
		return false
	}
	for _, k := range playersNotSupportRedirectURL {
		if strings.Contains(userAgent, k) {
			return false
		}
	}
	return true
}

func getLyric(c *gin.Context, format string) {
	providerName := c.Param("provider")
	id := c.Param("id")

	r := provider.GetSongIDPattern(providerName)
	if r == nil {
		c.AbortWithError(http.StatusNotFound, errUnsupportedProvider)
		return
	}
	if !r.MatchString(id) {
		c.AbortWithError(http.StatusNotFound, errInvalidSongID)
		return
	}
	// check cache first
	urlKey := fmt.Sprintf("%s:%s:%s:lyric", providerName, id, format)
	headerKey := fmt.Sprintf("%s:%s:header", providerName, id)
	if config.CacheEnabled {
		if h, err := redis.Get(headerKey); err == nil {
			if header, ok := h.(http.Header); ok {
				for k, v := range header {
					c.Writer.Header().Set(k, v[0])
				}
			}
		}

		if lyric, err := redis.GetString(urlKey); err == nil {
			c.Data(http.StatusOK, "application/octet-stream", []byte(lyric))
			return
		}
	}

	// resolve URL now
	p := provider.GetProvider(providerName)
	if p == nil {
		c.AbortWithError(http.StatusNotFound, errUnsupportedProvider)
		return
	}

	song, err := p.ResolveSongLyric(provider.Song{ID: id}, format)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	if song.Lyric == "" {
		c.AbortWithError(http.StatusNotFound, errLyricNotFound)
		return
	}

	// cache the resolved result
	if redis != nil {
		redis.Put(urlKey, song.Lyric)
	}

	c.Data(http.StatusOK, "application/octet-stream", []byte(song.Lyric))
}

func getSong(c *gin.Context) {
	providerName := c.Param("provider")
	id := c.Param("id")
	r := provider.GetSongIDPattern(providerName)
	if r == nil {
		c.AbortWithError(http.StatusNotFound, errUnsupportedProvider)
		return
	}
	if !r.MatchString(id) {
		c.AbortWithError(http.StatusNotFound, errInvalidSongID)
		return
	}
	canRedirect := supportRedirectURL(providerName, c.Request.UserAgent())
	// check cache first
	urlKey := fmt.Sprintf("%s:%s:url", providerName, id)
	headerKey := fmt.Sprintf("%s:%s:header", providerName, id)
	if config.CacheEnabled && config.RedirectURL && canRedirect {
		if h, err := redis.Get(headerKey); err == nil {
			if header, ok := h.(http.Header); ok {
				for k, v := range header {
					c.Writer.Header().Set(k, v[0])
				}
			}
		}

		if songURL, err := redis.GetString(urlKey); err == nil {
			c.Redirect(http.StatusFound, songURL)
			return
		}
	}

	// resolve URL now
	p := provider.GetProvider(providerName)
	if p == nil {
		c.AbortWithError(http.StatusNotFound, errUnsupportedProvider)
		return
	}

	song, err := p.ResolveSongURL(provider.Song{ID: id})
	for i := 0; i < config.ReverseProxyRetries && err != nil; i++ {
		song, err = p.ResolveSongURL(provider.Song{ID: id})
	}
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	if config.Debugging {
		log.Println(song.URL)
		log.Println("client ip:", c.ClientIP())
	}

	if canRedirect && (config.RedirectURL ||
		(!config.RedirectURL && config.AutoRedirectURL && (InChina(c.ClientIP()) || InLan(c.ClientIP())))) {
		c.Redirect(http.StatusFound, song.URL)
		return
	}

	// cache the resolved result
	if redis != nil {
		redis.PutWithTimeout(urlKey, song.URL, cacheTimeout)
	}

	req, err := http.NewRequest("GET", song.URL, nil)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	req.Header = c.Request.Header
	if r, err := url.Parse(song.URL); err == nil {
		req.Header.Set("Referer", fmt.Sprintf("%s://%s", r.Scheme, r.Hostname()))
		req.Header.Set("Origin", fmt.Sprintf("%s://%s", r.Scheme, r.Hostname()))
	}

	resp, err := client.Do(req)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}
	defer resp.Body.Close()

	if !strings.HasPrefix(resp.Header.Get("Content-Type"), "audio/") {
		_, mimeType := util.GetExtName(song.URL)
		resp.Header.Set("Content-Type", mimeType)
	}

	// cache the info
	if redis != nil {
		redis.PutWithTimeout(headerKey, resp.Header, cacheTimeout)
	}

	for k, v := range resp.Header {
		c.Writer.Header().Set(k, v[0])
	}

	c.Stream(func(w io.Writer) bool {
		_, e := io.Copy(w, resp.Body)
		return e != nil && e != io.EOF
	})
}
