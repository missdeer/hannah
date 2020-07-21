package provider

import (
	"testing"
)

func TestKuwo_HotPlaylist(t *testing.T) {
	p := GetProvider("kuwo")
	if p == nil {
		t.Error("can't get provider")
	}

}

func TestKuwo_PlaylistDetail(t *testing.T) {
	p := GetProvider("kuwo")
	if p == nil {
		t.Error("can't get provider")
	}

}

func TestKuwo_Search(t *testing.T) {
	p := GetProvider("kuwo")
	if p == nil {
		t.Error("can't get provider")
	}

	songs, err := p.Search("backstreet", 0, 25)
	if err != nil {
		t.Error(err)
	}

	if len(songs) == 0 {
		t.Error("can't found songs for backstreet")
	}

}

func TestKuwo_SongDetail(t *testing.T) {
	p := GetProvider("kuwo")
	if p == nil {
		t.Error("can't get provider")
	}

}

func TestKuwo_Name(t *testing.T) {
	p := GetProvider("kuwo")
	if p == nil {
		t.Error("can't get provider")
	}
	if p.Name() != "kuwo" {
		t.Error("provider name mismatched")
	}
}