package decode

import (
	"io"
	"strings"

	"github.com/faiface/beep"
	"github.com/faiface/beep/flac"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/vorbis"
	"github.com/faiface/beep/wav"

	"github.com/missdeer/hannah/config"
)

type builtinDecoder func(io.ReadCloser) (beep.StreamSeekCloser, beep.Format, error)

var (
	builtinDecoderMap = map[string]builtinDecoder{
		".mp3": func(r io.ReadCloser) (beep.StreamSeekCloser, beep.Format, error) {
			if config.Mpg123 {
				return Mpg123Decode(r)
			} else {
				return mp3.Decode(r)
			}
		},
		".ogg":  func(r io.ReadCloser) (beep.StreamSeekCloser, beep.Format, error) { return vorbis.Decode(r) },
		".flac": func(r io.ReadCloser) (beep.StreamSeekCloser, beep.Format, error) { return flac.Decode(r) },
		".wav":  func(r io.ReadCloser) (beep.StreamSeekCloser, beep.Format, error) { return wav.Decode(r) },
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
