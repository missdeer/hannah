package provider

import (
	"testing"
)

func TestNetease_Search(t *testing.T) {
	p := GetProvider("netease")
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

func TestNetease_SongDetail(t *testing.T) {
	p := GetProvider("netease")
	if p == nil {
		t.Error("can't get provider")
	}

	u, err := p.SongDetail(Song{ID: "12341234"})
	if err != nil {
		t.Error(err)
	}
	if u.URL != `http://music.163.com/song/media/outer/url?id=12341234.mp3` {
		t.Error("incorrect song URL")
	}
}

func TestNetease_HotPlaylist(t *testing.T) {
	p := GetProvider("netease")
	if p == nil {
		t.Error("can't get provider")
	}

}

func TestNetease_PlaylistDetail(t *testing.T) {
	p := GetProvider("netease")
	if p == nil {
		t.Error("can't get provider")
	}

}

func TestNetease_Name(t *testing.T) {
	p := GetProvider("netease")
	if p == nil {
		t.Error("can't get provider")
	}
	if p.Name() != "netease" {
		t.Error("provider name mismatched")
	}
}
