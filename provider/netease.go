package provider

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/missdeer/hannah/config"
	"github.com/missdeer/hannah/lyric"
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
	neteaseAPIGetSongsURL              = "http://music.163.com/weapi/song/enhance/player/url"
	neteaseAPISearch                   = `https://music.163.com/weapi/cloudsearch/get/web`
	neteaseAPIGetLyric                 = `http://music.163.com/weapi/song/lyric?csrf_token=`
	neteaseAPIHot                      = `http://music.163.com/discover/playlist/?order=hot&limit=%d&offset=%d`
	neteaseAPIPlaylistDetail           = `http://music.163.com/weapi/v3/playlist/detail`
	neteaseAPISongDetail               = `http://music.163.com/weapi/v3/song/detail`
	neteaseAPIGetArtistSongs           = `http://music.163.com/weapi/v1/artist/%s`
	neteaseAPIGetAlbumSongs            = `http://music.163.com/weapi/v1/album/%s`
	neteaseAPILoginCellphone           = `http://music.163.com/weapi/login/cellphone`
	neteaseAPILoginMail                = `http://music.163.com/weapi/login`
	neteaseAPILoginClientToken         = "1_jVUMqWEPke0/1/Vu56xCmJpo5vP1grjn_SOVVDzOc78w8OKLVZ2JH7IfkjSXqgfmh"
	neteaseAPIRefreshToken             = `https://music.163.com/weapi/login/token/refresh`
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

func (p *netease) SearchSongs(keyword string, page int, limit int) (SearchResult, error) {
	data := map[string]interface{}{
		"limit":  limit,
		"offset": (page - 1) * limit,
		"s":      keyword,
		"type":   1,
	}

	params := weapi(data)
	values := url.Values{}
	for k, vs := range params {
		values.Add(k, vs.(string))
	}
	postBody := values.Encode()
	req, err := http.NewRequest("POST", neteaseAPISearch, strings.NewReader(postBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", config.UserAgent)
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "http://music.163.com/")
	req.Header.Set("Origin", "http://music.163.com/")
	req.Header.Set("Accept-Language", "zh-CN,zh-HK;q=0.8,zh-TW;q=0.6,en-US;q=0.4,en;q=0.2")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")

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
		"ids": fmt.Sprintf("[%s]", song.ID),
		"br":  320000,
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

	req.Header.Set("User-Agent", config.UserAgent)
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "http://music.163.com/")
	req.Header.Set("Origin", "http://music.163.com/")
	req.Header.Set("Accept-Language", "zh-CN,zh-HK;q=0.8,zh-TW;q=0.6,en-US;q=0.4,en;q=0.2")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")

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

	var songInfo neteaseSongInfo
	if err = json.Unmarshal(content, &songInfo); err != nil {
		return song, err
	}

	if len(songInfo.Data) == 0 || songInfo.Data[0].URL == "" {
		return song, ErrEmptyPURL
	}

	song.Provider = "netease"
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

func (p *netease) ResolveSongLyric(song Song, format string) (Song, error) {
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

	req.Header.Set("User-Agent", config.UserAgent)
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "http://music.163.com/")
	req.Header.Set("Origin", "http://music.163.com/")
	req.Header.Set("Accept-Language", "zh-CN,zh-HK;q=0.8,zh-TW;q=0.6,en-US;q=0.4,en;q=0.2")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")

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

	var lrc neteaseLyricDetail
	if err = json.Unmarshal(content, &lrc); err != nil {
		return song, nil
	}

	song.Lyric = lyric.LyricConvert("lrc", format, lrc.LRC.Lyric)
	return song, nil
}

func (p *netease) HotPlaylist(page int, limit int) (res Playlists, err error) {
	u := fmt.Sprintf(neteaseAPIHot, limit, (page-1)*limit)
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return
	}

	req.Header.Set("User-Agent", config.UserAgent)
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "http://music.163.com/")
	req.Header.Set("Origin", "http://music.163.com/")
	req.Header.Set("Accept-Language", "zh-CN,zh-HK;q=0.8,zh-TW;q=0.6,en-US;q=0.4,en;q=0.2")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")

	httpClient := util.GetHttpClient()
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

type neteaseTrackIDs []struct {
	ID int `json:"id"`
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
		TrackIDs neteaseTrackIDs `json:"trackIds"`
	} `json:"playlist"`
}

