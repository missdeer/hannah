package bass

// #cgo LDFLAGS: -Llib/windows -lbass.dll
import "C"
import (
	"path"
	"regexp"
)

func pluginsPattern() ([]string, *regexp.Regexp) {
	return []string{
			"plugins",
			path.Join("bass", "plugins"),
			path.Join("output", "bass", "lib", "windows", "plugins"),
		},
		regexp.MustCompile(`^bass[^\.]+.dll`)
}
