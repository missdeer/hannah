package rp

import (
	"errors"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"

	"github.com/missdeer/hannah/cache"
	"github.com/missdeer/hannah/config"
	"github.com/missdeer/hannah/util"
)

const (
	cacheTimeout = 4 * time.Hour
)

var (
	client                 *http.Client
	redis                  *cache.RedisCache
	errUnsupportedProvider = errors.New("unsupported provider")
	errInvalidSongID       = errors.New("invaild song ID")
	errInvalidURL          = errors.New("invalid URL")
	errLyricNotFound       = errors.New("lyric not found")
	notFoundPage           = []byte(`<html><script type="text/javascript" src="//qzonestyle.gtimg.cn/qzone/hybrid/app/404/search_children.js" charset="utf-8"></script><body></body></html>`)
)

func getSongPlaylist(c *gin.Context) {
	if config.ShowUserAgent {
		log.Println(c.Request.UserAgent())
	}
	p := c.Param("provider")
	if p == "m3u" {
		generateM3ULink(c)
		return
	}
	requestType := c.Query("type")
	switch strings.ToLower(requestType) {
	case "playlist":
		getPlaylist(c)
	case "artist":
		getArtistSongs(c)
	case "album":
		getAlbumSongs(c)
	case "search":
		searchSongs(c)
	default:
		fn := c.Param("filename")
		lyricForamt := strings.ToLower(filepath.Ext(fn))
		if lyricForamt == "" {
			lyricForamt = fn
		}
		switch lyricForamt {
		case "lrc", ".lrc", "smi", ".smi":
			getLyric(c, lyricForamt)
		default:
			getSong(c)
		}
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
	r.GET("/:provider/:id", getSongPlaylist)
	r.HEAD("/:provider/:id/:filename", getSongInfo)
	r.HEAD("/:provider/:id", getSongInfo)

	r.NoRoute(func(c *gin.Context) {
		c.Data(http.StatusNotFound, "text/html; charset=UTF-8", notFoundPage)
	})
	return r.Run(addr)
}
