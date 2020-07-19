package input

import (
	"context"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"

	"golang.org/x/net/proxy"

	"github.com/missdeer/hannah/config"
)

func openLocalFile(filename string) (io.ReadCloser, error) {
	return os.Open(filename)
}

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

func socks5Proxy(addr string) *http.Transport {
	d := &dialer{addr: addr}
	return &http.Transport{
		DialContext: d.DialContext,
		Dial:        d.Dial,
	}
}
func openRemoteSource(u string) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	client := http.Client{}
	if config.HttpProxy != "" {
		if proxyUrl, err := url.Parse(config.HttpProxy); err == nil {
			client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
		}
	} else if config.Socks5Proxy != "" {
		client.Transport = socks5Proxy(config.Socks5Proxy)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func OpenSource(uri string) (io.ReadCloser, error) {
	if strings.HasPrefix(uri, "http://") || strings.HasPrefix(uri, "https://") {
		return openRemoteSource(uri)
	}
	return openLocalFile(uri)
}
