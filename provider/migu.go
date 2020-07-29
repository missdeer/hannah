package provider

import (
	"fmt"
	"net/http"

	"github.com/missdeer/hannah/util"
)

type migu struct {
}

type miguSearchResult struct {
	Musics []struct {
		AlbumName  string `json:"albumName"`
		AlbumID    string `json:"albumId"`
		MP3        string `json:"mp3"`
		SongName   string `json:"songName"`
		Lyrics     string `json:"lyrics"`
		ID         string `json:"id"`
		Title      string `json:"title"`
		Cover      string `json:"cover"`
		SingerName string `json:"singerName"`
		Artist     string `json:"artist"`
	} `json:"musics"`
	Pgt     int    `json:"pgt"`
	Keyword string `json:"keyword"`
	PageNo  string `json:"pageNo"`
	Success bool   `json:"success"`
}

func (p *migu) Search(keyword string, page int, limit int) (SearchResult, error) {
	u := fmt.Sprintf("http://m.music.migu.cn/migu/remoting/scr_search_tag?type=2&keyword=%s&pgc=%d&rows=%d", keyword, page, limit)

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:78.0) Gecko/20100101 Firefox/78.0")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "http://migu.cn/")
	req.Header.Set("Origin", "http://migu.cn/")
	req.Header.Set("Accept-Language", "zh-CN,zh-HK;q=0.8,zh-TW;q=0.6,en-US;q=0.4,en;q=0.2")
	req.Header.Set("Accept-Encoding", "gzip, deflate")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, ErrStatusNotOK
	}

	content, err := util.ReadHttpResponseBody(resp)
	if err != nil {
		return nil, err
	}

	var sr miguSearchResult
	err = json.Unmarshal(content, &sr)
	if err != nil {
		return nil, err
	}

	var res SearchResult
	for _, music := range sr.Musics {
		res = append(res, Song{
			ID:       music.ID,
			URL:      music.MP3,
			Title:    music.Title,
			Image:    music.Cover,
			Artist:   music.Artist,
			Lyric:    music.Lyrics,
			Provider: "migu",
		})
	}

	return res, nil
}

func (p *migu) ResolveSongURL(song Song) (Song, error) {
	return song, nil
}

func (p *migu) ResolveSongLyric(song Song) (Song, error) {
	return song, nil
}

func (p *migu) HotPlaylist(page int) (Playlists, error) {
	return nil, nil
}

func (p *migu) PlaylistDetail(pl Playlist) (Songs, error) {
	return nil, nil
}

func (p *migu) Name() string {
	return "migu"
}
