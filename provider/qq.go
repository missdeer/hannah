package provider

type qq struct {
}

func (p *qq) Search(keyword string, page int, limit int) (SearchResult, error) {
	return nil, nil
}

func (p *qq) SongURL(song Song) (string, error){
	return "", nil
}

func (p *qq) HotPlaylist(page int) (Playlists, error) {
	return nil, nil
}

func (p *qq) PlaylistDetail(pl Playlist) (Songs, error){
	return nil, nil
}

func (p *qq) Name() string {
	return "qq"
}
