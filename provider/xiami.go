package provider

import (
	"crypto/md5"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/missdeer/hannah/config"
	"github.com/missdeer/hannah/util"
)

const (
	xiamiAppKey          = "23649156"
	xiamiAPISearch       = "https://acs.m.xiami.com/h5/mtop.alimusic.search.searchservice.searchsongs/1.0/?appKey=23649156"
	xiamiBaseURL         = `https://www.xiami.com`
	xiamiAPIHot          = `/api/list/collect`
	xiamiAPIPlaylistInfo = `/api/collect/initialize`
	xiamiAPIAlbumInfo    = `/api/album/initialize`
	xiamiAPIArtistInfo   = `/api/artist/initialize`
	xiamiAPIArtistDetail = `/api/artist/getArtistDetail`
	xiamiAPIArtistSongs  = `/api/song/getArtistSongs`
)

var (
	ErrEmptyTrackList     = errors.New("empty track list")
	ErrTokenNotFound      = errors.New("xiami token not found")
	ErrCodeFieldNotExists = errors.New("code field does not exist")
	reqHeader             = map[string]interface{}{
		"appId":      200,
		"platformId": "h5",
	}
)

type xiami struct {
	token   string
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

func (p *xiami) getToken(u string, tokenKey string) (string, error) {
	parsedURL, _ := url.Parse(u)
	httpClient := util.GetHttpClient()
	c := httpClient.Jar.Cookies(parsedURL)
	if c == nil || len(c) == 0 {
		req, err := http.NewRequest("GET", u, nil)
		if err != nil {
			return "", err
		}

		resp, err := httpClient.Do(req)
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
		if cookie.Name == tokenKey {
			return strings.Split(cookie.Value, "_")[0], nil
		}
	}
	return "", ErrTokenNotFound
}

func signSearchPayload(token string, model interface{}) (map[string]string, error) {
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

func signPlaylistPayload(token string, model interface{}, api string) (string, error) {
	r, err := json.Marshal(model)
	if err != nil {
		return "", err
	}
	origin := fmt.Sprintf(`%s_xmMain_%s_%s`, strings.Split(token, "_")[0], api, string(r))
	sign := fmt.Sprintf("%x", md5.Sum([]byte(origin)))
	return fmt.Sprintf(`https://www.xiami.com%s?_q=%s&_s=%s`, api, url.QueryEscape(string(r)), sign), nil
}

func caesar(location string) (string, error) {
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

func (p *xiami) SearchSongs(keyword string, page int, limit int) (SearchResult, error) {
	token, err := p.getToken(xiamiAPISearch, "_m_h5_tk")
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
	params, err := signSearchPayload(token, model)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", xiamiAPISearch, nil)
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
	req.Header.Set("User-Agent", config.UserAgent)

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

func (p *xiami) ResolveSongURL(song Song) (Song, error) {
	u := fmt.Sprintf(`https://emumo.xiami.com/song/playlist/id/%s/object_name/default/object_id/0/cat/json`, song.ID)

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return song, err
	}

	req.Header.Set("User-Agent", config.UserAgent)
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "https://www.xiami.com/")
	req.Header.Set("Accept-Language", "zh-CN,zh-HK;q=0.8,zh-TW;q=0.6,en-US;q=0.4,en;q=0.2")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("TE", "Trailers")

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

	var sd xiamiSongDetail
	err = json.Unmarshal(content, &sd)
	if err != nil {
		return song, err
	}

	if len(sd.Data.TrackList) == 0 {
		return song, ErrEmptyTrackList
	}
	u, err = caesar(sd.Data.TrackList[0].Location)
	song.URL = "https:" + u
	return song, err
}

func (p *xiami) ResolveSongLyric(song Song) (Song, error) {
	if strings.HasPrefix(song.Lyric, "https://") || strings.HasPrefix(song.Lyric, "http://") {
		req, err := http.NewRequest("GET", song.Lyric, nil)
		if err != nil {
			return song, err
		}

		req.Header.Set("User-Agent", config.UserAgent)
		req.Header.Set("Accept", "application/json, text/plain, */*")
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Referer", "https://www.xiami.com/")
		req.Header.Set("Accept-Language", "zh-CN,zh-HK;q=0.8,zh-TW;q=0.6,en-US;q=0.4,en;q=0.2")
		req.Header.Set("Accept-Encoding", "gzip, deflate")
		req.Header.Set("TE", "Trailers")

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

		song.Lyric = string(content)
	}
	return song, nil
}

type xiamiHot struct {
	Code   string `json:"code"`
	Result struct {
		Status string `json:"status"`
		Data   struct {
			Collects []struct {
				ListID      int    `json:"listId"`
				CollectLogo string `json:"collectLogo"`
				CollectName string `json:"collectName"`
			} `json:"collects"`
			PagingVO struct {
				Page     int `json:"page"`
				PageSize int `json:"pageSize"`
				Pages    int `json:"pages"`
				Count    int `json:"count"`
			}
		} `json:"data"`
	} `json:"result"`
}

func (p *xiami) HotPlaylist(page int, limit int) (Playlists, error) {
	token, err := p.getToken(xiamiBaseURL+xiamiAPIHot, `xm_sg_tk`)
	if err != nil {
		return nil, err
	}

	model := map[string]interface{}{
		"pagingVO": map[string]int{
			"page":     page,
			"pageSize": limit,
		},
		"dataType": "system",
	}
	u, err := signPlaylistPayload(token, model, xiamiAPIHot)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Origin", "https://h.xiami.com")
	req.Header.Set("Referer", "https://h.xiami.com")
	req.Header.Set("User-Agent", config.UserAgent)

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

	var hot xiamiHot
	if err = json.Unmarshal(content, &hot); err != nil {
		return nil, err
	}

	var pls Playlists
	for _, pl := range hot.Result.Data.Collects {
		pls = append(pls, Playlist{
			ID:       strconv.Itoa(pl.ListID),
			Title:    pl.CollectName,
			Image:    pl.CollectLogo,
			Provider: "xiami",
		})
	}

	return pls, nil
}

type xiamiPlaylistDetail struct {
	Code   string `json:"code"`
	Result struct {
		Status string `json:"status"`
		Data   struct {
			CollectSongs []struct {
				SongID        int    `json:"songId"`
				SongStringID  string `json:"songStringId"`
				SongName      string `json:"songName"`
				AlbumID       int    `json:"albumId"`
				AlbumStringID string `json:"albumStringId"`
				AlbumLogo     string `json:"albumLogo"`
				AlbumName     string `json:"albumName"`
				ArtistID      int    `json:"artistId"`
				ArtistName    string `json:"artistName"`
				ArtistLogo    string `json:"artistLogo"`
				Singers       string `json:"singers"`
			} `json:"collectSongs"`
		} `json:"data"`
	} `json:"result"`
}

func (p *xiami) PlaylistDetail(pl Playlist) (Songs, error) {
	token, err := p.getToken(xiamiBaseURL+xiamiAPIPlaylistInfo, `xm_sg_tk`)
	if err != nil {
		return nil, err
	}

	model := map[string]interface{}{
		"listId": pl.ID,
	}
	u, err := signPlaylistPayload(token, model, xiamiAPIPlaylistInfo)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Origin", "https://h.xiami.com")
	req.Header.Set("Referer", "https://h.xiami.com")
	req.Header.Set("User-Agent", config.UserAgent)

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

	var pld xiamiPlaylistDetail
	if err = json.Unmarshal(content, &pld); err != nil {
		return nil, err
	}

	var songs Songs
	for _, pl := range pld.Result.Data.CollectSongs {
		songs = append(songs, Song{
			ID:       strconv.Itoa(pl.SongID),
			Title:    pl.SongName,
			Artist:   pl.Singers,
			Image:    pl.AlbumLogo,
			Provider: "xiami",
		})
	}

	return songs, nil
}

type xiamiArtistInfo struct {
	Code   string `json:"code"`
	Result struct {
		Status string `json:"status"`
		Data   struct {
			ArtistDetail struct {
				ArtistID       int    `json:"artistId"`
				ArtistStringID string `json:"artistStringId"`
				ArtistName     string `json:"artistName"`
				ArtistLogo     string `json:"artistLogo"`
			} `json:"artistDetail"`
		} `json:"data"`
	} `json:"result"`
}

func (p *xiami) getArtistIDFromArtistStringID(id string) (int, error) {
	token, err := p.getToken(xiamiBaseURL+xiamiAPIArtistInfo, `xm_sg_tk`)
	if err != nil {
		return 0, err
	}

	model := map[string]interface{}{
		"artistId": id,
	}
	u, err := signPlaylistPayload(token, model, xiamiAPIArtistInfo)
	if err != nil {
		return 0, err
	}
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return 0, err
	}

	req.Header.Set("Origin", "https://h.xiami.com")
	req.Header.Set("Referer", "https://h.xiami.com")
	req.Header.Set("User-Agent", config.UserAgent)

	httpClient := util.GetHttpClient()
	resp, err := httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return 0, ErrStatusNotOK
	}

	content, err := util.ReadHttpResponseBody(resp)
	if err != nil {
		return 0, err
	}

	var pld xiamiArtistInfo
	if err = json.Unmarshal(content, &pld); err != nil {
		return 0, err
	}

	return pld.Result.Data.ArtistDetail.ArtistID, nil
}

