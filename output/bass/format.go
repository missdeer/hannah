package bass

import (
	"strings"
)

var (
	bassSupportedExtensions = map[string]struct{}{
		".mp3":  {},
		".flac": {},
		".wav":  {},
		".ogg":  {},
		".m3u":  {},
		".m4a":  {},
		".aac":  {},
		".wma":  {},
		".ape":  {},
		".ac3":  {},
		".webm": {},
	}
)

func SupportedFileType(ext string) bool {
	_, ok := bassSupportedExtensions[strings.ToLower(ext)]
	return ok
}
