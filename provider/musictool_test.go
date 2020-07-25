package provider

import (
	"testing"
)

func TestMusictool_SetProvider(t *testing.T) {
	p := GetProvider("musictool")
	if p == nil {
		t.Error("can't get provider")
	}

	mt, ok := p.(*musictool)
	if !ok {
		t.Error("not a musictool instance")
	}

	if mt.provider != "" {
		t.Error("init provider field not nil")
	}

	mt.SetProvider("test")
	if mt.provider != "test" {
		t.Error("set provider failed")
	}
}

func TestMusictool_HotPlaylist(t *testing.T) {
	p := GetProvider("musictool")
	if p == nil {
		t.Error("can't get provider")
	}

}

func TestMusictool_PlaylistDetail(t *testing.T) {
	p := GetProvider("musictool")
	if p == nil {
		t.Error("can't get provider")
	}

}

func TestMusictool_Search(t *testing.T) {
	p := GetProvider("musictool")
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

func TestMusictool_ResolveSongURL(t *testing.T) {
	p := GetProvider("musictool")
	if p == nil {
		t.Error("can't get provider")
	}

}

func TestMusictool_ResolveSongLyric(t *testing.T) {
	p := GetProvider("musictool")
	if p == nil {
		t.Error("can't get provider")
	}

}

func TestMusictool_Name(t *testing.T) {
	p := GetProvider("musictool")
	if p == nil {
		t.Error("can't get provider")
	}
	if p.Name() != "musictool" {
		t.Error("provider name mismatched")
	}
}
