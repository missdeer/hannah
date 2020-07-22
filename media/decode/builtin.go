package decode

import (
	"io"
	"strings"

	"github.com/faiface/beep"
	"github.com/faiface/beep/flac"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/vorbis"
	"github.com/faiface/beep/wav"
)

type builtinDecoder func(io.ReadCloser) (beep.StreamSeekCloser, beep.Format, error)

var (
	builtinDecoderMap = map[string]builtinDecoder{
		".mp3":  mp3.Decode,
		".ogg":  vorbis.Decode,
		".flac": func(r io.ReadCloser) (beep.StreamSeekCloser, beep.Format, error) { return flac.Decode(r) },
		".wav":  func(r io.ReadCloser) (beep.StreamSeekCloser, beep.Format, error) { return wav.Decode(r) },
	}

	builtinSupportedExtensions = map[string]struct{}{
		".mp3":  {},
		".flac": {},
		".wav":  {},
		".ogg":  {},
		".m3u":  {},
	}
)

func GetBuiltinDecoder(uri string) builtinDecoder {
	for k, v := range builtinDecoderMap {
		if strings.HasPrefix(uri, "http://") || strings.HasPrefix(uri, "https://") {
			if strings.Contains(uri, k) {
				return v
			}
		} else {
			if strings.HasSuffix(uri, k) {
				return v
			}
		}
	}
	return nil
}

func BuiltinSupportedFileType(ext string) bool {
	_, ok := builtinSupportedExtensions[strings.ToLower(ext)]
	return ok
}
