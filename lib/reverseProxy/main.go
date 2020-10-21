package main

/*
#include <stdlib.h>
*/
import "C"
import (
	"log"
	"math/rand"
	"time"
	"unsafe"

	"github.com/missdeer/hannah/config"
	"github.com/missdeer/hannah/rp"
)

var (
	// Gitcommit contains the commit where we built reverseProxy from.
	GitCommit string
)

func main() {}

//export Free
func Free(c *C.char) {
	C.free(unsafe.Pointer(c))
}

//export StartReverseProxy
func StartReverseProxy(addr, limit string) bool {
	rand.Seed(time.Now().UnixNano())

	config.NetworkTimeout = 0 // no timeout, streaming costs much time
	if err := rp.Init(config.CacheAddr); err != nil {
		log.Println(err)
		return false
	}
	err := rp.Start(addr, limit)
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}

//export StopReverseProxy
func StopReverseProxy() {
	rp.Stop()
}
