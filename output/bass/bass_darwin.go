package bass

// #cgo LDFLAGS: -LmacOS/lib -lbass
import "C"
import (
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
)

var (
	fnRegexp = regexp.MustCompile(`^libbass[^\.]+.dylib$`)
)

func LoadAllPlugins() {
	dirs := []string{"plugins", "bass/plugins", "bass/macOS/plugins"}
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
