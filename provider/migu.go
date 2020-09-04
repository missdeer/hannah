package provider

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"crypto/rsa"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/missdeer/hannah/config"
	"github.com/missdeer/hannah/lyric"
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
	miguAPIArtistSongs    = `https://music.migu.cn/v3/music/artist/%s/song?page=%d`
	miguAPIAlbumSongs     = `https://music.migu.cn/v3/music/album/%s`
	miguAPIGetPlayInfo    = `https://m.music.migu.cn/migu/remoting/cms_detail_tag?cpid=%s`
	miguAPIGetLossless    = `http://music.migu.cn/v3/api/music/audioPlayer/getPlayInfo?dataType=2&`
	miguAPILyric          = `https://music.migu.cn/v3/api/music/audioPlayer/getLyric?copyrightId=%s`

	regPlaylist     = regexp.MustCompile(`data\-share='([^']+)'`)
	regPlaylistLink = regexp.MustCompile(`^\/v3\/music\/playlist\/([0-9]+)\?origin=[0-9]+$`)
	regSongs        = regexp.MustCompile(`(?m)data\-share='{\n"type":"song",\n"title":"[^"]+",\n"linkUrl":"\/v3\/music\/song\/(\w+)",\n"imgUrl":"([^"]+)",\n"summary":"([^"]+)",\n"singer":"([^"]+)",\n"album":"[^"]+"\n?}`)

	rsaPublicKey *rsa.PublicKey
)

func getRsaPublicKey() (*rsa.PublicKey, error) {
	var err error = nil
	if rsaPublicKey == nil {
		rsaPublicKey, err = cryptography.ParsePublicKey([]byte(miguRSAPublicKey))
	}
	return rsaPublicKey, err
}

type migu struct {
}

type miguSearchResult struct {
	Musics []struct {
		AlbumName   string `json:"albumName"`
		AlbumID     string `json:"albumId"`
		MP3         string `json:"mp3"`
		CopyrightID string `json:"copyrightId"`
		SongName    string `json:"songName"`
		Lyrics      string `json:"lyrics"`
		ID          string `json:"id"`
		Title       string `json:"title"`
		Cover       string `json:"cover"`
		SingerName  string `json:"singerName"`
		Artist      string `json:"artist"`
	} `json:"musics"`
	Pgt     int    `json:"pgt"`
	Keyword string `json:"keyword"`
	PageNo  string `json:"pageNo"`
	Success bool   `json:"success"`
}

func (p *migu) SearchSongs(keyword string, page int, limit int) (SearchResult, error) {
	u := fmt.Sprintf(miguAPISearch, url.QueryEscape(keyword), page, limit)

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", config.UserAgent)
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "http://migu.cn/")
	req.Header.Set("Origin", "http://migu.cn/")
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

	var sr miguSearchResult
	err = json.Unmarshal(content, &sr)
	if err != nil {
		return nil, err
	}

	var res SearchResult
	for _, music := range sr.Musics {
		res = append(res, Song{
			ID:       music.CopyrightID,
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
	Data struct {
		ListenURL  string   `json:"listenUrl"`
		Lyric      string   `json:"lyricLrc"`
		SongName   string   `json:"songName"`
		SingerName []string `json:"singerName"`
		SongID     string   `json:"songId"`
		PicL       string   `json:"picL"`
	} `json:"data"`
}

type miguSongURL struct {
	Data struct {
		PlayURL string `json:"playUrl"`
	} `json:"data"`
}

func (p *migu) ResolveSongURL(song Song) (Song, error) {
	u := fmt.Sprintf(miguAPIGetPlayInfo, song.ID)

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return song, err
	}

	req.Header.Set("User-Agent", config.UserAgent)
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "http://migu.cn/")
	req.Header.Set("Origin", "http://migu.cn/")
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

	var si miguSongInfo
	if err = json.Unmarshal(content, &si); err != nil {
		return song, err
	}
	song.URL = si.Data.ListenURL
	song.Title = si.Data.SongName
	song.Artist = strings.Join(si.Data.SingerName, "/")
	song.Image = si.Data.PicL

	u = fmt.Sprintf("%s%s", miguAPIGetLossless, p.encrypt(fmt.Sprintf(`{"copyrightId":"%s", "type":3}`, song.ID)))

	req, err = http.NewRequest("GET", u, nil)
	if err != nil {
		return song, err
	}

	req.Header.Set("User-Agent", config.UserAgent)
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "http://migu.cn/")
	req.Header.Set("Origin", "http://migu.cn/")
	req.Header.Set("Accept-Language", "zh-CN,zh-HK;q=0.8,zh-TW;q=0.6,en-US;q=0.4,en;q=0.2")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")

	resp, err = httpClient.Do(req)
	if err != nil {
		return song, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return song, ErrStatusNotOK
	}

	content, err = util.ReadHttpResponseBody(resp)
	if err != nil {
		return song, err
	}

	var songURL miguSongURL
	if err = json.Unmarshal(content, &songURL); err != nil {
		return song, err
	}
	if songURL.Data.PlayURL == "" {
		return song, ErrEmptyPURL
	}
	song.URL = "http:" + songURL.Data.PlayURL

	return song, nil
}

