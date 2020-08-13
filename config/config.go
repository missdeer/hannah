package config

import (
	"log"
	"os"
	"path/filepath"
	"reflect"

	"gopkg.in/ini.v1"
)

var (
	Shuffle             bool
	Repeat              bool
	ByExternalPlayer    bool
	CacheEnabled        bool
	RedirectURL         bool
	AutoRedirectURL     bool
	ReverseProxyEnabled bool
	Socks5Proxy         string
	HttpProxy           string
	Player              string
	DownloadDir         string
	M3UFileName         string
	BaseURL             string
	CacheAddr           = "127.0.0.1:6379"
	ReverseProxy        = "127.0.0.1:8123"
	Action              = "play"
	Provider            = "netease"
	Limit               = 35
	Page                = 1
	Engine              = "bass"
	NetworkTimeout      = 60

	m = map[string]interface{}{
		"baseurl":               &BaseURL,
		"cache":                 &CacheEnabled,
		"cache-addr":            &CacheAddr,
		"redirect":              &RedirectURL,
		"auto-redirect":         &AutoRedirectURL,
		"reverse-proxy-enabled": &ReverseProxyEnabled,
		"reverse-proxy":         &ReverseProxy,
		"dir":                   &DownloadDir,
		"m3u":                   &M3UFileName,
		"shuffle":               &Shuffle,
		"repeat":                &Repeat,
		"by-external-player":    &ByExternalPlayer,
		"driver":                &AudioDriver,
		"action":                &Action,
		"provider":              &Provider,
		"socks5":                &Socks5Proxy,
		"http-proxy":            &HttpProxy,
		"player":                &Player,
		"engine":                &Engine,
		"limit":                 &Limit,
		"page":                  &Page,
		"network-timeout":       &NetworkTimeout,
	}
)

func init() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	DownloadDir = filepath.Join(pwd, "download")
	M3UFileName = filepath.Join(pwd, "hannah.m3u")
}

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
