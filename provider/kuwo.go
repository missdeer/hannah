package provider

type kuwo struct {
}

func (p *kuwo) Search(keyword string, page int, limit int) (SearchResult, error) {
	return nil, nil
}

func (p *kuwo) SongURL(song Song) (string, error){
	return "", nil
}

func (p *kuwo) HotPlaylist(page int) (Playlists, error) {
	return nil, nil
}

func (p *kuwo) PlaylistDetail(pl Playlist) (Songs, error){
	return nil, nil
}
