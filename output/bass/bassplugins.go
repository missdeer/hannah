package bass

// #cgo windows,386 LDFLAGS: -Llib/windows/386 -lbass.dll -lbassasio.dll -lbassmix.dll -lbasswasapi.dll
// #cgo windows,amd64 LDFLAGS: -Llib/windows/amd64 -lbass.dll -lbassasio.dll -lbassmix.dll -lbasswasapi.dll
// #cgo darwin LDFLAGS: -Llib/darwin/amd64 -lbass
// #cgo linux,386 LDFLAGS: -Llib/linux/386 -lbass
// #cgo linux,amd64 LDFLAGS: -Llib/linux/amd64 -lbass
import "C"
import (
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
)

func pluginsPattern() ([]string, *regexp.Regexp) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return []string{
			dir,
			path.Join(dir, "output", "bass", "lib", runtime.GOOS, runtime.GOARCH, "plugins"),
		},
		regexp.MustCompile(`^(lib)?bass[0-9A-Za-z_]+.(dylib|dll|so)$`)
}
