package rp

import (
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/missdeer/hannah/input"
	"github.com/missdeer/hannah/provider"
)

var (
	mimeTypes = map[string]string{
		".mp3":  "audio/mpeg",
		".m4a":  "audio/mp4",
		".aac":  "audio/aac",
		".flac": "audio/flac",
		".ape":  "audio/x-ape",
		".wav":  "audio/wav",
		".ogg":  "audio/ogg",
		".m3u":  "audio/x-mpegurl",
	}
)

func getExtName(uri string) string {
	for k := range mimeTypes {
		if strings.HasPrefix(uri, "http://") || strings.HasPrefix(uri, "https://") {
			if strings.Contains(uri, k) {
				return k
			}
		} else {
			if strings.HasSuffix(uri, k) {
				return k
			}
		}
	}
	return ""
}

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

	ext := getExtName(song.URL)
	mt, ok := mimeTypes[ext]
	if !ok {
		mt = "application/octet-stream"
		log.Println(providerName, id, song, song.URL)
	}
	c.Writer.Header().Set("Content-Type", mt)
	c.Data(http.StatusOK, mt, []byte{})
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

	// TODO cache the resolved result for 30m ~ 60m

	r, err := input.OpenSource(song.URL)
	if err != nil {
		c.Abort()
		return
	}
	defer r.Close()

	ext := getExtName(song.URL)
	mt, ok := mimeTypes[ext]
	if !ok {
		mt = "application/octet-stream"
	}
	c.Writer.Header().Set("Content-Type", mt)
	c.Stream(func(w io.Writer) bool {
		_, e := io.Copy(w, r)
		return e == nil
	})
}

func StartReverseProxy(addr string) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/:provider/:id/:filename", getSong)
	r.HEAD("/:provider/:id/:filename", getSongInfo)
	r.GET("/:provider/:id", getSong)
	r.HEAD("/:provider/:id", getSongInfo)
	log.Fatal(r.Run(addr))
}
