package provider

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/missdeer/hannah/util"
)

var (
	ErrEmptyMidURLInfoField = errors.New("empty MidURLInfo field")
	ErrEmptyPURL            = errors.New("empty PURL, may be VIP needed")
)

type qq struct {
}

type qqSearchResult struct {
	Code int `json:"code"`
	Data struct {
		Keyword string `json:"keyword"`
		Song    struct {
			CurNum   int `json:"curnum"`
			CurPage  int `json:"curpage"`
			TotalNum int `json:"totalnum"`
			List     []struct {
				AlbumID   int    `json:"albumid"`
				AlbumMID  string `json:"albummid"`
				AlbumName string `json:"albumname"`
				SongID    int    `json:"songid"`
				SongMID   string `json:"songmid"`
				SongName  string `json:"songname"`
				Singer    []struct {
					ID   int    `json:"id"`
					MID  string `json:"mid"`
					Name string `json:"name"`
				} `json:"singer"`
			} `json:"list"`
		}
	} `json:"data"`
}

func (p *qq) Search(keyword string, page int, limit int) (SearchResult, error) {
	// http://i.y.qq.com/s.music/fcgi-bin/search_for_qq_cp?g_tk=938407465&uin=0&format=jsonp&inCharset=utf-8&outCharset=utf-8&notice=0&platform=h5&needNewCode=1&w=陈奕迅&
	// zhidaqu=1&catZhida=1&t=0&flag=1&ie=utf-8&sem=1&aggr=0&perpage=20&n=20&p=1&remoteplace=txt.mqq.all&_=1459991037831&jsonpCallback=jsonp4
	v := url.Values{}
	v.Add("g_tk", "938407465")
	v.Add("uin", "0")
	v.Add("format", "json")
	v.Add("inCharset", "utf-8")
	v.Add("outCharset", "utf-8")
	v.Add("notice", "0")
	v.Add("platform", "h5")
	v.Add("needNewCode", "1")
	v.Add("w", keyword)
	v.Add("zhidaqu", "1")
	v.Add("catZhida", "1")
	v.Add("t", "0")
	v.Add("flag", "1")
	v.Add("ie", "utf-8")
	v.Add("sem", "1")
	v.Add("aggr", "0")
	v.Add("perpage", strconv.Itoa(limit))
	v.Add("n", strconv.Itoa(limit))
	v.Add("p", strconv.Itoa(page))
	v.Add("remoteplace", "txt.mqq.all")
	v.Add("_", strconv.FormatInt(time.Now().UnixNano(), 10))
	u := "http://i.y.qq.com/s.music/fcgi-bin/search_for_qq_cp?" + v.Encode()
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:78.0) Gecko/20100101 Firefox/78.0")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "http://y.qq.com/")
	req.Header.Set("Origin", "http://y.qq.com/")
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

	var sr qqSearchResult
	err = json.Unmarshal(content, &sr)
	if err != nil {
		return nil, err
	}

	if sr.Code != 0 {
		return nil, fmt.Errorf("code = %d", sr.Code)
	}

	var songs SearchResult
	for _, s := range sr.Data.Song.List {
		var artists []string
		for _, a := range s.Singer {
			artists = append(artists, a.Name)
		}
		songs = append(songs, Song{
			ID:       s.SongMID,
			Title:    s.SongName,
			Artist:   strings.Join(artists, "/"),
			Provider: "qq",
		})
	}

	return songs, nil
}

type qqSongDetail struct {
	Code int `json:"code"`
	Req0 struct {
		Code int `json:"code"`
		Data struct {
			MidURLInfo []struct {
				FileName string `json:"filename"`
				PURL     string `json:"purl"`
			} `json:"midurlinfo"`
		} `json:"data"`
	} `json:"req_0"`
}

