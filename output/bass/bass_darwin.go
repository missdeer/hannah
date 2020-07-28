package bass

// #cgo LDFLAGS: -Llib/macOS -lbass
import "C"
import (
	"regexp"
)

func pluginsPattern() ([]string, *regexp.Regexp) {
	return []string{
			"plugins",
			"bass/plugins",
			"bass/lib/macOS/plugins",
		},
		regexp.MustCompile(`^libbass[^\.]+.dylib$`)
}
