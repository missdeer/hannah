package provider

import (
	"testing"
)

func TestNetease_Search(t *testing.T) {
	p := GetProvider("netease")
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

func TestNetease_ResolveSongURL(t *testing.T) {
	p := GetProvider("netease")
	if p == nil {
		t.Error("can't get provider")
	}

	u, err := p.ResolveSongURL(Song{ID: "1426649237"})
	if err != nil {
		t.Error(err)
	}
	if u.URL == `` {
		t.Error("incorrect song URL")
	}
}

func TestNetease_ResolveSongLyric(t *testing.T) {
	p := GetProvider("netease")
	if p == nil {
		t.Error("can't get provider")
	}

	u, err := p.ResolveSongLyric(Song{ID: "1426649237"})
	if err != nil {
		t.Error(err)
	}
	if u.Lyric == `` {
		t.Error("incorrect song lyric")
	}
}

func TestNetease_HotPlaylist(t *testing.T) {
	p := GetProvider("netease")
	if p == nil {
		t.Error("can't get provider")
	}

	pl, err := p.HotPlaylist(1, 50)
	if err != nil {
		t.Error(err)
	}
	if len(pl) == 0 {
		t.Error("can't get hot playlist")
	}
}

func TestNetease_PlaylistDetail(t *testing.T) {
	p := GetProvider("netease")
	if p == nil {
		t.Error("can't get provider")
	}

	songs, err := p.PlaylistDetail(Playlist{ID: `5038176324`})
	if err != nil {
		t.Error(err)
	}

	if len(songs) == 0 {
		t.Error("can't get playlist detail")
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

func TestNetease_ArtistSongs(t *testing.T) {
	p := GetProvider("netease")
	if p == nil {
		t.Error("can't get provider")
	}
	r, err := p.ArtistSongs("10559")
	if err != nil {
		t.Error(err)
	}
	if len(r) == 0 {
		t.Error("empty result")
	}
}

func TestNetease_AlbumSongs(t *testing.T) {
	p := GetProvider("netease")
	if p == nil {
		t.Error("can't get provider")
	}
	r, err := p.AlbumSongs("73315403")
	if err != nil {
		t.Error(err)
	}
	if len(r) == 0 {
		t.Error("empty result")
	}
}

func TestNetease_Login(t *testing.T) {
	p := &netease{}
	if p == nil {
		t.Error("can't get provider")
	}

	err := p.Login()
	if err != nil {
		t.Error(err)
	}
}
