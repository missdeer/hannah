package bass

// #cgo LDFLAGS: -Lwindows/lib -lbass
import "C"
import (
	"path"
	"regexp"
)

func pluginsPattern() ([]string, *regexp.Regexp) {
	return []string{
			"plugins",
			path.Join("bass", "plugins"),
			path.Join("bass", "windows", "plugins"),
		},
		regexp.MustCompile(`^bass[^\.]+.dll`)
}
