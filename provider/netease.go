package provider

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/missdeer/hannah/util"
)

type netease struct {
}

type neteaseMusicInfo struct {
	ID          int     `json:"id"`
	Size        int     `json:"size"`
	Extension   string  `json:"extension"`
	SR          int     `json:"sr"`
	DFSID       int     `json:"dfsId"`
	Bitrate     int     `json:"bitrate"`
	PlayTime    int     `json:"playTime"`
	VolumeDelta float64 `json:"volumeDelta"`
}

type neteaseArtist struct {
	Name      string `json:"name"`
	ID        int    `json:"id"`
	PicID     int    `json:"PicId"`
	Img1v1ID  int    `json:"img1v1Id"`
	BriefDesc string `json:"briefDesc"`
	PicURL    string `json:"picUrl"`
	Img1v1URL string `json:"img1v1Url"`
}

type neteaseSearchResult struct {
	Result struct {
		Songs []struct {
			Name    string `json:"name"`
			ID      int    `json:"id"`
			MP3URL  string `json:"mp3Url"`
			MVID    int    `json:"mvid"`
			BMusic  neteaseMusicInfo
			HMusic  neteaseMusicInfo
			MMusic  neteaseMusicInfo
			LMusic  neteaseMusicInfo
			Artists []neteaseArtist `json:"artists"`
			Album   struct {
				Name       string        `json:"name"`
				ID         int           `json:"id"`
				Type       string        `json:"type"`
				Size       int           `json:"size"`
				PicID      int           `json:"PicId"`
				BlurPicURL string        `json:"blurPicUrl"`
				PicURL     string        `json:"picUrl"`
				Artist     neteaseArtist `json:"artist"`
			} `json:"album"`
		} `json:"songs"`
		SongCount int `json:"songCount"`
	} `json:"result"`
	Code int `json:"code"`
}

func (p *netease) Search(keyword string, page int, limit int) (SearchResult, error) {
	u := `http://music.163.com/api/search/pc`
	body := fmt.Sprintf("limit=%d&offset=%d&s=%s&type=1", limit, page*limit, url.QueryEscape(keyword))
	req, err := http.NewRequest("POST", u, strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:78.0) Gecko/20100101 Firefox/78.0")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "http://music.163.com/")
	req.Header.Set("Origin", "http://music.163.com/")
	req.Header.Set("Accept-Language", "zh-CN,zh-HK;q=0.8,zh-TW;q=0.6,en-US;q=0.4,en;q=0.2")
	req.Header.Set("Accept-Encoding", "gzip, deflate")

	client := util.GetHttpClient()

	resp, err := client.Do(req)
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

	var sr neteaseSearchResult
	err = json.Unmarshal(content, &sr)
	if err != nil {
		return nil, err
	}

	var res SearchResult
	for _, r := range sr.Result.Songs {
		var artists []string
		for _, a := range r.Artists {
			artists = append(artists, a.Name)
		}
		res = append(res, Song{
			ID:       strconv.Itoa(r.ID),
			Title:    r.Name,
			Image:    r.Album.PicURL,
			Artist:   strings.Join(artists, "/"),
			Provider: "netease",
		})
	}

	return res, nil
}

func (p *netease) SongDetail(song Song) (Song, error) {
	song.URL = fmt.Sprintf(`http://music.163.com/song/media/outer/url?id=%s.mp3`, song.ID)
	return song, nil
}

func (p *netease) HotPlaylist(page int) (Playlists, error) {
	return nil, nil
}

func (p *netease) PlaylistDetail(pl Playlist) (Songs, error) {
	return nil, nil
}

func (p *netease) Name() string {
	return "netease"
}