type xiamiArtistSongs struct {
	Code   string `json:"code"`
	Result struct {
		Status string `json:"status"`
		Data   struct {
			Songs []struct {
				SongID       int    `json:"songId"`
				SongStringId string `json:"songStringId"`
				SongName     string `json:"songName"`
				AlbumLogo    string `json:"albumLogo"`
				ArtistName   string `json:"artistName"`
				Singers      string `json:"singers"`
				LyricInfo    struct {
					LyricFile string `json:"lyricFile"`
				} `json:"lyricInfo"`
				ArtistVOs []struct {
					ArtistID   int    `json:"artistId"`
					ArtistName string `json:"artistName"`
				} `json:"artistVOs"`
				SingerVOs []struct {
					ArtistID   int    `json:"artistId"`
					ArtistName string `json:"artistName"`
				} `json:"singerVOs"`
			} `json:"songs"`
			Total int `json:"total"`
		} `json:"data"`
	} `json:"result"`
}

func (p *xiami) ArtistSongs(id string) (res Songs, err error) {
	reg := regexp.MustCompile(`^[0-9]+$`)
	if !reg.MatchString(id) {
		idNr, err := p.getArtistIDFromArtistStringID(id)
		if err != nil {
			return nil, err
		}
		id = strconv.Itoa(idNr)
	}
	token, err := p.getToken(xiamiBaseURL+xiamiAPIArtistSongs, `xm_sg_tk`)
	if err != nil {
		return nil, err
	}

	model := map[string]interface{}{
		"artistId": id,
		"category": 0,
		"pagingVO": map[string]int{
			"page":     config.Page,
			"pageSize": config.Limit,
		},
	}
	u, err := signPlaylistPayload(token, model, xiamiAPIArtistSongs)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Origin", "https://h.xiami.com")
	req.Header.Set("Referer", "https://h.xiami.com")
	req.Header.Set("User-Agent", config.UserAgent)

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

	var pld xiamiArtistSongs
	if err = json.Unmarshal(content, &pld); err != nil {
		return nil, err
	}

	var songs Songs
	for _, pl := range pld.Result.Data.Songs {
		var artists []string
		for _, a := range pl.SingerVOs {
			artists = append(artists, a.ArtistName)
		}
		songs = append(songs, Song{
			ID:       strconv.Itoa(pl.SongID),
			Title:    pl.SongName,
			Artist:   strings.Join(artists, "/"),
			Image:    pl.AlbumLogo,
			Lyric:    pl.LyricInfo.LyricFile,
			Provider: "xiami",
		})
	}

	return songs, nil
}

