package rp

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
	"github.com/ushis/m3u"

	"github.com/missdeer/hannah/cache"
	"github.com/missdeer/hannah/config"
	"github.com/missdeer/hannah/provider"
	"github.com/missdeer/hannah/util"
)

const (
	cacheTimeout = 4 * time.Hour
)

var (
	client                   *http.Client
	redis                    *cache.RedisCache
	errUnsupportedProvider   = errors.New("unsupported provider")
	errInvalidSongID         = errors.New("invaild song ID")
	playerSupportRedirectURL = map[string]bool{
		"foobar2000": true,
		"libmpv":     false,
	}
)

func supportRedirectURL(userAgent string) bool {
	for k, v := range playerSupportRedirectURL {
		if strings.Contains(userAgent, k) {
			return v
		}
	}
	return true
}

func getSongInfo(c *gin.Context) {
	providerName := c.Param("provider")
	id := c.Param("id")

	// check cache first
	headerKey := fmt.Sprintf("%s:%s:header", providerName, id)
	if config.CacheEnabled {
		if h, err := redis.Get(headerKey); err == nil {
			if header, ok := h.(http.Header); ok {
				for k, v := range header {
					c.Writer.Header().Set(k, v[0])
				}
				c.Data(http.StatusOK, header["Content-Type"][0], []byte{})
				return
			}
		}
	}

	// resolve URL now
	p := provider.GetProvider(providerName)
	if p == nil {
		c.AbortWithError(http.StatusNotFound, errUnsupportedProvider)
		return
	}
	song, err := p.ResolveSongURL(provider.Song{ID: id})
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	req, err := http.NewRequest("HEAD", song.URL, nil)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}
	req.Header = c.Request.Header

	resp, err := client.Do(req)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	// cache the info
	if redis != nil {
		redis.PutWithTimeout(headerKey, resp.Header, cacheTimeout)
	}

	for k, v := range resp.Header {
		c.Writer.Header().Set(k, v[0])
	}
	c.Data(http.StatusOK, resp.Header.Get("Content-Type"), data)
}

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

func getSong(c *gin.Context) {
	providerName := c.Param("provider")
	id := c.Param("id")
	canRedirect := supportRedirectURL(c.Request.UserAgent())
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
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	if canRedirect && (config.RedirectURL || (!config.RedirectURL && config.AutoRedirectURL && InChina(c.ClientIP()))) {
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

	// cache the info
	if redis != nil {
		redis.PutWithTimeout(headerKey, resp.Header, cacheTimeout)
	}

	for k, v := range resp.Header {
		c.Writer.Header().Set(k, v[0])
	}
	c.Stream(func(w io.Writer) bool {
		_, e := io.Copy(w, resp.Body)
		return e == nil
	})
}

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

func getSongPlaylist(c *gin.Context) {
	requestType := c.Query("type")
	switch strings.ToLower(requestType) {
	case "playlist":
		getPlaylist(c)
	case "search":
		searchSongs(c)
	default:
		getSong(c)
	}
}

func Init(addr string) error {
	client = util.GetHttpClient()

	err := LoadChinaIPList()
	if err != nil {
		return err
	}

	if config.CacheEnabled {
		redis, err = cache.RedisInit(addr)
	}
	return err
}

func Start(addr string, limit string) error {
	r := gin.New()
	if gin.Mode() != gin.ReleaseMode {
		r.Use(gin.Logger())
	}
	if limit != "" {
		r.Use(CIDR(limit))
	}
	r.Use(location.Default())
	r.Use(gin.Recovery())
	r.GET("/:provider/:id/:filename", getSongPlaylist)
	r.HEAD("/:provider/:id/:filename", getSongInfo)
	r.GET("/:provider/:id", getSongPlaylist)
	r.HEAD("/:provider/:id", getSongInfo)

	r.NoRoute(func(c *gin.Context) {
		c.Data(http.StatusNotFound, "text/html; charset=UTF-8", []byte(`<html><script type="text/javascript" src="//qzonestyle.gtimg.cn/qzone/hybrid/app/404/search_children.js" charset="utf-8"></script><body></body></html>`))
	})
	return r.Run(addr)
}
