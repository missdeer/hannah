package util

import (
	"strings"
)

var (
	extNameMimeTypeMap = map[string]string{
		".mp3":  "audio/mpeg",
		".ogg":  "audio/ogg",
		".flac": "audio/flac",
		".wav":  "audio/wave",
		".m4a":  "audio/mp4",
		".aac":  "audio/aac",
		".wma":  "audio/x-ms-wma",
		".ape":  "audio/x-ape",
	}
)

func GetExtName(uri string) (ext string, mimetype string) {
	u := strings.ToLower(uri)
	for k, v := range extNameMimeTypeMap {
		if strings.HasPrefix(u, "http://") || strings.HasPrefix(u, "https://") {
			if strings.Contains(u, k) {
				return k, v
			}
		} else {
			if strings.HasSuffix(u, k) {
				return k, v
			}
		}
	}
	return "", "application/octet-stream"
}
