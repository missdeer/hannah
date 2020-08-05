package provider

import (
	"bufio"
	"bytes"
	"crypto/aes"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"

	"github.com/missdeer/hannah/util"
	"github.com/missdeer/hannah/util/cryptography"
)

const (
	miguAESPassphrase = "4ea5c508a6566e76240543f8feb06fd457777be39549c4016436afda65d2330e"
	miguRSAPublicKey  = "-----BEGIN PUBLIC KEY-----\nMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC8asrfSaoOb4je+DSmKdriQJKW\nVJ2oDZrs3wi5W67m3LwTB9QVR+cE3XWU21Nx+YBxS0yun8wDcjgQvYt625ZCcgin\n2ro/eOkNyUOTBIbuj9CvMnhUYiR61lC1f1IGbrSYYimqBVSjpifVufxtx/I3exRe\nZosTByYp4Xwpb1+WAQIDAQAB\n-----END PUBLIC KEY-----"
)

var (
	miguAPISearch         = `http://m.music.migu.cn/migu/remoting/scr_search_tag?type=2&keyword=%s&pgc=%d&rows=%d`
	miguAPIHot            = `https://music.migu.cn/v3/music/playlist?page=%d`
	miguAPIPlaylistDetail = `https://music.migu.cn/v3/music/playlist/%s`
	miguAPIGetPlayInfo    = `https://music.migu.cn/v3/api/music/audioPlayer/getPlayInfo?dataType=2&data=%s&secKey=%s`
	miguAPILyric          = `https://music.migu.cn/v3/api/music/audioPlayer/getLyric?copyrightId=%s`

	regPlaylist       = regexp.MustCompile(`data\-share='([^']+)'`)
	regPlaylistLink   = regexp.MustCompile(`^\/v3\/music\/playlist\/([0-9]+)\?origin=[0-9]+$`)
	regSongInPlaylist = regexp.MustCompile(`^<a\sclass="song\-name\-txt"\shref="([^"]+)"\stitle="([^"]+)"\starget="_blank">`)
	regSongLink       = regexp.MustCompile(`^\/v3\/music\/song\/([0-9]+)$`)
)

type migu struct {
}

type miguSearchResult struct {
	Musics []struct {
		AlbumName  string `json:"albumName"`
		AlbumID    string `json:"albumId"`
		MP3        string `json:"mp3"`
		SongName   string `json:"songName"`
		Lyrics     string `json:"lyrics"`
		ID         string `json:"id"`
		Title      string `json:"title"`
		Cover      string `json:"cover"`
		SingerName string `json:"singerName"`
		Artist     string `json:"artist"`
	} `json:"musics"`
	Pgt     int    `json:"pgt"`
	Keyword string `json:"keyword"`
	PageNo  string `json:"pageNo"`
	Success bool   `json:"success"`
}

func (p *migu) Search(keyword string, page int, limit int) (SearchResult, error) {
	u := fmt.Sprintf(miguAPISearch, keyword, page, limit)

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:78.0) Gecko/20100101 Firefox/78.0")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "http://migu.cn/")
	req.Header.Set("Origin", "http://migu.cn/")
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

	var sr miguSearchResult
	err = json.Unmarshal(content, &sr)
	if err != nil {
		return nil, err
	}

	var res SearchResult
	for _, music := range sr.Musics {
		res = append(res, Song{
			ID:       music.ID,
			URL:      music.MP3,
			Title:    music.Title,
			Image:    music.Cover,
			Artist:   music.Artist,
			Lyric:    music.Lyrics,
			Provider: "migu",
		})
	}

	return res, nil
}

type miguSongInfo struct {
	ReturnCode string `json:"returncode"`
	Msg        string `json:"msg"`
	Data       struct {
		PlayURL string `json:"playUrl"`
	} `json:"data"`
}

