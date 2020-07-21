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

}

func TestMigu_SongURL(t *testing.T) {
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