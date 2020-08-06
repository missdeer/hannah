package provider

type musictool struct {
	provider string
}

func (p *musictool) SetProvider(provider string) {
	p.provider = provider
}

func (p *musictool) Search(keyword string, page int, limit int) (SearchResult, error) {
	return nil, nil
}

func (p *musictool) ResolveSongURL(song Song) (Song, error) {
	return song, nil
}

func (p *musictool) ResolveSongLyric(song Song) (Song, error){
	return song, nil
}

func (p *musictool) HotPlaylist(page int, limit int) (Playlists, error) {
	return nil, nil
}

func (p *musictool) PlaylistDetail(pl Playlist) (Songs, error) {
	return nil, nil
}

func (p *musictool) Name() string {
	return "musictool"
}
