package config

import (
	"log"
	"os"
	"reflect"

	"gopkg.in/ini.v1"
)

var (
	Shuffle          bool
	Repeat           bool
	ByExternalPlayer bool
	Socks5Proxy      string
	HttpProxy        string
	Player           string
	Action           = "play"
	Provider         = "netease"
	Limit            = 25
	Page             = 1
	Engine           = "builtin"
	Mpg123           = true
	NetworkTimeout   = 30

	m = map[string]interface{}{
		"mpg123":             &Mpg123,
		"shuffle":            &Shuffle,
		"repeat":             &Repeat,
		"by-external-player": &ByExternalPlayer,
		"driver":             &AudioDriver,
		"action":             &Action,
		"provider":           &Provider,
		"socks5":             &Socks5Proxy,
		"http-proxy":         &HttpProxy,
		"player":             &Player,
		"engine":             &Engine,
		"limit":              &Limit,
		"page":               &Page,
		"network-timeout":    &NetworkTimeout,
	}
)

func LoadConfigurationFromFile(fn string) error {
	cfg, err := ini.Load(fn)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	for key, variable := range m {
		if !cfg.Section("").HasKey(key) {
			continue
		}
		switch v := reflect.Indirect(reflect.ValueOf(variable)); v.Kind() {
		case reflect.Int:
			if b, err := cfg.Section("").Key(key).Int(); err == nil {
				*(variable.(*int)) = b
			}
		case reflect.String:
			if b := cfg.Section("").Key(key).String(); b != "" {
				*(variable.(*string)) = b
			}
		case reflect.Bool:
			if b, err := cfg.Section("").Key(key).Bool(); err == nil {
				*(variable.(*bool)) = b
			}
		default:
			log.Fatalf("unsupported type:%s,%s\n", key, v.String())
		}
	}

	return nil
}
