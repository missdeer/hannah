package util

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"context"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"time"

	"golang.org/x/net/proxy"
	"golang.org/x/net/publicsuffix"

	"github.com/missdeer/hannah/config"
)

type dialer struct {
	addr   string
	socks5 proxy.Dialer
}

func (d *dialer) DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	// TODO: golang.org/x/net/proxy need to add DialContext
	return d.Dial(network, addr)
}

func (d *dialer) Dial(network, addr string) (net.Conn, error) {
	var err error
	if d.socks5 == nil {
		d.socks5, err = proxy.SOCKS5("tcp", d.addr, nil, proxy.Direct)
		if err != nil {
			return nil, err
		}
	}
	return d.socks5.Dial(network, addr)
}

func socks5ProxyTransport(addr string) *http.Transport {
	d := &dialer{addr: addr}
	return &http.Transport{
		DialContext: d.DialContext,
		Dial:        d.Dial,
	}
}

func GetHttpClient() *http.Client {
	jar, _ := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})
	client := &http.Client{
		Transport: http.DefaultTransport,
		Jar:       jar,
		Timeout:   config.NetworkTimeout * time.Second,
	}
	httpProxy := os.Getenv("HTTP_PROXY")
	if config.HttpProxy != "" {
		httpProxy = config.HttpProxy
	}
	socks5Proxy := os.Getenv("SOCKS5_PROXY")
	if config.Socks5Proxy != "" {
		socks5Proxy = config.Socks5Proxy
	}
	if httpProxy != "" {
		if proxyUrl, err := url.Parse(httpProxy); err == nil {
			client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
		}
	} else if socks5Proxy != "" {
		client.Transport = socks5ProxyTransport(socks5Proxy)
	}
	return client
}

func ReadHttpResponseBody(r *http.Response) (b []byte, err error) {
	var (
		header string
		reader io.Reader
	)
	defer r.Body.Close()
	header = strings.ToLower(r.Header.Get("Content-Encoding"))
	switch header {
	case "":
		reader = r.Body
	case "gzip":
		if reader, err = gzip.NewReader(r.Body); err != nil {
			log.Fatalln("creating gzip reader failed:", err)
			return
		}
	case "deflate":
		content, e := ioutil.ReadAll(r.Body)
		if e != nil {
			log.Fatalln("reading inflate failed:", e)
			return []byte{}, e
		}

		if reader = flate.NewReader(bytes.NewReader(content[2:])); reader == nil {
			log.Fatalln("creating deflate reader failed")
			return []byte{}, errors.New("creating deflate reader failed")
		}
	}

	b, err = ioutil.ReadAll(reader)
	return
}
