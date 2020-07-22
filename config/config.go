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
	Player      string

	builtinSupportedExtensions = map[string]struct{}{
		".mp3":  {},
		".flac": {},
		".wav":  {},
		".ogg":  {},
	}
)

func BuiltinSupportedFileType(ext string) bool {
	_, ok := builtinSupportedExtensions[strings.ToLower(ext)]
	return ok
}
