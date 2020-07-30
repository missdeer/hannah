package bass

import (
	"strings"
)

var (
	bassSupportedExtensions = map[string]struct{}{
		".mp3":  {},
		".mp2":  {},
		".mp1":  {},
		".aiff": {},
		".wav":  {},
		".ogg":  {},
		".m3u":  {},
	}
)

func SupportedFileType(ext string) bool {
	_, ok := bassSupportedExtensions[strings.ToLower(ext)]
	return ok
}

func AddSupportedFileType(ext string) {
	e := strings.TrimSpace(ext)
	if e != "" {
		bassSupportedExtensions[e] = struct{}{}
	}
}
