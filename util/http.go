package util

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/proxy"
	"golang.org/x/net/publicsuffix"

	"github.com/missdeer/hannah/config"
)

var (
	errorNotIP    = errors.New("addr is not an IP")
	resolveResult = sync.Map{}
)

func patchAddress(addr string) (string, error) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return addr, err
	}
	ip := net.ParseIP(host)
	if ip.To4() != nil || ip.To16() != nil {
		return addr, errorNotIP
	}
	// query from cache
	if rr, ok := resolveResult.Load(host); ok {
		ips := rr.([]string)
		if len(ips) > 0 {
			return net.JoinHostPort(ips[rand.Intn(len(ips))], port), nil
		}
	}
	// resolve it via http://119.29.29.29/d?dn=api.baidu.com
	client := GetHttpClient()
	req, err := http.NewRequest("GET", fmt.Sprintf("http://119.29.29.29/d?dn=%s", host), nil)
	if err != nil {
		log.Println(err)
		return addr, err
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return addr, err
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return addr, err
	}
	ips := string(content)
	ss := strings.Split(ips, ";")
	if len(ss) == 0 {
		return addr, err
	}
	resolveResult.Store(host, ss)
	return net.JoinHostPort(ss[0], port), nil
}

type dialer struct {
	addr   string
	socks5 proxy.Dialer
}

func (d *dialer) socks5DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	// TODO: golang.org/x/net/proxy need to add socks5DialContext
	return d.socks5Dial(network, addr)
}

func (d *dialer) socks5Dial(network, addr string) (net.Conn, error) {
	var err error
	if d.socks5 == nil {
		d.socks5, err = proxy.SOCKS5("tcp", d.addr, nil, proxy.Direct)
		if err != nil {
			return nil, err
		}
	}

	addr, _ = patchAddress(addr)
	return d.socks5.Dial(network, addr)
}

func socks5ProxyTransport(addr string) *http.Transport {
	d := &dialer{addr: addr}
	return &http.Transport{
		DialContext: d.socks5DialContext,
		Dial:        d.socks5Dial,
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
		if proxyURL, err := url.Parse(httpProxy); err == nil {
			transport := &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
			}
			transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
				addr, _ = patchAddress(addr)
				return transport.DialContext(ctx, network, addr)
			}
			transport.Dial = func(network, addr string) (net.Conn, error) {
				addr, _ = patchAddress(addr)
				return transport.Dial(network, addr)
			}
			client.Transport = transport
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
