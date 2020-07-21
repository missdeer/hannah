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

}

func TestKugou_SongURL(t *testing.T) {
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