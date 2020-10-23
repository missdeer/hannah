package rp

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
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
	srv                    *http.Server
	client                 *http.Client
	redis                  *cache.RedisCache
	errUnsupportedProvider = errors.New("unsupported provider")
	errInvalidSongID       = errors.New("invaild song ID")
	errInvalidURL          = errors.New("invalid URL")
	errLyricNotFound       = errors.New("lyric not found")
	notFoundPage           = []byte(`<html><script type="text/javascript" src="//qzonestyle.gtimg.cn/qzone/hybrid/app/404/search_children.js" charset="utf-8"></script><body></body></html>`)
)

func getSongPlaylist(c *gin.Context) {
	if config.Debugging {
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
		} else {
			lyricForamt = lyricForamt[1:]
		}
		switch lyricForamt {
		case "lrc", "smi":
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

func StartDaemon(addr string, limit string) {
	Start(addr, limit)

	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	Stop()
}

func Start(addr string, limit string) {
	if srv != nil {
		Stop()
	}
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

	srv = &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
}

func Stop() {
	if srv == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
