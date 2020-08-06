package provider

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/missdeer/hannah/util"
	"github.com/missdeer/hannah/util/cryptography"
)

const (
	Base62                             = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	neteasePresetKey                   = "0CoJUm6Qyw8W8jud"
	neteaseIV                          = "0102030405060708"
	neteaseLinuxAPIKey                 = "rFgB&h#%2?^eDg:Q"
	neteaseEAPIKey                     = "e82ckenh8dichen8"
	neteaseDefaultRSAPublicKeyModulus  = "e0b509f6259df8642dbc35662901477df22677ec152b5ff68ace615bb7b725152b3ab17a876aea8a5aa76d2e417629ec4ee341f56135fccf695280104e0312ecbda92557c93870114af6c9d05c4f7f0c3685b7a46bee255932575cce10b424d813cfe4875d3e82047b97ddef52741d546b8e289dc6935b3ece0462db0a22b8e7"
	neteaseDefaultRSAPublicKeyExponent = 0x10001
	neteaseAPIGetSongsURL              = "http://music.163.com/weapi/song/enhance/player/url/v1?csrf_token="
	neteaseAPISearch                   = `http://music.163.com/api/search/pc`
	neteaseAPIGetLyric                 = `http://music.163.com/weapi/song/lyric?csrf_token=`
	neteaseAPIHot                      = `http://music.163.com/discover/playlist/?order=hot&limit=%d&offset=%d`
	neteaseAPIPlaylistDetail           = `http://music.163.com/weapi/v3/playlist/detail`
	neteaseAPISongDetail               = `http://music.163.com/weapi/v3/song/detail`
)

func weapi(origData interface{}) map[string]interface{} {
	plainText, _ := json.Marshal(origData)
	params := base64.StdEncoding.EncodeToString(cryptography.AESCBCEncrypt(plainText, []byte(neteasePresetKey), []byte(neteaseIV)))
	secKey := createSecretKey(16, Base62)
	params = base64.StdEncoding.EncodeToString(cryptography.AESCBCEncrypt([]byte(params), secKey, []byte(neteaseIV)))
	return map[string]interface{}{
		"params":    params,
		"encSecKey": cryptography.RSAEncrypt(bytesReverse(secKey), neteaseDefaultRSAPublicKeyModulus, neteaseDefaultRSAPublicKeyExponent),
	}
}

func linuxapi(origData interface{}) map[string]interface{} {
	plainText, _ := json.Marshal(origData)
	return map[string]interface{}{
		"eparams": strings.ToUpper(hex.EncodeToString(cryptography.AESECBEncrypt(plainText, []byte(neteaseLinuxAPIKey)))),
	}
}

func eapi(url string, origData interface{}) map[string]interface{} {
	plainText, _ := json.Marshal(origData)
	text := string(plainText)
	message := fmt.Sprintf("nobody%suse%smd5forencrypt", url, text)
	digest := fmt.Sprintf("%x", md5.Sum([]byte(message)))
	data := fmt.Sprintf("%s-36cd479b6b5-%s-36cd479b6b5-%s", url, text, digest)
	return map[string]interface{}{
		"params": strings.ToUpper(hex.EncodeToString(cryptography.AESECBEncrypt([]byte(data), []byte(neteaseEAPIKey)))),
	}
}

func createSecretKey(size int, charset string) []byte {
	secKey, n := make([]byte, size), len(charset)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range secKey {
		secKey[i] = charset[r.Intn(n)]
	}
	return secKey
}

func bytesReverse(b []byte) []byte {
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return b
}

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

type neteaseSongInfo struct {
	Code int `json:"code"`
	Data []struct {
		ID         int    `json:"id"`
		URL        string `json:"url"`
		BR         int    `json:"br"`
		Size       int    `json:"size"`
		Type       string `json:"type"`
		Level      string `json:"level"`
		EncodeType string `json:"encodeType"`
	} `json:"data"`
}

