package decode

import (
	"io"
	"strings"

	"github.com/faiface/beep"
	"github.com/faiface/beep/flac"
	"github.com/faiface/beep/vorbis"
	"github.com/faiface/beep/wav"
)

type builtinDecoder func(io.ReadCloser) (beep.StreamSeekCloser, beep.Format, error)

var (
	builtinDecoderMap = map[string]builtinDecoder{
		".mp3":  Mpg123Decode,
		".ogg":  vorbis.Decode,
		".flac": func(r io.ReadCloser) (beep.StreamSeekCloser, beep.Format, error) { return flac.Decode(r) },
		".wav":  func(r io.ReadCloser) (beep.StreamSeekCloser, beep.Format, error) { return wav.Decode(r) },
		".m4a":  FAADDecode,
		".aac":  FAADDecode,
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
