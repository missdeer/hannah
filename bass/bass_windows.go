package bass

// #cgo LDFLAGS: -Lwindows/lib -lbass
import "C"
import (
	"io/ioutil"
	"log"
	"path"
	"path/filepath"
	"regexp"
)

var (
	fnRegexp = regexp.MustCompile(`^bass[^\.]+.dll`)
)

func LoadAllPlugins() {
	dirs := []string{"plugins", path.Join("bass", "plugins"), path.Join("bass", "windows", "plugins")}
	for _, dir := range dirs {
		fi, err := ioutil.ReadDir(dir)
		if err != nil {
			log.Println(err)
			continue
		}
		for _, f := range fi {
			if f.IsDir() {
				continue
			}
			if !fnRegexp.MatchString(f.Name()) {
				continue
			}
			fn := filepath.Join(dir, f.Name())
			PluginLoad(fn)
		}
	}
}
