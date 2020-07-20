package config

import (
	"strings"
)

var (
	Shuffle     bool
	Repeat      bool
	Action      string
	Provider    string
	Socks5Proxy string
	HttpProxy   string
	Limit       int
	Page        int

	supportedExtensions = map[string]struct{}{
		".mp3":  {},
		".flac": {},
		".wav":  {},
		".ogg":  {},
	}
)

func SupportedFileType(ext string) bool {
	_, ok := supportedExtensions[strings.ToLower(ext)]
	return ok
}
