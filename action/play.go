package action

import (
	"bufio"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strings"

	"github.com/missdeer/golib/fsutil"

	"github.com/missdeer/hannah/config"
	"github.com/missdeer/hannah/media"
	"github.com/missdeer/hannah/media/decode/mpg123"
	"github.com/missdeer/hannah/provider"
)

var (
	supportedRemote = map[string]struct{}{
		"http://":  {},
		"https://": {},
	}
	supportedService = map[string]struct{}{
		"netease":  {},
		"qq":       {},
		"xiami":    {},
		"bilibili": {},
		"kugou":    {},
		"kuwo":     {},
		"migu":     {},
	}
)

func resolve(song string) provider.Song {
	// local filesystem
	if _, err := os.Stat(song); !os.IsNotExist(err) {
		f, err := os.Open(song)
		if err != nil {
			return provider.Song{
				URL:      song,
				Provider: "local filesystem",
			}
		}
		defer f.Close()

		r := mpg123.NewReaderConfig(bufio.NewReader(f), mpg123.ReaderConfig{
			OutputFormat: &mpg123.OutputFormat{
				Channels: 2,
				Rate:     44100,
				Encoding: mpg123.EncodingInt16,
			},
		})

		r.Read(nil)
		if r.Meta().ID3v2 != nil {
			return provider.Song{
				URL:      song,
				Artist:   r.Meta().ID3v2.Artist,
				Title:    r.Meta().ID3v2.Title,
				Provider: "local filesystem",
			}
		}
		return provider.Song{
			URL:      song,
			Provider: "local filesystem",
		}
	}
	// http/https
	for k, _ := range supportedRemote {
		if strings.HasPrefix(song, k) {
			return provider.Song{URL: song}
		}
	}
	// services
	ss := strings.Split(song, "://")
	if len(ss) == 2 {
		schema := ss[0]
		if _, ok := supportedService[schema]; ok {
			p := provider.GetProvider(schema)
			if s, err := p.ResolveSongURL(provider.Song{ID: ss[1]}); err == nil {
				return s
			}
		}
	}
	return provider.Song{URL: song}
}

func play(songs []string) error {
	for played := false; !played || config.Repeat; played = true {
		if config.Shuffle {
			rand.Shuffle(len(songs), func(i, j int) { songs[i], songs[j] = songs[j], songs[i] })
		}
		for i := 0; i < len(songs); i++ {
			song := resolve(songs[i])
			err := media.PlayMedia(song.URL, i+1, len(songs), song.Artist, song.Title) // TODO: extract from file name or ID3v1/v2 tag
			switch err {
			case media.ShouldQuit:
				return err
			case media.PreviousSong:
				i -= 2
			case media.NextSong:
			// auto next
			case media.UnsupportedMediaType:
				if b, e := fsutil.FileExists(config.Player); e == nil && b {
					log.Println(err, song, ", try to use external player", config.Player)
					cmd := exec.Command(config.Player, song.URL)
					cmd.Run()
				} else {
					log.Println(e, song)
				}
			default:
				log.Println(err)
			}
		}
	}
	return nil
}
