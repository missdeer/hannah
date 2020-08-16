package provider

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/missdeer/hannah/util"
)

const (
	kugouAPISearch         = `http://songsearch.kugou.com/song_search_v2?keyword=%s&page=%d&pagesize=%d`
	kugouAPISongInfo       = `http://m.kugou.com/app/i/getSongInfo.php?cmd=playInfo&hash=%s`
	kugouAPIHot            = `http://m.kugou.com/plist/index&json=true&page=%d`
	kugouAPIPlaylistDetail = `http://m.kugou.com/plist/list/%s?json=true&page=%d`
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

func (p *kugou) SearchSongs(keyword string, page int, limit int) (SearchResult, error) {
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

	var si kugouSongInfo
	if err = json.Unmarshal(content, &si); err != nil {
		return song, err
	}
	song.URL = si.URL
	song.Image = si.ImgURL
	song.Artist = si.SingerName
	song.Title = si.SongName

	return song, nil
}

func (p *kugou) ResolveSongLyric(song Song) (Song, error) {
	return song, nil
}

type kugouHotPlaylist struct {
	PList struct {
		List struct {
			HasNext int `json:"has_next"`
			Total   int `json:"total"`
			Info    []struct {
				Intro       string `json:"intro"`
				Img         string `json:"imgurl"`
				SpecialID   int    `json:"specialid"`
				SUID        int    `json:"suid"`
				SpecialName string `json:"specialname"`
			} `json:"info"`
		} `json:"list"`
	} `json:"plist"`
	PageSize int `json:"pagesize"`
}

func (p *kugou) HotPlaylist(page int, limit int) (Playlists, error) {
	u := fmt.Sprintf(kugouAPIHot, (page-1)*32)

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

	var pld kugouHotPlaylist
	if err = json.Unmarshal(content, &pld); err != nil {
		return nil, err
	}
	var res Playlists
	for _, pl := range pld.PList.List.Info {
		res = append(res, Playlist{
			ID:       strconv.Itoa(pl.SpecialID),
			Image:    pl.Img,
			Title:    pl.SpecialName,
			Provider: "kugou",
		})
	}

	return res, nil
}

type kugouPlaylistDetail struct {
	List struct {
		List struct {
			Info []struct {
				Hash    string `json:"hash"`
				ExtName string `json:"extname"`
			} `json:"info"`
			Total int `json:"total"`
		} `json:"list"`
		PageSize int `json:"pagesize"`
		Page     int `json:"page"`
	} `json:"list"`
}

func (p *kugou) PlaylistDetail(pl Playlist) (songs Songs, err error) {
	for page := 1; ; page++ {
		u := fmt.Sprintf(kugouAPIPlaylistDetail, pl.ID, page)

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

		var pld kugouPlaylistDetail
		if err = json.Unmarshal(content, &pld); err != nil {
			return nil, err
		}
		for _, p := range pld.List.List.Info {
			songs = append(songs, Song{
				ID:       p.Hash,
				Provider: "kugou",
			})
		}
		if len(songs) == pld.List.List.Total {
			break
		}
	}
	return songs, nil
}

func (p *kugou) ArtistSongs(id string) (res Songs, err error) {
	return nil, ErrNotImplemented
}

func (p *kugou) AlbumSongs(id string) (res Songs, err error) {
	return nil, ErrNotImplemented
}

func (p *kugou) Login() error {
	return  ErrNotImplemented
}

func (p *kugou) Name() string {
	return "kugou"
}
