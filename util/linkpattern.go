package util

import (
	"regexp"
)

var (
	playlistPatterns = map[*regexp.Regexp]string{
		regexp.MustCompile(`^https?:\/\/music\.163\.com\/#\/discover\/toplist\?id=(\d+)`):     "netease",
		regexp.MustCompile(`^https?:\/\/music\.163\.com\/#\/playlist\?id=(\d+)`):              "netease",
		regexp.MustCompile(`^https?:\/\/music\.163\.com\/#/my\/m\/music\/playlist\?id=(\d+)`): "netease",
		regexp.MustCompile(`^https?:\/\/www\.xiami\.com\/collect\/(\d+)`):                     "xiami",
		regexp.MustCompile(`^https?:\/\/y\.qq\.com\/n\/yqq\/playlist\/(\d+)\.html`):           "qq",
		regexp.MustCompile(`^https?:\/\/www\.kugou\.com\/yy\/special\/single\/(\d+)\.html`):   "kugou",
		regexp.MustCompile(`^https?:\/\/(www\.)?kuwo\.cn\/playlist_detail\/(\d+)`):            "kuwo",
		regexp.MustCompile(`^https?:\/\/music\.migu\.cn\/v3\/music\/playlist\/(\d+)`):         "migu",
	}
	songPatterns = map[*regexp.Regexp]string{
		regexp.MustCompile(`^https?:\/\/music\.163\.com\/#\/song\?id=(\d+)`):       "netease",
		regexp.MustCompile(`^https?:\/\/www\.xiami\.com\/song\/(\w+)`):             "xiami",
		regexp.MustCompile(`^https?:\/\/y\.qq\.com/n/yqq\/song\/(\w+)\.html`):      "qq",
		regexp.MustCompile(`^https?:\/\/www\.kugou\.com\/song\/#hash=([0-9A-F]+)`): "kugou",
		regexp.MustCompile(`^https?:\/\/(www\.)kuwo.cn\/play_detail\/(\d+)`):       "kuwo",
		regexp.MustCompile(`^https?:\/\/music\.migu\.cn\/v3\/music\/song\/(\d+)`):  "migu",
	}
	artistPatterns = map[*regexp.Regexp]string{
		regexp.MustCompile(`^https?:\/\/music\.163\.com\/weapi\/v1\/artist\/(\d+)`):                                       "netease",
		regexp.MustCompile(`^https?:\/\/music\.163\.com\/#\/artist\?id=(\d+)`):                                            "netease",
		regexp.MustCompile(`^https?:\/\/y\.qq\.com\/n\/yqq\/singer\/(\w+)\.html`):                                         "qq",
		regexp.MustCompile(`^https?:\/\/www\.xiami\.com\/artist\/(\w+)`):                                                  "xiami",
		regexp.MustCompile(`^https?:\/\/www\.xiami\.com\/list\?scene=artist&type=\w+&query={%22artistId%22:%22(\d+)%22}`): "xiami",
		regexp.MustCompile(`^https?:\/\/www\.xiami\.com\/list\?scene=artist&type=\w+&query={"artistId":"(\d+)"}`):         "xiami",
		regexp.MustCompile(`^https?:\/\/(www\.)?kuwo\.cn\/singer_detail\/(\d+)`):                                          "kuwo",
		regexp.MustCompile(`^https?:\/\/music\.migu\.cn\/v3\/music\/artist\/(\d+)`):                                       "migu",
	}
	albumPatterns = map[*regexp.Regexp]string{
		regexp.MustCompile(`^https?:\/\/music\.163\.com\/weapi\/v1\/album\/(\d+)`): "netease",
		regexp.MustCompile(`^https?:\/\/music\.163\.com\/#\/album\?id=(\d+)`):      "netease",
		regexp.MustCompile(`^https?:\/\/y\.qq\.com\/n\/yqq\/album\/(\w+)\.html`):   "qq",
		regexp.MustCompile(`^https?:\/\/www\.xiami\.com\/album\/(\w+)`):            "xiami",
		regexp.MustCompile(`^https?:\/\/(www\.)?kuwo\.cn\/album_detail\/(\d+)`):    "kuwo",
		regexp.MustCompile(`^https?:\/\/music\.migu\.cn\/v3\/music\/album\/(\d+)`): "migu",
	}
)

func patternMatch(u string, patterns map[*regexp.Regexp]string) (string, string, bool) {
	for pattern, providerName := range patterns {
		if pattern.MatchString(u) {
			ss := pattern.FindAllStringSubmatch(u, -1)
			if len(ss) == 1 && len(ss[0]) >= 2 {
				return ss[0][len(ss[0])-1], providerName, true
			}
		}
	}
	return "", "", false
}

func PlaylistMatch(u string) (string, string, bool) {
	return patternMatch(u, playlistPatterns)
}

func SingleSongMatch(u string) (string, string, bool) {
	return patternMatch(u, songPatterns)
}

func ArtistMatch(u string) (string, string, bool) {
	return patternMatch(u, artistPatterns)
}

func AlbumMatch(u string) (string, string, bool) {
	return patternMatch(u, albumPatterns)
}
