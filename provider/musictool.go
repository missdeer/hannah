package provider

type musictool struct {
	provider string
}

func (p *musictool) SetProvider(provider string) {
	p.provider = provider
}

func (p *musictool) SearchSongs(keyword string, page int, limit int) (SearchResult, error) {
	return nil, ErrNotImplemented
}

func (p *musictool) ResolveSongURL(song Song) (Song, error) {
	return song, ErrNotImplemented
}

func (p *musictool) ResolveSongLyric(song Song) (Song, error){
	return song, ErrNotImplemented
}

func (p *musictool) HotPlaylist(page int, limit int) (Playlists, error) {
	return nil, ErrNotImplemented
}

func (p *musictool) PlaylistDetail(pl Playlist) (Songs, error) {
	return nil, ErrNotImplemented
}

func (p *musictool) ArtistSongs(id string) (res Songs, err error) {
	return nil, ErrNotImplemented
}

func (p *musictool) AlbumSongs(id string) (res Songs, err error) {
	return nil, ErrNotImplemented
}

func (p *musictool) Login() error {
	return  ErrNotImplemented
}

func (p *musictool) Name() string {
	return "musictool"
}
