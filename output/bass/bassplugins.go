package bass

// #cgo windows,386 LDFLAGS: -Llib/windows/x86 -lbass.dll -lbassasio.dll -lbassmix.dll -lbasswasapi.dll
// #cgo windows,amd64 LDFLAGS: -Llib/windows -lbass.dll -lbassasio.dll -lbassmix.dll -lbasswasapi.dll
// #cgo darwin LDFLAGS: -Llib/macOS -lbass
import "C"
import (
	"path"
	"regexp"
)

func pluginsPattern() ([]string, *regexp.Regexp) {
	return []string{
			".",
			"plugins",
			"output/bass/lib/macOS/plugins",
			path.Join("bass", "plugins"),
			path.Join("output", "bass", "lib", "windows", "plugins"),
			path.Join("output", "bass", "lib", "windows", "x86", "plugins"),
		},
		regexp.MustCompile(`^(lib)?bass[0-9A-Za-z_]+.(dylib|dll|so)$`)
}