func (p *qq) ResolveSongURL(song Song) (Song, error) {
	// https://u.y.qq.com/cgi-bin/musicu.fcg?loginUin=0&hostUin=0&format=json&inCharset=utf8&outCharset=utf-8&notice=0&platform=yqq.json&needNewCode=0
	// &data=%7B%22req_0%22%3A%7B%22module%22%3A%22vkey.GetVkeyServer%22%2C%22method%22%3A%22CgiGetVkey%22%2C%22param%22%3A%7B%22guid%22%3A%2210000%22%2C%22songmid%22%3A%5B%22003VQrF72a0DGb%22%5D%2C%22songtype%22%3A%5B0%5D%2C%22uin%22%3A%220%22%2C%22loginflag%22%3A1%2C%22platform%22%3A%2220%22%7D%7D%2C%22comm%22%3A%7B%22uin%22%3A0%2C%22format%22%3A%22json%22%2C%22ct%22%3A20%2C%22cv%22%3A0%7D%7D
	u := `https://u.y.qq.com/cgi-bin/musicu.fcg?loginUin=0&hostUin=0&format=json&inCharset=utf8&outCharset=utf-8&notice=0&platform=yqq.json&needNewCode=0&data=%7B%22req_0%22%3A%7B%22module%22%3A%22vkey.GetVkeyServer%22%2C%22method%22%3A%22CgiGetVkey%22%2C%22param%22%3A%7B%22guid%22%3A%2210000%22%2C%22songmid%22%3A%5B%22` + song.ID + `%22%5D%2C%22songtype%22%3A%5B0%5D%2C%22uin%22%3A%220%22%2C%22loginflag%22%3A1%2C%22platform%22%3A%2220%22%7D%7D%2C%22comm%22%3A%7B%22uin%22%3A0%2C%22format%22%3A%22json%22%2C%22ct%22%3A20%2C%22cv%22%3A0%7D%7D`

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return song, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:78.0) Gecko/20100101 Firefox/78.0")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "http://y.qq.com/")
	req.Header.Set("Origin", "http://y.qq.com/")
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

	var detail qqSongDetail
	err = json.Unmarshal(content, &detail)
	if err != nil {
		return song, err
	}

	if detail.Code != 0 {
		return song, fmt.Errorf("code = %d", detail.Code)
	}

	if len(detail.Req0.Data.MidURLInfo) == 0 {
		return song, ErrEmptyMidURLInfoField
	}

	if detail.Req0.Data.MidURLInfo[0].PURL == "" {
		return song, ErrEmptyPURL
	}

	song.URL = `http://ws.stream.qqmusic.qq.com/` + detail.Req0.Data.MidURLInfo[0].PURL
	return song, nil
}

func (p *qq) ResolveSongLyric(song Song) (Song, error) {
	return song, nil
}

type qqHot struct {
	Code int `json:"code"`
	Data struct {
		SIN  int `json:"sin"`
		EIN  int `json:"ein"`
		List []struct {
			DissID   string `json:"dissid"`
			DissName string `json:"dissname"`
			ImgURL   string `json:"imgurl`
		} `json:"list"`
	} `json:"data"`
}

func (p *qq) HotPlaylist(page int, limit int) (Playlists, error) {
	// https://c.y.qq.com/splcloud/fcgi-bin/fcg_get_diss_by_tag.fcg?picmid=1&rnd=%d&g_tk=732560869&loginUin=0&hostUin=0&format=json&inCharset=utf8&outCharset=utf-8&notice=0&platform=yqq.json&needNewCode=0&categoryId=10000000&sortId=5&sin=%d&ein=%d
	u := fmt.Sprintf(`https://c.y.qq.com/splcloud/fcgi-bin/fcg_get_diss_by_tag.fcg?picmid=1&rnd=%d&g_tk=732560869&loginUin=0&hostUin=0&format=json&inCharset=utf8&outCharset=utf-8&notice=0&platform=yqq.json&needNewCode=0&categoryId=10000000&sortId=5&sin=%d&ein=%d`,
		rand.Int63(), page*limit, (page+1)*limit-1)
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:78.0) Gecko/20100101 Firefox/78.0")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "http://y.qq.com/")
	req.Header.Set("Origin", "http://y.qq.com/")
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

	var hot qqHot
	if err = json.Unmarshal(content, &hot); err != nil {
		return nil, err
	}

	var pls Playlists
	for _, pl := range hot.Data.List {
		pls = append(pls, Playlist{
			ID:       pl.DissID,
			Title:    pl.DissName,
			Image:    pl.ImgURL,
			Provider: "qq",
		})
	}

	return pls, nil
}

type qqPlaylistDetail struct {
	Code   int `json:"code"`
	CDList []struct {
		DissTID  string `json:"disstid"`
		DissID   int    `json:"dissid"`
		DissName string `json:"DissName"`
		Logo     string `json:"logo"`
		SongList []struct {
			AlbumID   int    `json:"albumid"`
			AlbumName string `json:"albumname"`
			Singer    []struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			} `json:"singer"`
			SongID   int    `json:"songid"`
			SongMID  string `json:"songmid"`
			SongName string `json:"songname"`
		} `json:"songlist"`
	} `json:"cdlist"`
}