type neteaseSongs []struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
	AL   struct {
		Name   string `json:"name"`
		PicURL string `json:"picUrl"`
	} `json:"al"`
	AR []struct {
		Name string `json:"name"`
	} `json:"ar"`
}

type neteaseSongDetail struct {
	Songs neteaseSongs `json:"songs"`
}

func (p *netease) getSongList(trackIDs neteaseTrackIDs) (res neteaseSongs, err error) {
	var ids []string
	var c []string
	for _, trackID := range trackIDs {
		ids = append(ids, strconv.Itoa(trackID.ID))
		c = append(c, fmt.Sprintf(`{"id":%d}`, trackID.ID))
	}

	// song detail
	data := map[string]interface{}{
		"ids": fmt.Sprintf(`[%s]`, strings.Join(ids, ",")),
		"c":   fmt.Sprintf(`[%s]`, strings.Join(c, ",")),
	}

	params := weapi(data)
	values := url.Values{}
	for k, vs := range params {
		values.Add(k, vs.(string))
	}
	postBody := values.Encode()
	req, err := http.NewRequest("POST", neteaseAPISongDetail, strings.NewReader(postBody))
	if err != nil {
		return
	}

	req.Header.Set("User-Agent", config.UserAgent)
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "http://music.163.com/")
	req.Header.Set("Origin", "http://music.163.com/")
	req.Header.Set("Accept-Language", "zh-CN,zh-HK;q=0.8,zh-TW;q=0.6,en-US;q=0.4,en;q=0.2")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")

	httpClient := util.GetHttpClient()
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

	var sd neteaseSongDetail
	if err = json.Unmarshal(content, &sd); err != nil {
		return nil, err
	}
	return sd.Songs, nil
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

	req.Header.Set("User-Agent", config.UserAgent)
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "http://music.163.com/")
	req.Header.Set("Origin", "http://music.163.com/")
	req.Header.Set("Accept-Language", "zh-CN,zh-HK;q=0.8,zh-TW;q=0.6,en-US;q=0.4,en;q=0.2")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")

	httpClient := util.GetHttpClient()
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

	for i := 0; i < len(plds.Playlist.TrackIDs); i += 200 {
		end := i + 200
		if end > len(plds.Playlist.TrackIDs) {
			end = len(plds.Playlist.TrackIDs)
		}
		sd, err := p.getSongList(plds.Playlist.TrackIDs[i:end])
		if err != nil {
			return res, err
		}

		for _, d := range sd {
			res = append(res, Song{
				ID:       strconv.Itoa(d.ID),
				Title:    d.Name,
				Image:    d.AL.PicURL,
				Provider: "netease",
				Artist:   d.AR[0].Name,
			})
		}
	}

	return
}

