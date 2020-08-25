package provider

import (
	"math/rand"
	"net/url"
	"path/filepath"
	"testing"
	"time"
)

func TestQq_HotPlaylist(t *testing.T) {
	p := GetProvider("qq")
	if p == nil {
		t.Error("can't get provider")
	}

}

func TestQq_PlaylistDetail(t *testing.T) {
	p := GetProvider("qq")
	if p == nil {
		t.Error("can't get provider")
	}

	songs, err := p.PlaylistDetail(Playlist{ID: `3602407677`})
	if err != nil {
		t.Error(err)
	}

	if len(songs) == 0 {
		t.Error("can't get playlist detail")
	}
}

func TestQq_Search(t *testing.T) {
	p := GetProvider("qq")
	if p == nil {
		t.Error("can't get provider")
	}

	songs, err := p.SearchSongs("backstreet", 1, 25)
	if err != nil {
		t.Error(err)
	}

	if len(songs) == 0 {
		t.Error("can't found songs for backstreet")
	}
}

func TestQq_ResolveSongURL(t *testing.T) {
	p := GetProvider("qq")
	if p == nil {
		t.Error("can't get provider")
	}

	rand.Seed(time.Now().UnixNano())
	u, err := p.ResolveSongURL(Song{ID: "003VQrF72a0DGb"})
	if err != nil {
		t.Error(err)
	}

	parsedURL, err := url.Parse(u.URL)
	if err != nil {
		t.Error(err)
	}
	if filepath.Base(parsedURL.Path) == ".m4a" {
		t.Error("incorrect song URL")
	}
}

func TestQq_ResolveSongLyric(t *testing.T) {
	p := GetProvider("qq")
	if p == nil {
		t.Error("can't get provider")
	}

	_, err := p.ResolveSongLyric(Song{ID: "003VQrF72a0DGb"})
	if err != nil {
		t.Error(err)
	}
}

func TestQq_Name(t *testing.T) {
	p := GetProvider("qq")
	if p == nil {
		t.Error("can't get provider")
	}
	if p.Name() != "qq" {
		t.Error("provider name mismatched")
	}
}

func TestQq_AlbumSongs(t *testing.T) {
	p := GetProvider("qq")
	if p == nil {
		t.Error("can't get provider")
	}

	r, err := p.AlbumSongs("001IskfD3Vncxo")
	if err != nil {
		t.Error(err)
	}
	if len(r) == 0 {
		t.Error("empty result")
	}
}

func TestQq_ArtistSongs(t *testing.T) {
	p := GetProvider("qq")
	if p == nil {
		t.Error("can't get provider")
	}

	r, err := p.ArtistSongs("004aRKga0CXIPm")
	if err != nil {
		t.Error(err)
	}
	if len(r) == 0 {
		t.Error("empty result")
	}
}
