package provider

import (
	"net/url"
	"path/filepath"
	"testing"
)

func TestKugou_HotPlaylist(t *testing.T) {
	p := GetProvider("kugou")
	if p == nil {
		t.Error("can't get provider")
	}

}

func TestKugou_PlaylistDetail(t *testing.T) {
	p := GetProvider("kugou")
	if p == nil {
		t.Error("can't get provider")
	}

}

func TestKugou_Search(t *testing.T) {
	p := GetProvider("kugou")
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

func TestKugou_ResolveSongURL(t *testing.T) {
	p := GetProvider("kugou")
	if p == nil {
		t.Error("can't get provider")
	}

	u, err:= p.ResolveSongURL(Song{ID: "F3EA661A19E9A0C965AD64049040BBAC"})
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

func TestKugou_ResolveSongLyric(t *testing.T) {
	p := GetProvider("kugou")
	if p == nil {
		t.Error("can't get provider")
	}

	_, err := p.ResolveSongLyric(Song{ID: "F3EA661A19E9A0C965AD64049040BBAC"}, "lrc")
	if err != nil {
		t.Error(err)
	}
}

func TestKugou_Name(t *testing.T) {
	p := GetProvider("kugou")
	if p == nil {
		t.Error("can't get provider")
	}
	if p.Name() != "kugou" {
		t.Error("provider name mismatched")
	}
}
