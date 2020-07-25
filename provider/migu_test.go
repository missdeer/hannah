package provider

import (
	"testing"
)

func TestMigu_HotPlaylist(t *testing.T) {
	p := GetProvider("migu")
	if p == nil {
		t.Error("can't get provider")
	}

}

func TestMigu_PlaylistDetail(t *testing.T) {
	p := GetProvider("migu")
	if p == nil {
		t.Error("can't get provider")
	}

}

func TestMigu_Search(t *testing.T) {
	p := GetProvider("migu")
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

func TestMigu_ResolveSongURL(t *testing.T) {
	p := GetProvider("migu")
	if p == nil {
		t.Error("can't get provider")
	}

}

func TestMigu_ResolveSongLyric(t *testing.T) {
	p := GetProvider("migu")
	if p == nil {
		t.Error("can't get provider")
	}

}

func TestMigu_Name(t *testing.T) {
	p := GetProvider("migu")
	if p == nil {
		t.Error("can't get provider")
	}
	if p.Name() != "migu" {
		t.Error("provider name mismatched")
	}
}