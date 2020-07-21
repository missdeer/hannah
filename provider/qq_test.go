package provider

import (
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

}

func TestQq_SongURL(t *testing.T) {
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