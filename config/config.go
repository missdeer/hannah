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

	supportedExtentions = map[string]struct{}{
		".mp3":  {},
		".flac": {},
		".wav":  {},
		".ogg":  {},
	}
)

func SupportedFileType(ext string) bool {
	_, ok := supportedExtentions[strings.ToLower(ext)]
	return ok
}
