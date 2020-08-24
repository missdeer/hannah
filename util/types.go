package util

import (
	"net/url"
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
	if strings.HasPrefix(u, "http://") || strings.HasPrefix(u, "https://") {
		if u, err := url.Parse(uri); err == nil {
			for k, v := range extNameMimeTypeMap {
				if strings.HasSuffix(strings.ToLower(u.Path), k) {
					return k, v
				}
			}

		}
	} else {
		for k, v := range extNameMimeTypeMap {
			if strings.HasSuffix(u, k) {
				return k, v
			}
		}
	}
	return "", "application/octet-stream"
}
