package provider

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/missdeer/hannah/util"
)

const (
	bilibiliAPIHot            = `https://www.bilibili.com/audio/music-service-c/web/menu/hit?ps=%d&pn=%d`
	bilibiliAPIPlaylistDetail = `https://www.bilibili.com/audio/music-service-c/web/song/of-menu?pn=1&ps=100&sid=%s`
	bilibiliAPISongURL        = `https://www.bilibili.com/audio/music-service-c/web/url?sid=%s`
)

type bilibili struct {
}

func (p *bilibili) SearchSongs(keyword string, page int, limit int) (SearchResult, error) {
	// not supported
	return nil, nil
}

type bilibiliSongURL struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		CDNs []string `json:"cdns"`
	} `json:"data"`
}

func (p *bilibili) ResolveSongURL(song Song) (Song, error) {
	u := fmt.Sprintf(bilibiliAPISongURL, song.ID)
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return song, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:78.0) Gecko/20100101 Firefox/78.0")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "http://www.bilibili.com/")
	req.Header.Set("Origin", "http://www.bilibili.com/")
	req.Header.Set("Accept-Language", "zh-CN,zh-HK;q=0.8,zh-TW;q=0.6,en-US;q=0.4,en;q=0.2")
	req.Header.Set("Accept-Encoding", "gzip, deflate")

	httpClient := util.GetHttpClient()
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

	var songURL bilibiliSongURL
	if err = json.Unmarshal(content, &songURL); err != nil {
		return song, err
	}
	if len(songURL.Data.CDNs) == 0 {
		return song, errors.New("no song URL")
	}
	song.URL = songURL.Data.CDNs[0]
	return song, nil
}

func (p *bilibili) ResolveSongLyric(song Song) (Song, error) {
	return song, nil
}

type bilibiliHot struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		CurPage   int `json:"curPage"`
		PageCount int `json:"pageCount"`
		TotalSize int `json:"totalSize"`
		PageSize  int `json:"pageSize"`
		Data      []struct {
			Title  string `json:"title"`
			MenuID int    `json:"menuId"`
			Cover  string `json:"cover"`
		} `json:"data"`
	} `json:"data"`
}

func (p *bilibili) HotPlaylist(page int, limit int) (Playlists, error) {
	u := fmt.Sprintf(bilibiliAPIHot, limit, page)

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:78.0) Gecko/20100101 Firefox/78.0")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "http://www.bilibili.com/")
	req.Header.Set("Origin", "http://www.bilibili.com/")
	req.Header.Set("Accept-Language", "zh-CN,zh-HK;q=0.8,zh-TW;q=0.6,en-US;q=0.4,en;q=0.2")
	req.Header.Set("Accept-Encoding", "gzip, deflate")

	httpClient := util.GetHttpClient()
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

	var hot bilibiliHot
	if err = json.Unmarshal(content, &hot); err != nil {
		return nil, err
	}
	var pls Playlists
	for _, pl := range hot.Data.Data {
		pls = append(pls, Playlist{
			ID:       strconv.Itoa(pl.MenuID),
			Title:    pl.Title,
			Image:    pl.Cover,
			Provider: "bilibili",
		})
	}
	return pls, nil
}

type bilibiliPlaylistDetail struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		CurPage   int `json:"curPage"`
		PageCount int `json:"pageCount"`
		TotalSize int `json:"totalSize"`
		PageSize  int `json:"pageSize"`
		Data      []struct {
			Title string `json:"title"`
			ID    int    `json:"id"`
			Cover string `json:"cover"`
			UName string `json:"uname"`
			Lyric string `jsoN:"lyric"`
		} `json:"data"`
	} `json:"data"`
}

func (p *bilibili) PlaylistDetail(pl Playlist) (Songs, error) {
	u := fmt.Sprintf(bilibiliAPIPlaylistDetail, pl.ID)

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:78.0) Gecko/20100101 Firefox/78.0")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "http://www.bilibili.com/")
	req.Header.Set("Origin", "http://www.bilibili.com/")
	req.Header.Set("Accept-Language", "zh-CN,zh-HK;q=0.8,zh-TW;q=0.6,en-US;q=0.4,en;q=0.2")
	req.Header.Set("Accept-Encoding", "gzip, deflate")

	httpClient := util.GetHttpClient()
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

	var pld bilibiliPlaylistDetail
	if err = json.Unmarshal(content, &pld); err != nil {
		return nil, err
	}
	var songs Songs
	for _, pl := range pld.Data.Data {
		songs = append(songs, Song{
			ID:       strconv.Itoa(pl.ID),
			Title:    pl.Title,
			Artist:   pl.UName,
			Image:    pl.Cover,
			Lyric:    pl.Lyric,
			Provider: "bilibili",
		})
	}

	return songs, nil
}

func (p *bilibili) ArtistSongs(id string) (res Songs, err error) {
	return nil, ErrNotImplemented
}

func (p *bilibili) AlbumSongs(id string) (res Songs, err error) {
	return nil, ErrNotImplemented
}

func (p *bilibili) Login() error {
	return  ErrNotImplemented
}

func (p *bilibili) Name() string {
	return "bilibili"
}
