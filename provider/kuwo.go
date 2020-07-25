package provider

type kuwo struct {
}

func (p *kuwo) Search(keyword string, page int, limit int) (SearchResult, error) {
	return nil, nil
}

func (p *kuwo) ResolveSongURL(song Song) (Song, error) {
	return song, nil
}

func (p *kuwo) ResolveSongLyric(song Song) (Song, error){
	return song, nil
}

func (p *kuwo) HotPlaylist(page int) (Playlists, error) {
	return nil, nil
}

func (p *kuwo) PlaylistDetail(pl Playlist) (Songs, error) {
	return nil, nil
}

func (p *kuwo) Name() string {
	return "kuwo"
}