type neteaseArtistSongs struct {
	Artist struct {
		Name      string `json:"name"`
		PicURL    string `json:"picUrl"`
		Img1v1URL string `json:"img1v1Url"`
	} `json:"artist"`
	HotSongs []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"hotSongs"`
}

func (p *netease) ArtistSongs(id string) (res Songs, err error) {
	data := map[string]interface{}{}

	params := weapi(data)
	values := url.Values{}
	for k, vs := range params {
		values.Add(k, vs.(string))
	}
	postBody := values.Encode()
	u := fmt.Sprintf(neteaseAPIGetArtistSongs, id)
	req, err := http.NewRequest("POST", u, strings.NewReader(postBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", config.UserAgent)
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "http://music.163.com/")
	req.Header.Set("Origin", "http://music.163.com/")
	req.Header.Set("Accept-Language", "zh-CN,zh-HK;q=0.8,zh-TW;q=0.6,en-US;q=0.4,en;q=0.2")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")

	httpClient := util.GetHttpClient()
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return res, ErrStatusNotOK
	}

	content, err := util.ReadHttpResponseBody(resp)
	if err != nil {
		return nil, err
	}

	var sd neteaseArtistSongs
	if err = json.Unmarshal(content, &sd); err != nil {
		return
	}

	for _, d := range sd.HotSongs {
		res = append(res, Song{
			ID:       strconv.Itoa(d.ID),
			Title:    d.Name,
			Image:    sd.Artist.PicURL,
			Provider: "netease",
			Artist:   sd.Artist.Name,
		})
	}
	return res, nil
}

type neteaseAlbumSongs struct {
	Code  int `json:"code"`
	Songs []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		AL   struct {
			PicURL string `json:"picUrl"`
		} `json:"al"`
		AR []struct {
			Name string `json:"name"`
		} `json:"ar"`
	} `json:"songs"`
}

func (p *netease) AlbumSongs(id string) (res Songs, err error) {
	data := map[string]interface{}{}

	params := weapi(data)
	values := url.Values{}
	for k, vs := range params {
		values.Add(k, vs.(string))
	}
	postBody := values.Encode()
	u := fmt.Sprintf(neteaseAPIGetAlbumSongs, id)
	req, err := http.NewRequest("POST", u, strings.NewReader(postBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", config.UserAgent)
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "http://music.163.com/")
	req.Header.Set("Origin", "http://music.163.com/")
	req.Header.Set("Accept-Language", "zh-CN,zh-HK;q=0.8,zh-TW;q=0.6,en-US;q=0.4,en;q=0.2")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")

	httpClient := util.GetHttpClient()
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return res, ErrStatusNotOK
	}

	content, err := util.ReadHttpResponseBody(resp)
	if err != nil {
		return nil, err
	}

	var sd neteaseAlbumSongs
	if err = json.Unmarshal(content, &sd); err != nil {
		return
	}

	for _, d := range sd.Songs {
		var artists []string
		for _, a := range d.AR {
			artists = append(artists, a.Name)
		}
		res = append(res, Song{
			ID:       strconv.Itoa(d.ID),
			Title:    d.Name,
			Image:    d.AL.PicURL,
			Provider: "netease",
			Artist:   strings.Join(artists, "/"),
		})
	}
	return res, nil
}

func (p *netease) Login() error {
	username := config.NeteaseUsername
	password := config.NeteasePassword
	if username == "" || password == "" {
		return ErrNoAuthorizeInfo
	}
	sum := md5.Sum([]byte(password))
	data := map[string]interface{}{
		"password":      hex.EncodeToString(sum[:]),
		"rememberLogin": "true",
	}
	r := regexp.MustCompile(`^[0-9]+$`)
	var u string
	if r.MatchString(username) {
		u = neteaseAPILoginCellphone
		data["phone"] = username
	} else {
		u = neteaseAPILoginMail
		data["username"] = username
		data["clientToken"] = neteaseAPILoginClientToken
	}

	params := weapi(data)
	values := url.Values{}
	for k, vs := range params {
		values.Add(k, vs.(string))
	}
	postBody := values.Encode()
	req, err := http.NewRequest("POST", u, strings.NewReader(postBody))
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", config.UserAgent)
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "http://music.163.com/")
	req.Header.Set("Origin", "http://music.163.com/")
	req.Header.Set("Accept-Language", "zh-CN,zh-HK;q=0.8,zh-TW;q=0.6,en-US;q=0.4,en;q=0.2")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Cookie", "os=pc; osver=Microsoft-Windows-10-Professional-build-10586-64bit; appver=2.0.3.131777; channel=netease; __remember_me=true;")

	httpClient := util.GetHttpClient()
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return ErrStatusNotOK
	}
	b, err := util.ReadHttpResponseBody(resp)
	if err != nil {
		return err
	}
	var res map[string]interface{}
	err = json.Unmarshal(b, &res)
	if err != nil {
		return err
	}
	code, ok := res["code"]
	if ok {
		if code.(float64) != 200 {
			return errors.New(res["message"].(string))
		}
	}
	log.Println(string(b))

	return nil
}

func (p *netease) RefreshToken() error {
	data := map[string]interface{}{}

	params := weapi(data)
	values := url.Values{}
	for k, vs := range params {
		values.Add(k, vs.(string))
	}
	postBody := values.Encode()
	req, err := http.NewRequest("POST", neteaseAPIRefreshToken, strings.NewReader(postBody))
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", config.UserAgent)
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "http://music.163.com/")
	req.Header.Set("Origin", "http://music.163.com/")
	req.Header.Set("Accept-Language", "zh-CN,zh-HK;q=0.8,zh-TW;q=0.6,en-US;q=0.4,en;q=0.2")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")

	httpClient := util.GetHttpClient()
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return ErrStatusNotOK
	}

	_, err = util.ReadHttpResponseBody(resp)
	if err != nil {
		return err
	}

	return nil
}

func (p *netease) Name() string {
	return "netease"
}
