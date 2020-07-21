package provider

type kugou struct {
}

func (p *kugou) Search(keyword string, page int, limit int) (SearchResult, error) {
	return nil, nil
}

func (p *kugou) SongURL(song Song) (string, error){
	return "", nil
}

func (p *kugou) HotPlaylist(page int) (Playlists, error) {
	return nil, nil
}

func (p *kugou) PlaylistDetail(pl Playlist) (Songs, error){
	return nil, nil
}

func (p *kugou) Name() string {
	return "kugou"
}
