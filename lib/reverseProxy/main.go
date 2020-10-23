package main

/*
#include <stdlib.h>
*/
import "C"
import (
	"log"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/missdeer/hannah/config"
	"github.com/missdeer/hannah/rp"
)

var (
	// Gitcommit contains the commit where we built reverseProxy from.
	GitCommit string
)

func main() {}

//export StartReverseProxy
func StartReverseProxy(addr, limit string) bool {
	rand.Seed(time.Now().UnixNano())
	gin.SetMode(gin.ReleaseMode)
	config.NetworkTimeout = 0 // no timeout, streaming costs much time
	if err := rp.Init(config.CacheAddr); err != nil {
		log.Println(err)
		return false
	}
	rp.Start(addr, limit)
	return true
}

//export StopReverseProxy
func StopReverseProxy() {
	rp.Stop()
}

//export SetSocks5Proxy
func SetSocks5Proxy(socks5 string) {
	config.Socks5Proxy = socks5
}

//export SetHttpProxy
func SetHttpProxy(httpProxy string) {
	config.HttpProxy = httpProxy
}

//export SetNetworkInterface
func SetNetworkInterface(networkInterface string) {
	config.NetworkInterface = networkInterface
}

//export SetAutoRedirect
func SetAutoRedirect(autoRedirect bool) {
	config.AutoRedirectURL = autoRedirect
}

//export SetRedirect
func SetRedirect(redirect bool) {
	config.RedirectURL = redirect
}
