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

	songs, err := p.SearchSongs("backstreet", 1, 25)
	if err != nil {
		t.Error(err)
	}

	if len(songs) == 0 {
		t.Error("can't found songs for backstreet")
	}

}

func TestKuwo_ResolveSongURL(t *testing.T) {
	p := GetProvider("kuwo")
	if p == nil {
		t.Error("can't get provider")
	}

	_, err := p.ResolveSongURL(Song{ID: `15195332`})
	if err != nil {
		t.Error(err)
	}
}

func TestKuwo_ResolveSongLyric(t *testing.T) {
	p := GetProvider("kuwo")
	if p == nil {
		t.Error("can't get provider")
	}

	_, err := p.ResolveSongLyric(Song{ID: `147917739`}, "lrc")
	if err != nil {
		t.Error(err)
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

func TestKuwo_ArtistSongs(t *testing.T) {
	p := GetProvider("kuwo")
	if p == nil {
		t.Error("can't get provider")
	}

	songs, err := p.ArtistSongs("5335193")
	if err != nil {
		t.Error(err)
	}

	if len(songs) == 0 {
		t.Error("can't found artist songs")
	}
}

func TestKuwo_AlbumSongs(t *testing.T) {
	p := GetProvider("kuwo")
	if p == nil {
		t.Error("can't get provider")
	}

	songs, err := p.AlbumSongs("12997")
	if err != nil {
		t.Error(err)
	}

	if len(songs) == 0 {
		t.Error("can't found album songs")
	}
}
