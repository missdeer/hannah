package bass

// #cgo LDFLAGS: -L/usr/local/lib/bass -lbass
import "C"
import (
	"regexp"
)

func pluginsPattern() ([]string, *regexp.Regexp) {
	return []string{
			"plugins",
			"bass/plugins",
			"bass/macOS/lib/plugins",
		},
		regexp.MustCompile(`^libbass[^\.]+.dylib$`)
}
