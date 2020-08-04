package beep

import (
	"strings"
)

var (
	builtinSupportedExtensions = map[string]struct{}{
		".mp3":  {},
		".flac": {},
		".wav":  {},
		".ogg":  {},
		".acc":  {},
		".m4a":  {},
		".m3u":  {},
	}
)

func SupportedFileType(ext string) bool {
	_, ok := builtinSupportedExtensions[strings.ToLower(ext)]
	return ok
}
