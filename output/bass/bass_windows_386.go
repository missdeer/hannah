// package bass link against bass.dll
// generate import lib for MinGW
// gendef.exe bass.dll
// dlltool.exe --no-leading-underscore -d bass.def -D bass.dll -l libbass.dll.a
package bass

// #cgo LDFLAGS: -Llib/windows/x86 -lbass.dll
import "C"
import (
	"path"
	"regexp"
)

func pluginsPattern() ([]string, *regexp.Regexp) {
	return []string{
			"plugins",
			path.Join("bass", "plugins"),
			path.Join("output", "bass", "lib", "windows", "x86", "plugins"),
		},
		regexp.MustCompile(`^bass[^\.]+.dll`)
}
