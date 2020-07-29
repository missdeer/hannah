package provider

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/missdeer/hannah/util"
)

const (
	kuwoAPISearch     = `http://www.kuwo.cn/api/www/search/searchMusicBykeyWord?key=%s&pn=%d&rn=%d`
	kuwoAPIToken      = `http://www.kuwo.cn/search/list?key=`
	kuwoAPIConvertURL = `http://antiserver.kuwo.cn/anti.s?type=convert_url&format=mp3|wma|aac&response=url&rid=%s`
)

type kuwo struct {
}

type kuwoSearchResult struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Total string `json:"total"`
		List  []struct {
			MusicRID string `json:"musicrid"`
			Artist   string `json:"artist"`
			Pic      string `json:"pic"`
			RID      int    `json:"rid"`
			Album    string `json:"album"`
			AlbumID  int    `json:"albumid"`
			AlbumPic string `json:"albumpic"`
			Pic120   string `json:"pic120"`
			Name     string `json:"name"`
		} `json:"list"`
	} `json:"data"`
}

func (p *kuwo) getToken() (string, error) {
	req, err := http.NewRequest("GET", kuwoAPIToken, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:78.0) Gecko/20100101 Firefox/78.0")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "http://www.kuwo.cn/")
	req.Header.Set("Origin", "http://www.kuwo.cn/")
	req.Header.Set("Accept-Language", "zh-CN,zh-HK;q=0.8,zh-TW;q=0.6,en-US;q=0.4,en;q=0.2")
	req.Header.Set("Accept-Encoding", "gzip, deflate")

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", ErrStatusNotOK
	}

	parsedURL, _ := url.Parse(kuwoAPIToken)
	c := httpClient.Jar.Cookies(parsedURL)
	const kuwoToken = "kw_token"
	for _, cookie := range c {
		if cookie.Name == kuwoToken {
			return cookie.Value, nil
		}
	}
	return "", nil
}

func (p *kuwo) Search(keyword string, page int, limit int) (SearchResult, error) {
	token, err := p.getToken()
	if err != nil {
		return nil, err
	}
	u := fmt.Sprintf(kuwoAPISearch, keyword, page, limit)
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:78.0) Gecko/20100101 Firefox/78.0")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "http://www.kuwo.cn/")
	req.Header.Set("Origin", "http://www.kuwo.cn/")
	req.Header.Set("Accept-Language", "zh-CN,zh-HK;q=0.8,zh-TW;q=0.6,en-US;q=0.4,en;q=0.2")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("csrf", token)

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

	var sr kuwoSearchResult
	err = json.Unmarshal(content, &sr)
	if err != nil {
		return nil, err
	}

	var res SearchResult
	for _, r := range sr.Data.List {
		res = append(res, Song{
			ID:       r.MusicRID,
			Title:    r.Name,
			Image:    r.Pic120,
			Artist:   r.Artist,
			Provider: "kuwo",
		})
	}

	return res, nil
}

func (p *kuwo) ResolveSongURL(song Song) (Song, error) {
	u := fmt.Sprintf(kuwoAPIConvertURL, song.ID)
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return song, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:78.0) Gecko/20100101 Firefox/78.0")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "http://www.kuwo.cn/")
	req.Header.Set("Origin", "http://www.kuwo.cn/")
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
	song.URL = string(content)
	return song, nil
}

func (p *kuwo) ResolveSongLyric(song Song) (Song, error) {
	return song, nil
}

func (p *kuwo) HotPlaylist(page int) (Playlists, error) {
	return nil, nil
}

func (p *kuwo) PlaylistDetail(pl Playlist) (Songs, error) {
	return nil, nil
}

func (p *kuwo) Name() string {
	return "kuwo"
}
