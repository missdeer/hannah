package provider

import (
	"testing"
)

func TestBilibili_HotPlaylist(t *testing.T) {
	p := GetProvider("bilibili")
	if p == nil {
		t.Error("can't get provider")
	}

}

func TestBilibili_PlaylistDetail(t *testing.T) {
	p := GetProvider("bilibili")
	if p == nil {
		t.Error("can't get provider")
	}

}

func TestBilibili_Search(t *testing.T) {
	p := GetProvider("bilibili")
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

func TestBilibili_ResolveSongURL(t *testing.T) {
	p := GetProvider("bilibili")
	if p == nil {
		t.Error("can't get provider")
	}

}

func TestBilibili_ResolveSongLyric(t *testing.T) {
	p := GetProvider("bilibili")
	if p == nil {
		t.Error("can't get provider")
	}

}

func TestBilibili_Name(t *testing.T) {
	p := GetProvider("bilibili")
	if p == nil {
		t.Error("can't get provider")
	}
	if p.Name() != "bilibili" {
		t.Error("provider name mismatched")
	}
}
