package provider

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/missdeer/hannah/util"
)

type xiami struct {
	cookies []*http.Cookie
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

type xiamiSearchResult struct {
	Code   string `json:"code"`
	Result struct {
		Status string `json:"status"`
		Data   struct {
			PagingVO struct {
				Page     int `json:"page"`
				PageSize int `json:"pageSize"`
				Pages    int `json:"pages"`
				Count    int `json:"count"`
			} `json:"pagingVO"`
			Songs []struct {
				SongID        int    `json:"songId"`
				SongStringID  string `json:"songStringId"`
				SongName      string `json:"songName"`
				AlbumID       int    `json:"albumId"`
				AlbumStringID string `json:"albumStringId"`
				AlbumLogo     string `json:"albumLogo"`
				AlbumLogoS    string `json:"albumLogoS"`
				AlbumName     string `json:"albumName"`
				AlbumSubName  string `json:"albumSubName"`
				ArtistID      int    `json:"artistId"`
				ArtistName    string `json:"artistName"`
				ArtistLogo    string `json:"artistLogo"`
				Singers       string `json:"singers"`
				PinYin        string `json:"pinyin"`
			} `json:"songs"`
		} `json:"data"`
	} `json:"result"`
}

func caesar(location string) (string, error) {
	// caesar(location)
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
start:
	// curl 'https://www.xiami.com/api/search/searchSongs?_q=%7B%22pagingVO%22:%7B%22page%22:%221%22,%22pageSize%22:60%7D,%22key%22:%22jay%22%7D&_s=366bc6054c0f94e3561642d06e651017'
	// -H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:79.0) Gecko/20100101 Firefox/79.0'
	// -H 'Accept: application/json, text/plain, */*'
	// -H 'Accept-Language: zh-CN,zh-HK;q=0.8,zh-TW;q=0.6,en-US;q=0.4,en;q=0.2'
	// --compressed
	// -H 'Connection: keep-alive'
	// -H 'Cookie: gid=151166315123687; xmgid=723ac295-e03a-4109-9c53-22b49371aea8; _uab_collina=157266137909801636000205; _xm_cf_=qltiVwysrtKMw3W0p_Z0fQ-U; xm_sg_tk=628da31835e19e3d3a65c49ea6a0f9f9_1595162817446; xm_sg_tk.sig=m58QVuGn8hcLNcuaWO3vHOlZXZaC-Mjp0O1oJzy1gG4'
	// -H 'Referer: https://www.xiami.com/search?key=jay'
	// -H 'Pragma: no-cache'
	// -H 'Cache-Control: no-cache'

	// curl 'https://www.xiami.com/api/search/searchSongs?_q=%7B%22pagingVO%22:%7B%22page%22:%221%22,%22pageSize%22:60%7D,%22key%22:%22jay%22%7D&_s=e5c324f991fea3ce9f904f49505d6499'
	// -H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:79.0) Gecko/20100101 Firefox/79.0'
	// -H 'Accept: application/json, text/plain, */*'
	// -H 'Accept-Language: zh-CN,zh-HK;q=0.8,zh-TW;q=0.6,en-US;q=0.4,en;q=0.2'
	// --compressed
	// -H 'Connection: keep-alive'
	// -H 'Cookie: gid=151166315123687; xmgid=723ac295-e03a-4109-9c53-22b49371aea8; _uab_collina=157266137909801636000205; _xm_cf_=qltiVwysrtKMw3W0p_Z0fQ-U; xm_sg_tk=e787be0b681f0ef5339f53186085eca3_1595253790619; xm_sg_tk.sig=DmgnD-KBRRnABMEpKGTiA61CEB-3qDXzDMVpjRa9Yhc'
	// -H 'Referer: https://www.xiami.com/search?key=jay'
	// -H 'Pragma: no-cache' -H 'Cache-Control: no-cache'
	// -H 'TE: Trailers'
	v := url.Values{}
	v.Add("_q", fmt.Sprintf(`{"pagingVO":{"page":"%d", "pageSize":"%d"},"key":"%s"}`, page, limit, keyword))
	v.Add("_s", strconv.FormatInt(time.Now().UnixNano(), 10))
	u := "https://www.xiami.com/api/search/searchSongs?" + v.Encode()
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	for _, c := range p.cookies {
		req.AddCookie(c)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:78.0) Gecko/20100101 Firefox/78.0")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", fmt.Sprintf("https://www.xiami.com/search?key=%s", url.QueryEscape(keyword)))
	req.Header.Set("Accept-Language", "zh-CN,zh-HK;q=0.8,zh-TW;q=0.6,en-US;q=0.4,en;q=0.2")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("TE", "Trailers")

	client := util.GetHttpClient()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("status != 200")
	}

	content, err := util.ReadHttpResponseBody(resp)
	if err != nil {
		return nil, err
	}

	var simpleResp map[string]interface{}
	if err = json.Unmarshal(content, &simpleResp); err != nil {
		return nil, err
	}

	code, ok := simpleResp["code"]
	if !ok {
		return nil, errors.New("code field does not exist")
	}
	codeStr, ok := code.(string)
	if !ok {
		return nil, err
	}
	switch codeStr {
	case "SG_TOKEN_EXPIRED", "SG_TOKEN_EMPTY", "SG_INVALID":
		// extract cookies
		p.cookies = resp.Cookies()
		goto start
	case "SUCCESS":
	default:
		return nil, errors.New(codeStr)
	}

	p.cookies = resp.Cookies()

	var sr xiamiSearchResult
	if err = json.Unmarshal(content, &sr); err != nil {
		return nil, err
	}

	var songs SearchResult
	for _, s := range sr.Result.Data.Songs {
		songs = append(songs, Song{
			ID:       strconv.Itoa(s.SongID),
			Title:    s.SongName,
			Image:    s.AlbumLogo,
			Artist:   s.ArtistName,
			Provider: "xiami",
		})
	}

	return songs, nil
}

func (p *xiami) SongURL(song Song) (string, error) {
	u := fmt.Sprintf(`https://emumo.xiami.com/song/playlist/id/%s/object_name/default/object_id/0/cat/json`, song.ID)

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return "", err
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

	client := util.GetHttpClient()

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", errors.New("status != 200")
	}

	content, err := util.ReadHttpResponseBody(resp)
	if err != nil {
		return "", err
	}

	var sd xiamiSongDetail
	err = json.Unmarshal(content, &sd)
	if err != nil {
		return "", err
	}

	return caesar(sd.Data.TrackList[0].Location)
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
