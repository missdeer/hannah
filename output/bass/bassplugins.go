package bass

// #cgo windows,386 LDFLAGS: -Llib/windows/x86 -lbass.dll -lbassasio.dll -lbassmix.dll -lbasswasapi.dll
// #cgo windows,amd64 LDFLAGS: -Llib/windows -lbass.dll -lbassasio.dll -lbassmix.dll -lbasswasapi.dll
// #cgo darwin LDFLAGS: -Llib/macOS -lbass
import "C"
import (
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
)

func pluginsPattern() ([]string, *regexp.Regexp) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return []string{
			dir,
			path.Join(dir, "plugins"),
			path.Join(dir, "bass", "plugins"),
			path.Join(dir, "output", "bass", "lib", "macOS", "plugins"),
			path.Join(dir, "output", "bass", "lib", "windows", "plugins"),
			path.Join(dir, "output", "bass", "lib", "windows", "x86", "plugins"),
		},
		regexp.MustCompile(`^(lib)?bass[0-9A-Za-z_]+.(dylib|dll|so)$`)
}
