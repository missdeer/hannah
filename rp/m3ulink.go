package rp

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	patterns = []string{
		`https://music.163.com/#/discover/toplist?id=2884035`,
		`https://music.163.com/#/playlist?id=5067260983`,
		`https://music.163.com/#/my/m/music/playlist?id=41220530`,
		`https://music.163.com/#/song?id=1463165983`,
	}
)

func generateM3ULink(c *gin.Context) {
	u := c.Query("u")
	if u == "" {
		c.AbortWithError(http.StatusNotFound, errInvalidURL)
		return
	}
}
