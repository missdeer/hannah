package media

import (
	"fmt"
	"os"

	"github.com/ushis/m3u"

	"github.com/missdeer/hannah/config"
	"github.com/missdeer/hannah/provider"
)

func insertToM3U(song provider.Song) error {
	if _, err := os.Stat(config.M3UFileName); os.IsNotExist(err) {
		f, err := os.Create(config.M3UFileName)
		if err != nil {
			return err
		}
		f.Close()
	}

	f, err := os.OpenFile(config.M3UFileName, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	f.Seek(0, 0)
	pl, err := m3u.Parse(f)
	if err != nil {
		return err
	}

	track := m3u.Track{
		Path:  song.URL,
		Title: song.Title,
	}
	if song.Provider != "local filesystem" && song.Provider != "http(s)" {
		track.Path = fmt.Sprintf("%s://%s", song.Provider, song.ID)
	}

	for _, t := range pl {
		if t.Path == track.Path {
			return nil
		}
	}

	pl = append(pl, track)

	f.Seek(0, 0)
	if _, err := pl.WriteTo(f); err != nil {
		return err
	}

	return nil
}
