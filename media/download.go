package media

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/missdeer/hannah/config"
	"github.com/missdeer/hannah/media/decode"
	"github.com/missdeer/hannah/provider"
	"github.com/missdeer/hannah/util"
)

var (
	ErrNotDir     = errors.New("path exists but not a directory")
	ErrFileExists = errors.New("file exists")
)

func downloadSong(song provider.Song, done chan string) error {
	stat, err := os.Stat(config.DownloadDir)
	if os.IsNotExist(err) {
		if err = os.MkdirAll(config.DownloadDir, 0755); err != nil {
			return err
		}
	} else if !stat.IsDir() {
		return ErrNotDir
	}

	filename := fmt.Sprintf("%s-%s%s", song.Title, song.Artist, decode.GetExtName(song.URL))
	fn := filepath.Join(config.DownloadDir, filename)
	if _, err = os.Stat(fn); !os.IsNotExist(err) {
		return ErrFileExists
	}

	req, err := http.NewRequest("GET", song.URL, nil)
	u, err := url.Parse(song.URL)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:78.0) Gecko/20100101 Firefox/78.0")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", fmt.Sprintf("%s://%s", u.Scheme, u.Host))
	req.Header.Set("Origin", fmt.Sprintf("%s://%s", u.Scheme, u.Host))
	req.Header.Set("Accept-Language", "zh-CN,zh-HK;q=0.8,zh-TW;q=0.6,en-US;q=0.4,en;q=0.2")
	req.Header.Set("Accept-Encoding", "gzip, deflate")

	client := util.GetHttpClient()
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	f, err := os.Create(fn)
	if err != nil {
		return err
	}

	defer func() {
		f.Close()
		done <- filename
	}()
	return util.CopyHttpResponseBody(resp, f)
}
