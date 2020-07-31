// package bass link against bass.dll
// copy official bass.lib to libbass.dll.a
package bass

// #cgo LDFLAGS: -Llib/windows/x86 -lbass.dll
import "C"
import (
	"path"
	"regexp"
)

func pluginsPattern() ([]string, *regexp.Regexp) {
	return []string{
			".",
			"plugins",
			path.Join("bass", "plugins"),
			path.Join("output", "bass", "lib", "windows", "x86", "plugins"),
		},
		regexp.MustCompile(`^bass[0-9A-Za-z_]+.dll$`)
}