func (p *netease) Search(keyword string, page int, limit int) (SearchResult, error) {
	body := fmt.Sprintf("limit=%d&offset=%d&s=%s&type=1", limit, (page-1)*limit, url.QueryEscape(keyword))
	req, err := http.NewRequest("POST", neteaseAPISearch, strings.NewReader(body))
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

func (p *netease) ResolveSongURL(song Song) (Song, error) {
	data := map[string]interface{}{
		"ids":        fmt.Sprintf("[%s]", song.ID),
		"level":      "standard",
		"encodeType": "aac",
		"csrf_token": "",
	}

	params := weapi(data)
	values := url.Values{}
	for k, vs := range params {
		values.Add(k, vs.(string))
	}
	postBody := values.Encode()
	req, err := http.NewRequest("POST", neteaseAPIGetSongsURL, strings.NewReader(postBody))
	if err != nil {
		return song, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:78.0) Gecko/20100101 Firefox/78.0")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "http://music.163.com/")
	req.Header.Set("Origin", "http://music.163.com/")
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

	var songInfo neteaseSongInfo
	if err = json.Unmarshal(content, &songInfo); err != nil {
		return song, err
	}

	if len(songInfo.Data) == 0 || songInfo.Data[0].URL == "" {
		return song, err
	}

	song.URL = songInfo.Data[0].URL
	return song, nil
}

type neteaseLyricDetail struct {
	SGC bool `json:"sgc"`
	SFY bool `json:"sfy"`
	QFY bool `json:"qfy"`
	LRC struct {
		Version int    `json:"version"`
		Lyric   string `json:"lyric"`
	} `json:"lrc"`
	Code int `json:"code"`
}

func (p *netease) ResolveSongLyric(song Song) (Song, error) {
	data := map[string]interface{}{
		"id":         song.ID,
		"lv":         -1,
		"tv":         -1,
		"csrf_token": "",
	}

	params := weapi(data)
	values := url.Values{}
	for k, vs := range params {
		values.Add(k, vs.(string))
	}
	postBody := values.Encode()
	req, err := http.NewRequest("POST", neteaseAPIGetLyric, strings.NewReader(postBody))
	if err != nil {
		return song, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:78.0) Gecko/20100101 Firefox/78.0")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "http://music.163.com/")
	req.Header.Set("Origin", "http://music.163.com/")
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

	var lyric neteaseLyricDetail
	if err = json.Unmarshal(content, &lyric); err != nil {
		return song, nil
	}

	song.Lyric = lyric.LRC.Lyric
	return song, nil
}

func (p *netease) HotPlaylist(page int, limit int) (res Playlists, err error) {
	u := fmt.Sprintf(neteaseAPIHot, limit, (page-1)*limit)
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:78.0) Gecko/20100101 Firefox/78.0")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "http://music.163.com/")
	req.Header.Set("Origin", "http://music.163.com/")
	req.Header.Set("Accept-Language", "zh-CN,zh-HK;q=0.8,zh-TW;q=0.6,en-US;q=0.4,en;q=0.2")
	req.Header.Set("Accept-Encoding", "gzip, deflate")

	resp, err := httpClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return res, ErrStatusNotOK
	}

	content, err := util.ReadHttpResponseBody(resp)
	if err != nil {
		return
	}

	regPlaylistInfo := regexp.MustCompile(`^<a\stitle="([^"]+)"\shref="\/playlist\?id=(\d+)"\sclass="msk"><\/a>$`)
	regPlaylistImage := regexp.MustCompile(`^\<img\sclass="j\-flag"\ssrc="([^"]+)"\/\>$`)
	var images []string
	scanner := bufio.NewScanner(bytes.NewReader(content))
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		ss := regPlaylistInfo.FindAllStringSubmatch(line, -1)
		if len(ss) == 1 && len(ss[0]) == 3 {
			res = append(res, Playlist{
				ID:       ss[0][2],
				Title:    ss[0][1],
				Provider: "netease",
				URL:      fmt.Sprintf(`https://music.163.com/#/playlist?id=%s`, ss[0][2]),
			})
		}
		ss = regPlaylistImage.FindAllStringSubmatch(line, -1)
		if len(ss) == 1 && len(ss[0]) == 2 {
			images = append(images, ss[0][1])
		}
	}
	for i := 0; i < len(res) && i < len(images); i++ {
		res[i].Image = images[i]
	}
	return
}

