package provider

type migu struct {
}

func (p *migu) Search(keyword string, page int, limit int) (SearchResult, error) {
	return nil, nil
}

func (p *migu) ResolveSongURL(song Song) (Song, error){
	return song, nil
}

func (p *migu) ResolveSongLyric(song Song) (Song, error){
	return song, nil
}

func (p *migu) HotPlaylist(page int) (Playlists, error) {
	return nil, nil
}

func (p *migu) PlaylistDetail(pl Playlist) (Songs, error) {
	return nil, nil
}

func (p *migu) Name() string {
	return "migu"
}


