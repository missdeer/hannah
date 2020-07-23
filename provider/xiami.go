package provider

import (
	"crypto/md5"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/missdeer/hannah/util"
)

const (
	xiamiAppKey    = "23649156"
	xiamiApiSearch = "https://acs.m.xiami.com/h5/mtop.alimusic.search.searchservice.searchsongs/1.0/?appKey=23649156"
)

var (
	ErrTokenNotFound      = errors.New("xiami token not found")
	ErrCodeFieldNotExists = errors.New("code field does not exist")
	reqHeader             = map[string]interface{}{
		"appId":      200,
		"platformId": "h5",
	}
)

type xiami struct {
	token   string
	cookies []*http.Cookie
	client  *http.Client
}

type xiamiSongDetail struct {
	Status bool `json:"status"`
	Data   struct {
		TrackList []struct {
			SongID        string `json:"songId"`
			SongStringID  string `json:"songStringId"`
			SongName      string `json:"songName"`
			AlbumID       int    `json:"albumId"`
			AlbumStringID string `json:"albumStringId"`
			ArtistID      int    `json:"artistId"`
			Artist        string `json:"artist"`
			Singers       string `json:"singers"`
			Location      string `json:"location"`
			Pic           string `json:"pic"`
			LyricInfo     struct {
				LyricID   int    `json:"lyricId"`
				LyricFile string `json:"lyricFile"`
			} `json:"lyricInfo"`
		} `json:"trackList"`
	} `json:"data"`
}

type xiamiListenFile struct {
	DownloadFileSize string `json:"downloadFileSize"`
	FileSize         string `json:"fileSize"`
	Format           string `json:"format"`
	ListenFile       string `json:"listenFile"`
	PlayVolume       string `json:"playVolume"`
	Quality          string `json:"quality"`
}

type xiamiListenFiles []xiamiListenFile

func (s xiamiListenFiles) Len() int {
	return len(s)
}

func (s xiamiListenFiles) Less(i, j int) bool {
	if s[i].Format == s[j].Format && s[j].Format == `mp3` {
		if s[i].Quality == `h` { // only first 2 minutes are available in `h` quality
			return false
		}
		return true
	}
	if s[i].Format == s[j].Format && s[j].Format == `m4a` {
		if s[i].Quality == `f` {
			return false
		}
		return true
	}
	if s[i].Format == `mp3` {
		return false
	}
	return true
}

func (s xiamiListenFiles) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type xiamiSearchResult struct {
	API  string `json:"api"`
	Data struct {
		Data struct {
			PagingVO struct {
				Page     string `json:"page"`
				PageSize string `json:"pageSize"`
				Pages    string `json:"pages"`
				Count    string `json:"count"`
			} `json:"pagingVO"`
			Songs []struct {
				SongID        string           `json:"songId"`
				SongStringID  string           `json:"songStringId"`
				SongName      string           `json:"songName"`
				AlbumID       string           `json:"albumId"`
				AlbumStringID string           `json:"albumStringId"`
				AlbumLogo     string           `json:"albumLogo"`
				AlbumLogoS    string           `json:"albumLogoS"`
				AlbumName     string           `json:"albumName"`
				ArtistID      string           `json:"artistId"`
				ArtistName    string           `json:"artistName"`
				ArtistLogo    string           `json:"artistLogo"`
				Singers       string           `json:"singers"`
				ListenFiles   xiamiListenFiles `json:"listenFiles"`
			} `json:"songs"`
		} `json:"data"`
	} `json:"data"`
}

func (p *xiami) getToken(u string) (string, error) {
	if p.token != "" {
		return p.token, nil
	}
	parsedURL, _ := url.Parse(u)
	c := p.client.Jar.Cookies(parsedURL)
	const XiaMiToken = "_m_h5_tk"
	if c == nil || len(c) == 0 {
		req, err := http.NewRequest("GET", u, nil)
		if err != nil {
			return "", err
		}

		resp, err := p.client.Do(req)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return "", ErrStatusNotOK
		}

		c = resp.Cookies()
	}
	for _, cookie := range c {
		if cookie.Name == XiaMiToken {
			return strings.Split(cookie.Value, "_")[0], nil
		}
	}
	return "", ErrTokenNotFound
}

