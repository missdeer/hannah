package rp

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/missdeer/hannah/cache"
	"github.com/missdeer/hannah/config"
	"github.com/missdeer/hannah/provider"
	"github.com/missdeer/hannah/util"
)

const (
	cacheTimeout = 4 * time.Hour
)

var (
	client *http.Client
	redis  *cache.RedisCache
)

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
		c.Abort()
		return
	}
	song, err := p.ResolveSongURL(provider.Song{ID: id})
	if err != nil {
		c.Abort()
		return
	}

	req, err := http.NewRequest("HEAD", song.URL, nil)
	if err != nil {
		c.Abort()
		return
	}
	req.Header = c.Request.Header

	resp, err := client.Do(req)
	if err != nil {
		c.Abort()
		return
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.Abort()
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

func getSong(c *gin.Context) {
	providerName := c.Param("provider")
	id := c.Param("id")

	// check cache first
	urlKey := fmt.Sprintf("%s:%s:url", providerName, id)
	headerKey := fmt.Sprintf("%s:%s:header", providerName, id)
	if config.CacheEnabled && config.RedirectURL {
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
		c.Abort()
		return
	}
	song, err := p.ResolveSongURL(provider.Song{ID: id})
	if err != nil {
		c.Abort()
		return
	}

	if config.RedirectURL {
		c.Redirect(http.StatusFound, song.URL)
		return
	}

	// cache the resolved result
	if redis != nil {
		redis.PutWithTimeout(urlKey, song.URL, cacheTimeout)
	}

	req, err := http.NewRequest("GET", song.URL, nil)
	if err != nil {
		c.Abort()
		return
	}

	req.Header = c.Request.Header
	if r, err := url.Parse(song.URL); err == nil {
		req.Header.Set("Referer", fmt.Sprintf("%s://%s", r.Scheme, r.Hostname()))
		req.Header.Set("Origin", fmt.Sprintf("%s://%s", r.Scheme, r.Hostname()))
	}

	resp, err := client.Do(req)
	if err != nil {
		c.Abort()
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

func Init(addr string) error {
	client = util.GetHttpClient()

	var err error
	if config.CacheEnabled {
		redis, err = cache.RedisInit(addr)
	}
	return err
}

func Start(addr string) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/:provider/:id/:filename", getSong)
	r.HEAD("/:provider/:id/:filename", getSongInfo)
	r.GET("/:provider/:id", getSong)
	r.HEAD("/:provider/:id", getSongInfo)
	log.Fatal(r.Run(addr))
}
