package provider

import (
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

	songs, err := p.Search("backstreet", 0, 25)
	if err != nil {
		t.Error(err)
	}

	if len(songs) == 0 {
		t.Error("can't found songs for backstreet")
	}

}

func TestKugou_SongDetail(t *testing.T) {
	p := GetProvider("kugou")
	if p == nil {
		t.Error("can't get provider")
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