func signPayload(token string, model interface{}) (map[string]string, error) {
	payload := map[string]interface{}{
		"header": reqHeader,
		"model":  model,
	}
	r, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	data := map[string]string{
		"requestStr": string(r),
	}
	r, err = json.Marshal(data)
	if err != nil {
		return nil, err
	}

	t := time.Now().UnixNano() / (1e6)
	signStr := fmt.Sprintf("%s&%d&%s&%s", token, t, xiamiAppKey, string(r))
	sign := fmt.Sprintf("%x", md5.Sum([]byte(signStr)))

	return map[string]string{
		"t":    strconv.FormatInt(t, 10),
		"sign": sign,
		"data": string(r),
	}, nil
}

func caesar(location string) (string, error) {
	// https://github.com/listen1/listen1_chrome_extension/blob/f2e1b4376d3770333816668d98808ae90f796217/js/provider/xiami.js#L5
	num := int(location[0] - '0')
	avgLen := (len(location) - 1) / num
	remainder := (len(location) - 1) % num

	var result []string
	for i := 0; i < remainder; i++ {
		line := location[i*(avgLen+1)+1 : (i+1)*(avgLen+1)+1]
		result = append(result, line)
	}

	for i := 0; i < num-remainder; i++ {
		line := location[(avgLen+1)*remainder:]
		line = line[i*avgLen+1 : (i+1)*avgLen+1]
		result = append(result, line)
	}

	var s []byte
	for i := 0; i < avgLen; i++ {
		for j := 0; j < num; j++ {
			s = append(s, result[j][i])
		}
	}

	for i := 0; i < remainder; i++ {
		line := result[i]
		s = append(s, line[len(line)-1])
	}

	res, err := url.QueryUnescape(string(s))
	if err != nil {
		log.Println(err)
		return "", err
	}
	return strings.Replace(res, "^", "0", -1), nil
}

func (p *xiami) Search(keyword string, page int, limit int) (SearchResult, error) {
	token, err := p.getToken(xiamiApiSearch)
	if err != nil {
		return nil, err
	}

	model := map[string]interface{}{
		"key": keyword,
		"pagingVO": map[string]int{
			"page":     page,
			"pageSize": limit,
		},
	}
	params, err := signPayload(token, model)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", xiamiApiSearch, nil)
	if err != nil {
		return nil, err
	}
	query := req.URL.Query()
	for k, vs := range params {
		query[k] = []string{vs}
	}
	req.URL.RawQuery = query.Encode()

	req.Header.Set("Origin", "https://h.xiami.com")
	req.Header.Set("Referer", "https://h.xiami.com")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:78.0) Gecko/20100101 Firefox/78.0")

	resp, err := p.client.Do(req)
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

	var sr xiamiSearchResult
	if err = json.Unmarshal(content, &sr); err != nil {
		return nil, err
	}

	var songs SearchResult
	for _, s := range sr.Data.Data.Songs {
		listenFiles := s.ListenFiles
		if len(listenFiles) > 0 {
			sort.Sort(sort.Reverse(listenFiles))
			songs = append(songs, Song{
				ID:       s.SongID,
				URL:      listenFiles[0].ListenFile,
				Title:    s.SongName,
				Image:    s.AlbumLogo,
				Artist:   s.ArtistName,
				Provider: "xiami",
			})
		}
	}

	return songs, nil
}

func (p *xiami) SongDetail(song Song) (Song, error) {
	u := fmt.Sprintf(`https://emumo.xiami.com/song/playlist/id/%s/object_name/default/object_id/0/cat/json`, song.ID)

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return song, err
	}

	for _, c := range p.cookies {
		req.AddCookie(c)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:78.0) Gecko/20100101 Firefox/78.0")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "https://www.xiami.com/")
	req.Header.Set("Accept-Language", "zh-CN,zh-HK;q=0.8,zh-TW;q=0.6,en-US;q=0.4,en;q=0.2")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("TE", "Trailers")

	resp, err := p.client.Do(req)
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

	var sd xiamiSongDetail
	err = json.Unmarshal(content, &sd)
	if err != nil {
		return song, err
	}

	u, err = caesar(sd.Data.TrackList[0].Location)
	song.URL = "https:" + u
	return song, err
}

func (p *xiami) HotPlaylist(page int) (Playlists, error) {
	return nil, nil
}

func (p *xiami) PlaylistDetail(pl Playlist) (Songs, error) {
	return nil, nil
}

func (p *xiami) Name() string {
	return "xiami"
}
