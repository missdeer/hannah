package rp

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"

	"github.com/missdeer/hannah/config"
	"github.com/missdeer/hannah/provider"
	"github.com/missdeer/hannah/util"
)

var (
	client *http.Client
)

func getSongInfo(c *gin.Context) {
	providerName := c.Param("provider")
	id := c.Param("id")
	// TODO check cache first

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

	for k, v := range resp.Header {
		c.Writer.Header().Set(k, v[0])
	}
	c.Data(http.StatusOK, resp.Header.Get("Content-Type"), data)
}

func getSong(c *gin.Context) {
	providerName := c.Param("provider")
	id := c.Param("id")
	// TODO check cache first

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

	// TODO cache the resolved result for 30m ~ 60m

	req, err := http.NewRequest("GET", song.URL, nil)
	if err != nil {
		c.Abort()
		return
	}

	req.Header = c.Request.Header
	if r, err := url.Parse(song.URL); err == nil {
		req.Header.Set("Referer", fmt.Sprintf("%s://%s", r.Scheme, r.Hostname()))
	}

	resp, err := client.Do(req)
	if err != nil {
		c.Abort()
		return
	}
	defer resp.Body.Close()

	for k, v := range resp.Header {
		c.Writer.Header().Set(k, v[0])
	}
	c.Stream(func(w io.Writer) bool {
		_, e := io.Copy(w, resp.Body)
		return e == nil
	})
}

func StartReverseProxy(addr string) {
	client = util.GetHttpClient()

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/:provider/:id/:filename", getSong)
	r.HEAD("/:provider/:id/:filename", getSongInfo)
	r.GET("/:provider/:id", getSong)
	r.HEAD("/:provider/:id", getSongInfo)
	log.Fatal(r.Run(addr))
}
