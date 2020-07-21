package provider

type bilibili struct {
}

func (p *bilibili) Search(keyword string, page int, limit int) (SearchResult, error) {
	return nil, nil
}

func (p *bilibili) SongDetail(song Song) (Song, error) {
	return song, nil
}

func (p *bilibili) HotPlaylist(page int) (Playlists, error) {
	return nil, nil
}

func (p *bilibili) PlaylistDetail(pl Playlist) (Songs, error) {
	return nil, nil
}

func (p *bilibili) Name() string {
	return "bilibili"
}