func (p *migu) ResolveSongURL(song Song) (Song, error) {
	if song.URL != "" {
		return song, nil
	}
	aesInput := fmt.Sprintf(`{"copyrightId":"%s","type":2,"auditionsFlag":0}`, song.ID)
	aesKey, _ := hex.DecodeString(miguAESPassphrase)
	plaintext := []byte(aesInput)

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return song, err
	}
	cipherText := cryptography.AESCBCEncrypt(plaintext, aesKey, iv)
	// TODO https://github.com/listen1/listen1_chrome_extension/blob/6c66da63e0e6f77de047392f2cc1d4cf9fc85951/js/provider/migu.js#L128
	aesResult := base64.StdEncoding.EncodeToString(append(iv, cipherText...))

	rsaPubKey, err := cryptography.BytesToPublicKey([]byte(miguRSAPublicKey))
	if err != nil {
		return song, err
	}
	secKey, err := cryptography.EncryptWithPublicKey(aesKey, rsaPubKey)
	if err != nil {
		return song, err
	}
	u := fmt.Sprintf(miguAPIGetPlayInfo, aesResult, url.QueryEscape(string(secKey)))

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return song, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:78.0) Gecko/20100101 Firefox/78.0")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "http://migu.cn/")
	req.Header.Set("Origin", "http://migu.cn/")
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

	var si miguSongInfo
	if err = json.Unmarshal(content, &si); err != nil {
		return song, err
	}
	song.URL = si.Data.PlayURL

	return song, nil
}

type miguLyric struct {
	ReturnCode string `json:"returncode"`
	Msg        string `json:"msg"`
	Lyric      string `json:"lyric"`
}

func (p *migu) ResolveSongLyric(song Song) (Song, error) {
	u := fmt.Sprintf(miguAPILyric, song.ID)

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return song, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:78.0) Gecko/20100101 Firefox/78.0")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "http://migu.cn/")
	req.Header.Set("Origin", "http://migu.cn/")
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

	var lyric miguLyric
	if err = json.Unmarshal(content, &lyric); err != nil {
		return song, err
	}
	song.Lyric = lyric.Lyric

	return song, nil
}

type miguPlaylist struct {
	Type    string `json:"type"`
	Title   string `json:"title"`
	LinkURL string `json:"linkUrl"`
	ImgURL  string `json:"imgUrl"`
}

func (p *migu) HotPlaylist(page int) (res Playlists, err error) {
	u := fmt.Sprintf(miguAPIHot, page)

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:78.0) Gecko/20100101 Firefox/78.0")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "http://migu.cn/")
	req.Header.Set("Origin", "http://migu.cn/")
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
	scanner := bufio.NewScanner(bytes.NewReader(content))
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		ss := regPlaylist.FindAllStringSubmatch(line, -1)
		if len(ss) == 1 && len(ss[0]) == 2 {
			var pl miguPlaylist
			if err = json.Unmarshal([]byte(ss[0][1]), &pl); err != nil {
				continue
			}
			ss := regPlaylistLink.FindAllStringSubmatch(pl.LinkURL, -1)
			if len(ss) == 1 && len(ss[0]) == 2 {
				res = append(res, Playlist{
					ID:       ss[0][1],
					Image:    "http:" + pl.ImgURL,
					URL:      "https://music.migu.cn" + pl.LinkURL,
					Title:    pl.Title,
					Provider: "migu",
				})
			}
		}
	}
	return res, nil
}

func (p *migu) PlaylistDetail(pl Playlist) (songs Songs, err error) {
	u := fmt.Sprintf(miguAPIPlaylistDetail, pl.ID)
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:78.0) Gecko/20100101 Firefox/78.0")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "http://migu.cn/")
	req.Header.Set("Origin", "http://migu.cn/")
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
	scanner := bufio.NewScanner(bytes.NewReader(content))
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		ss := regSongInPlaylist.FindAllStringSubmatch(line, -1)
		if len(ss) == 1 && len(ss[0]) == 3 {
			sss := regSongLink.FindAllStringSubmatch(ss[0][1], -1)
			if len(sss) == 1 && len(sss[0]) == 2 {
				songs = append(songs, Song{
					ID:       sss[0][1],
					Title:    ss[0][2],
					Provider: "migu",
				})
			}
		}
	}
	return songs, nil
}

func (p *migu) Name() string {
	return "migu"
}