type miguLyric struct {
	ReturnCode string `json:"returncode"`
	Msg        string `json:"msg"`
	Lyric      string `json:"lyric"`
}

func (p *migu) ResolveSongLyric(song Song, format string) (Song, error) {
	if song.Lyric != "" {
		return song, nil
	}
	u := fmt.Sprintf(miguAPILyric, song.ID)

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return song, err
	}

	req.Header.Set("User-Agent", config.UserAgent)
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "http://migu.cn/")
	req.Header.Set("Origin", "http://migu.cn/")
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

	var lrc miguLyric
	if err = json.Unmarshal(content, &lrc); err != nil {
		return song, err
	}
	song.Lyric = lyric.LyricConvert("lrc", format, lrc.Lyric)

	return song, nil
}

type miguPlaylist struct {
	Type    string `json:"type"`
	Title   string `json:"title"`
	LinkURL string `json:"linkUrl"`
	ImgURL  string `json:"imgUrl"`
}

func (p *migu) HotPlaylist(page int, limit int) (res Playlists, err error) {
	u := fmt.Sprintf(miguAPIHot, page)

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", config.UserAgent)
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "http://migu.cn/")
	req.Header.Set("Origin", "http://migu.cn/")
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

	req.Header.Set("User-Agent", config.UserAgent)
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "http://migu.cn/")
	req.Header.Set("Origin", "http://migu.cn/")
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
	ss := regSongs.FindAllSubmatch(content, -1)
	for _, s := range ss {
		if len(s) == 5 {
			songs = append(songs, Song{
				ID:       string(s[1]),
				Image:    "http:" + string(s[2]),
				Title:    string(s[3]),
				Artist:   string(s[4]),
				Provider: "migu",
			})
		}
	}
	return songs, nil
}

func (p *migu) ArtistSongs(id string) (res Songs, err error) {
	u := fmt.Sprintf(miguAPIArtistSongs, id, config.Page)

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", config.UserAgent)
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "http://migu.cn/")
	req.Header.Set("Origin", "http://migu.cn/")
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
	ss := regSongs.FindAllSubmatch(content, -1)
	for _, s := range ss {
		if len(s) == 5 {
			res = append(res, Song{
				ID:       string(s[1]),
				Image:    "http:" + string(s[2]),
				Title:    string(s[3]),
				Artist:   string(s[4]),
				Provider: "migu",
			})
		}
	}
	if len(res) == 0 {
		return nil, ErrEmptyTrackList
	}
	return
}

