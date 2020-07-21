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

}

func TestBilibili_SongURL(t *testing.T) {
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