type neteasePlaylistDetail struct {
	Code     int `json:"code"`
	Playlist struct {
		Tracks []struct {
			Name string `json:"name"`
			ID   int    `json:"id"`
			AL   struct {
				Name   string `json:"name"`
				PicURL string `json:"picUrl"`
			} `json:"al"`
			AR []struct {
				Name string `json:"name"`
			} `json:"ar"`
		} `json:"tracks"`
		TrackIDs []struct {
			ID int `json:"id"`
		} `json:"trackIds"`
	} `json:"playlist"`
}

type neteaseSongDetail struct {
	Songs []struct {
		Name string `json:"name"`
		ID   int    `json:"id"`
		AL   struct {
			Name   string `json:"name"`
			PicURL string `json:"picUrl"`
		} `json:"al"`
		AR []struct {
			Name string `json:"name"`
		} `json:"ar"`
	} `json:"songs"`
}

func (p *netease) PlaylistDetail(pl Playlist) (res Songs, err error) {
	data := map[string]interface{}{
		"id":         pl.ID,
		"csrf_token": "",
		"offset":     0,
		"total":      true,
		"limit":      1000,
		"n":          1000,
	}

	params := weapi(data)
	values := url.Values{}
	for k, vs := range params {
		values.Add(k, vs.(string))
	}
	postBody := values.Encode()
	req, err := http.NewRequest("POST", neteaseAPIPlaylistDetail, strings.NewReader(postBody))
	if err != nil {
		return
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:78.0) Gecko/20100101 Firefox/78.0")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "http://music.163.com/")
	req.Header.Set("Origin", "http://music.163.com/")
	req.Header.Set("Accept-Language", "zh-CN,zh-HK;q=0.8,zh-TW;q=0.6,en-US;q=0.4,en;q=0.2")
	req.Header.Set("Accept-Encoding", "gzip, deflate")

	resp, err := httpClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return res, ErrStatusNotOK
	}

	content, err := util.ReadHttpResponseBody(resp)
	if err != nil {
		return
	}

	var plds neteasePlaylistDetail
	if err = json.Unmarshal(content, &plds); err != nil {
		return
	}

	var ids []string
	var c []string
	for _, pld := range plds.Playlist.TrackIDs {
		ids = append(ids, strconv.Itoa(pld.ID))
		c = append(c, fmt.Sprintf(`{"id":%d}`, pld.ID))
	}

	// song detail
	data = map[string]interface{}{
		"ids": fmt.Sprintf(`[%s]`, strings.Join(ids, ",")),
		"c":   fmt.Sprintf(`[%s]`, strings.Join(c, ",")),
	}

	params = weapi(data)
	values = url.Values{}
	for k, vs := range params {
		values.Add(k, vs.(string))
	}
	postBody = values.Encode()
	req, err = http.NewRequest("POST", neteaseAPISongDetail, strings.NewReader(postBody))
	if err != nil {
		return
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:78.0) Gecko/20100101 Firefox/78.0")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "http://music.163.com/")
	req.Header.Set("Origin", "http://music.163.com/")
	req.Header.Set("Accept-Language", "zh-CN,zh-HK;q=0.8,zh-TW;q=0.6,en-US;q=0.4,en;q=0.2")
	req.Header.Set("Accept-Encoding", "gzip, deflate")

	resp, err = httpClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return res, ErrStatusNotOK
	}

	content, err = util.ReadHttpResponseBody(resp)
	if err != nil {
		return
	}

	var sd neteaseSongDetail
	if err = json.Unmarshal(content, &sd); err != nil {
		return
	}

	for _, d := range sd.Songs {
		res = append(res, Song{
			ID:       strconv.Itoa(d.ID),
			Title:    d.Name,
			Image:    d.AL.PicURL,
			Provider: "netease",
			Artist:   d.AR[0].Name,
		})
	}

	return
}

func (p *netease) Name() string {
	return "netease"
}
