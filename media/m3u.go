package media

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/ushis/m3u"

	"github.com/missdeer/hannah/config"
	"github.com/missdeer/hannah/provider"
)

func composeTrack(pl m3u.Playlist, song provider.Song, origin bool) *m3u.Track {
	track := &m3u.Track{
		Path:  song.URL,
		Title: song.Title,
	}
	if song.Provider != "local filesystem" && song.Provider != "http(s)" && origin {
		if config.ReverseProxyEnabled {
			scheme := `http`
			host := config.ReverseProxy
			if u, err := url.Parse(config.ReverseProxy); err == nil {
				scheme = u.Scheme
				host = u.Host
			}
			filename := strings.Replace(fmt.Sprintf("%s - %s", song.Title, song.Artist), "/", "-", -1)
			track.Path = fmt.Sprintf("%s://%s/%s/%s/%s",
				scheme, host, song.Provider, song.ID, url.PathEscape(filename))
		} else {
			track.Path = fmt.Sprintf("%s://%s", song.Provider, song.ID)
		}
	}

	for _, t := range pl {
		if t.Path == track.Path {
			return nil
		}
	}

	return track
}

// AppendSongToM3U append song to M3U
// song song info struct
// origin origin URI or final URI, netease://12345 or https://music.163.com/12345.mp3
// done notifier
func AppendSongToM3U(song provider.Song, origin bool, done chan string) error {
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

	track := composeTrack(pl, song, origin)
	if track == nil {
		return nil
	}

	pl = append(pl, *track)

	f.Seek(0, 0)
	if _, err := pl.WriteSimpleTo(f); err != nil {
		return err
	}

	defer func() { done <- track.Path }()
	return nil
}

func AppendSongsToM3U(songs provider.Songs, origin bool) error {
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

	for _, song := range songs {
		if track := composeTrack(pl, song, origin); track != nil {
			pl = append(pl, *track)
		}
	}
	f.Seek(0, 0)
	if _, err := pl.WriteSimpleTo(f); err != nil {
		return err
	}

	return nil
}
