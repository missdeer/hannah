package provider

type xiami struct {
}

func (p *xiami) Search(keyword string, page int, limit int) (SearchResult, error) {
	return nil, nil
}

func (p *xiami) HotPlaylist(page int) (Playlists, error) {
	return nil, nil
}

func (p *xiami) PlaylistDetail(pl Playlist) (Songs, error){
	return nil, nil
}