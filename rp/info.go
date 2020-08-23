package rp

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/missdeer/hannah/config"
	"github.com/missdeer/hannah/provider"
	"github.com/missdeer/hannah/util"
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

	if !strings.HasPrefix(resp.Header.Get("Content-Type"), "audio/") {
		_, mimeType := util.GetExtName(song.URL)
		resp.Header.Set("Content-Type", mimeType)
	}

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
	fmt.Println(resp.Header.Get("Content-Type"))
	c.Data(http.StatusOK, resp.Header.Get("Content-Type"), data)
}
