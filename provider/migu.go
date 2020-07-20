package provider

type migu struct {
}

func (p *migu) Search(keyword string, page int, limit int) (SearchResult, error) {
	return nil, nil
}

func (p *migu) SongURL(song Song) (string, error){
	return "", nil
}

func (p *migu) HotPlaylist(page int) (Playlists, error) {
	return nil, nil
}

func (p *migu) PlaylistDetail(pl Playlist) (Songs, error) {
	return nil, nil
}