type xiamiAlbumSongs struct {
	Code   string `json:"code"`
	Result struct {
		Status string `json:"status"`
		Data   struct {
			AlbumDetail struct {
				AlbumLogo  string `json:"albumLogo"`
				ArtistName string `json:"artistName"`
				Songs      []struct {
					SongID       int    `json:"songId"`
					SongStringId string `json:"songStringId"`
					SongName     string `json:"songName"`
					LyricInfo    struct {
						LyricFile string `json:"lyricFile"`
					} `json:"lyricInfo"`
				} `json:"songs"`
				Artists []struct {
					ArtistName string `json:"artistName"`
					ArtistLogo string `json:"artistLogo"`
				} `json:"artists"`
			} `json:"albumdetail"`
		} `json:"data"`
	} `json:"result"`
}

func (p *xiami) AlbumSongs(id string) (res Songs, err error) {
	token, err := p.getToken(xiamiBaseURL+xiamiAPIAlbumInfo, `xm_sg_tk`)
	if err != nil {
		return nil, err
	}

	model := map[string]interface{}{
		"albumId": id,
	}
	u, err := signPlaylistPayload(token, model, xiamiAPIAlbumInfo)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Origin", "https://h.xiami.com")
	req.Header.Set("Referer", "https://h.xiami.com")
	req.Header.Set("User-Agent", config.UserAgent)

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

	var pld xiamiAlbumSongs
	if err = json.Unmarshal(content, &pld); err != nil {
		return nil, err
	}

	var artists []string
	for _, a := range pld.Result.Data.AlbumDetail.Artists {
		artists = append(artists, a.ArtistName)
	}
	var songs Songs
	for _, pl := range pld.Result.Data.AlbumDetail.Songs {
		songs = append(songs, Song{
			ID:       strconv.Itoa(pl.SongID),
			Title:    pl.SongName,
			Artist:   strings.Join(artists, "/"),
			Image:    pld.Result.Data.AlbumDetail.AlbumLogo,
			Lyric:    pl.LyricInfo.LyricFile,
			Provider: "xiami",
		})
	}

	return songs, nil
}

func (p *xiami) Login() error {
	return ErrNotImplemented
}

func (p *xiami) Name() string {
	return "xiami"
}