func (p *qq) PlaylistDetail(pl Playlist) (Songs, error) {
	// http://i.y.qq.com/qzone-music/fcg-bin/fcg_ucc_getcdinfo_byids_cp.fcg?type=1&json=1&utf8=1&onlysong=0&nosign=1&disstid=%s&g_tk=5381&loginUin=0&hostUin=0&format=json&inCharset=GB2312&outCharset=utf-8&notice=0&platform=yqq&jsonpCallback=jsonCallback&needNewCode=0
	u := fmt.Sprintf(`http://i.y.qq.com/qzone-music/fcg-bin/fcg_ucc_getcdinfo_byids_cp.fcg?type=1&json=1&utf8=1&onlysong=0&nosign=1&disstid=%s&g_tk=5381&loginUin=0&hostUin=0&format=json&inCharset=GB2312&outCharset=utf-8&notice=0&platform=yqq&jsonpCallback=jsonCallback&needNewCode=0`,
		pl.ID)
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:78.0) Gecko/20100101 Firefox/78.0")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "http://y.qq.com/")
	req.Header.Set("Origin", "http://y.qq.com/")
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

	if bytes.HasPrefix(content, []byte(`jsonCallback(`)) && bytes.HasSuffix(content, []byte(`)`)) {
		content = content[len(`jsonCallback(`):]
		content = content[:len(content)-1]
	}

	var pld qqPlaylistDetail
	if err = json.Unmarshal(content, &pld); err != nil {
		return nil, err
	}

	if len(pld.CDList) == 0 {
		return nil, errors.New("empty playlist")
	}

	var songs Songs
	for _, pl := range pld.CDList[0].SongList {
		var singers []string
		for _, s := range pl.Singer {
			singers = append(singers, s.Name)
		}
		songs = append(songs, Song{
			ID:       pl.SongMID,
			Title:    pl.SongName,
			Image:    pld.CDList[0].Logo,
			Artist:   strings.Join(singers, "/"),
			Provider: "qq",
		})
	}

	return songs, nil
}

func (p *qq) ArtistSongs(id string) (res Songs, err error) {
	return nil, ErrNotImplemented
}

type qqAlbumSongs struct {
	Code int `json:"code"`
	Data struct {
		SingerName   string `json:"singername"`
		SingerMID    string `json:"singermid"`
		ID           int    `json:"id"`
		Total        int    `json:"total"`
		TotalSongNum int    `json:"total_song_num"`
		MID          string `json:"mid"`
		List         []struct {
			SongName  string `json:"songname"`
			SongMID   string `json:"songmid"`
			AlbumMID  string `json:"albummid"`
			AlbumName string `json:"albumname"`
			Singer    []struct {
				MID  string `json:"mid"`
				Name string `json:"name"`
			} `json:"singer"`
		} `json:"list"`
	} `json:"data"`
}

func (p *qq) AlbumSongs(id string) (res Songs, err error) {
	// https://c.y.qq.com/v8/fcg-bin/fcg_v8_album_info_cp.fcg?albummid=001IskfD3Vncxo&g_tk=1278911659&hostUin=0&format=jsonp&jsonpCallback=callback&inCharset=utf8&outCharset=utf-8&notice=0&platform=yqq&needNewCode=0
	u := fmt.Sprintf(`https://c.y.qq.com/v8/fcg-bin/fcg_v8_album_info_cp.fcg?albummid=%s&g_tk=1278911659&hostUin=0&format=jsonp&jsonpCallback=callback&inCharset=utf8&outCharset=utf-8&notice=0&platform=yqq&needNewCode=0`, id)
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:78.0) Gecko/20100101 Firefox/78.0")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "http://y.qq.com/")
	req.Header.Set("Origin", "http://y.qq.com/")
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
	content = bytes.TrimSpace(content)
	if bytes.HasPrefix(content, []byte(`callback(`)) && bytes.HasSuffix(content, []byte(`})`)) {
		content = content[len(`callback(`):]
		content = content[:len(content)-1]
	}

	var pld qqAlbumSongs
	if err = json.Unmarshal(content, &pld); err != nil {
		return nil, err
	}

	if len(pld.Data.List) == 0 {
		return nil, errors.New("empty playlist")
	}

	var songs Songs
	for _, pl := range pld.Data.List {
		var singers []string
		for _, s := range pl.Singer {
			singers = append(singers, s.Name)
		}
		songs = append(songs, Song{
			ID:       pl.SongMID,
			Title:    pl.SongName,
			Artist:   strings.Join(singers, "/"),
			Provider: "qq",
		})
	}

	return songs, nil
}

func (p *qq) Name() string {
	return "qq"
}
