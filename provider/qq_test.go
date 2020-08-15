package provider

import (
	"net/url"
	"path/filepath"
	"testing"
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

}

func TestQq_Search(t *testing.T) {
	p := GetProvider("qq")
	if p == nil {
		t.Error("can't get provider")
	}

	songs, err := p.Search("backstreet", 1, 25)
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