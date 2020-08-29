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

	songs, err := p.SearchSongs("Lydia", 1, 25)
	if err != nil {
		t.Error(err)
	}

	if len(songs) == 0 {
		t.Error("can't found songs for Lydia")
	}

}

func TestMigu_ResolveSongURL(t *testing.T) {
	p := GetProvider("migu")
	if p == nil {
		t.Error("can't get provider")
	}

	_, err := p.ResolveSongURL(Song{ID: `6005752J00T`})
	if err != nil {
		t.Error(err)
	}
}

func TestMigu_ResolveSongLyric(t *testing.T) {
	p := GetProvider("migu")
	if p == nil {
		t.Error("can't get provider")
	}

	_, err := p.ResolveSongLyric(Song{ID: `6005752J00T`}, "lrc")
	if err != nil {
		t.Error(err)
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

func TestMigu_AlbumSongs(t *testing.T) {
	p := GetProvider("migu")
	if p == nil {
		t.Error("can't get provider")
	}

	songs, err := p.AlbumSongs("1135110570")
	if err != nil {
		t.Error(err)
	}

	if len(songs) == 0 {
		t.Error("can't found songs for album")
	}
}

func TestMigu_ArtistSongs(t *testing.T) {
	p := GetProvider("migu")
	if p == nil {
		t.Error("can't get provider")
	}

	songs, err := p.ArtistSongs("266")
	if err != nil {
		t.Error(err)
	}

	if len(songs) == 0 {
		t.Error("can't found songs for artist")
	}
}