func (p *migu) AlbumSongs(id string) (res Songs, err error) {
	u := fmt.Sprintf(miguAPIAlbumSongs, id)

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", config.UserAgent)
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "http://migu.cn/")
	req.Header.Set("Origin", "http://migu.cn/")
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
	ss := regSongs.FindAllSubmatch(content, -1)
	for _, s := range ss {
		if len(s) == 5 {
			res = append(res, Song{
				ID:       string(s[1]),
				Image:    "http:" + string(s[2]),
				Title:    string(s[3]),
				Artist:   string(s[4]),
				Provider: "migu",
			})
		}
	}
	if len(res) == 0 {
		return nil, ErrEmptyTrackList
	}
	return
}

func (p *migu) Login() error {
	return ErrNotImplemented
}

func (p *migu) Name() string {
	return "migu"
}

func (p *migu) encrypt(text string) (encryptedData string) {
	// fmt.Println(text)
	text = util.ToJson(util.ParseJson(bytes.NewBufferString(text).Bytes()))
	randomBytes, err := util.GenRandomBytes(32)
	if err != nil {
		fmt.Println(err)
		return encryptedData
	}
	pwd := bytes.NewBufferString(hex.EncodeToString(randomBytes)).Bytes()
	salt, err := util.GenRandomBytes(8)
	if err != nil {
		fmt.Println(err)
		return encryptedData
	}
	// key = []byte{0xaf, 0xb3, 0xac, 0x50, 0xcd, 0x1d, 0x23, 0x81, 0x58, 0x5f, 0xa7, 0xbc, 0xbd, 0x8c, 0xbe, 0x02, 0x56, 0x0f, 0xad, 0xe7, 0xd1, 0x7e, 0x2e, 0xb1, 0x14, 0x81, 0x6f, 0x27, 0xab, 0x7b, 0x6a, 0x75}
	// iv = []byte{0xfb, 0x10, 0x89, 0xb0, 0x13, 0x32, 0xf2, 0xa7, 0x02, 0x51, 0x49, 0xff, 0xbc, 0x16, 0xf0, 0x40}
	// pwd = bytes.NewBufferString("d8e28215ed6573e0fd5eb8b8ae8062542589e96f669bee6503af003c63cdfbd4").Bytes()
	// salt = []byte{0xde, 0xfc, 0x9f, 0x26, 0x29, 0xdd, 0xec, 0x37}
	key, iv := p.derive(pwd, salt, 256, 16)
	var data []byte
	data = append(data, bytes.NewBufferString("Salted__").Bytes()...)
	data = append(data, salt...)
	encryptedD := cryptography.AesEncryptCBCWithIv(bytes.NewBufferString(text).Bytes(), key, iv)
	data = append(data, encryptedD...)
	dat := base64.StdEncoding.EncodeToString(data)
	var rsaB []byte
	pubKey, err := getRsaPublicKey()
	if err == nil {
		rsaB = cryptography.RSAEncryptV2(pwd, pubKey)
	}
	sec := base64.StdEncoding.EncodeToString(rsaB)
	// fmt.Println("data:", dat)
	// fmt.Println("sec:", sec)
	encryptedData = "data=" + url.QueryEscape(dat)
	encryptedData = encryptedData + "&secKey=" + url.QueryEscape(sec)
	return encryptedData
}

func (p *migu) derive(password []byte, salt []byte, keyLength int, ivSize int) ([]byte, []byte) {
	keySize := keyLength / 8
	repeat := math.Ceil(float64(keySize+ivSize*8) / 32)
	var data []byte
	var lastData []byte
	for i := 0.0; i < repeat; i++ {
		var md5Data []byte
		md5Data = append(md5Data, lastData...)
		md5Data = append(md5Data, password...)
		md5Data = append(md5Data, salt...)
		h := md5.New()
		h.Write(md5Data)
		md5Data = h.Sum(nil)
		data = append(data, md5Data...)
		lastData = md5Data
	}
	return data[:keySize], data[keySize : keySize+ivSize]
}
