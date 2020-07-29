package provider

import (
	"fmt"
	"net/http"

	"github.com/missdeer/hannah/util"
)

const (
	kugouAPISearch   = `http://songsearch.kugou.com/song_search_v2?keyword=%s&page=%d&pagesize=%d`
	kugouAPISongInfo = `http://m.kugou.com/app/i/getSongInfo.php?cmd=playInfo&hash=%s`
)

type kugou struct {
}

type kugouSearchResult struct {
	Status    int `json:"status"`
	ErrorCode int `json:"error_code"`
	Data      struct {
		Total    int `json:"total"`
		Page     int `json:"page"`
		PageSize int `json:"pagesize"`
		Lists    []struct {
			ExtName     string `json:"ExtName"`
			OriSongName string `json:"OriSongName"`
			AlbumID     string `json:"AlbumID"`
			MixSongID   string `json:"MixSongID"`
			ID          string `json:"ID"`
			FileName    string `json:"FileName"`
			SongName    string `json:"SongName"`
			SingerName  string `json:"SingerName"`
			FileHash    string `json:"FileHash"`
			HQFileHash  string `json:"HQFileHash"`
		} `json:"lists"`
	} `json:"data"`
}

type kugouSongInfo struct {
	URL        string `json:"url"`
	SingerName string `json:"singerName"`
	ExtName    string `json:"extname"`
	Hash       string `json:"hash"`
	ImgURL     string `json:"imgUrl"`
	FileName   string `json:"fileName"`
	SongName   string `json:"songName"`
	Extra      struct {
		Hash128 string `json:"128hash"`
		Hash320 string `json:"320hash"`
	} `json:"extra"`
}

func (p *kugou) Search(keyword string, page int, limit int) (SearchResult, error) {
	u := fmt.Sprintf(kugouAPISearch, keyword, page, limit)
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:78.0) Gecko/20100101 Firefox/78.0")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "http://www.kugou.com/")
	req.Header.Set("Origin", "http://www.kugou.com/")
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

	var sr kugouSearchResult
	err = json.Unmarshal(content, &sr)
	if err != nil {
		return nil, err
	}

	var res SearchResult
	for _, r := range sr.Data.Lists {
		res = append(res, Song{
			ID:       r.FileHash,
			Title:    r.SongName,
			Artist:   r.SingerName,
			Provider: "kugou",
		})
	}

	return res, nil
}

func (p *kugou) ResolveSongURL(song Song) (Song, error) {
	u := fmt.Sprintf(kugouAPISongInfo, song.ID)
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return song, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:78.0) Gecko/20100101 Firefox/78.0")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "http://www.kugou.com/")
	req.Header.Set("Origin", "http://www.kugou.com/")
	req.Header.Set("Accept-Language", "zh-CN,zh-HK;q=0.8,zh-TW;q=0.6,en-US;q=0.4,en;q=0.2")
	req.Header.Set("Accept-Encoding", "gzip, deflate")

	resp, err := httpClient.Do(req)
	if err != nil {
		return song, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return song, ErrStatusNotOK
	}

	content, err := util.ReadHttpResponseBody(resp)
	if err != nil {
		return song, err
	}

	var si kugouSongInfo
	if err = json.Unmarshal(content, &si); err != nil {
		return song, err
	}
	song.URL = si.URL
	song.Image = si.ImgURL

	return song, nil
}

func (p *kugou) ResolveSongLyric(song Song) (Song, error) {
	return song, nil
}

func (p *kugou) HotPlaylist(page int) (Playlists, error) {
	return nil, nil
}

func (p *kugou) PlaylistDetail(pl Playlist) (Songs, error) {
	return nil, nil
}

func (p *kugou) Name() string {
	return "kugou"
}
