package util

import (
	"bytes"
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

	"github.com/klauspost/compress/flate"
	"github.com/klauspost/compress/gzip"

	"github.com/andybalholm/brotli"
	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/proxy"
	"golang.org/x/net/publicsuffix"

	"github.com/missdeer/hannah/config"
)

var (
	errorNotIP    = errors.New("addr is not an IP")
	resolveResult = sync.Map{}
	once          = sync.Once{}
	globalClient  *http.Client
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

func createHttpClient() *http.Client {
	jar, _ := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})
	client := &http.Client{
		Transport: http.DefaultTransport,
		Jar:       jar,
		Timeout:   time.Duration(config.NetworkTimeout) * time.Second,
	}

	var localAddr net.Addr
	if config.NetworkInterface != "" {
		if ip := net.ParseIP(config.NetworkInterface); ip != nil {
			localAddr, _ = net.ResolveTCPAddr("tcp", config.NetworkInterface)
		}
		if i, err := net.InterfaceByName(config.NetworkInterface); err == nil {
			if addrs, err := i.Addrs(); err == nil {
				for _, addr := range addrs {
					ip, _, err := net.ParseCIDR(addr.String())
					if err == nil && ip != nil && !ip.IsLinkLocalUnicast() && !ip.IsLinkLocalMulticast() && !ip.IsLoopback() {
						localAddr = addr
						break
					}
				}
			}
		}
	}
	if localAddr != nil {
		ip, _, _ := net.ParseCIDR(localAddr.String())
		if ipaddr, err := net.ResolveIPAddr("ip", ip.String()); err == nil {
			client.Transport = &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					LocalAddr: &net.TCPAddr{IP: ipaddr.IP},
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
					DualStack: true,
				}).DialContext,
				ForceAttemptHTTP2:     true,
				MaxIdleConns:          100,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
			}
			net.DefaultResolver = &net.Resolver{
				PreferGo: true,
				Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
					d := net.Dialer{
						LocalAddr: &net.UDPAddr{IP: ipaddr.IP},
					}
					return d.DialContext(ctx, "udp", "119.29.29.29:53")
				},
			}
		}
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

func GetHttpClient() *http.Client {
	once.Do(func() { globalClient = createHttpClient() })
	return globalClient
}

func uncompressReader(r *http.Response) (io.ReadCloser, bool, error) {
	header := strings.ToLower(r.Header.Get("Content-Encoding"))
	switch header {
	case "":
		return r.Body, false, nil
	case "br":
		rc := brotli.NewReader(r.Body)
		if rc == nil {
			log.Println("creating brotli reader failed")
			return nil, false, errors.New("creating brotli reader failed")
		}
		return ioutil.NopCloser(rc), true, nil
	case "gzip":
		rc, err := gzip.NewReader(r.Body)
		if err != nil {
			log.Println("creating gzip reader failed:", err)
			return nil, false, err
		}
		return rc, true, nil
	case "deflate":
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("reading inflate failed:", err)
			return nil, false, err
		}
		rc := flate.NewReader(bytes.NewReader(content[2:]))
		if rc == nil {
			log.Println("creating deflate reader failed")
			return nil, false, errors.New("creating deflate reader failed")
		}
		return rc, true, nil
	}
	return nil, false, errors.New("unexpected encoding type")
}

func CopyHttpResponseBody(r *http.Response, w io.Writer) error {
	reader, needClose, err := uncompressReader(r)
	if err != nil {
		return err
	}
	_, err = io.Copy(w, reader)
	if needClose {
		reader.Close()
	}
	return err
}

func ReadHttpResponseBody(r *http.Response) (b []byte, err error) {
	reader, needClose, err := uncompressReader(r)
	if err != nil {
		return nil, err
	}
	b, err = ioutil.ReadAll(reader)
	if needClose {
		reader.Close()
	}
	return
}

func GetBaseURL(c *gin.Context) (baseURL string) {
	baseURL = config.BaseURL
	if baseURL == "" {
		scheme := c.Request.Header.Get("X-Forwarded-Proto")
		if scheme == "" {
			originURL := location.Get(c)
			scheme = originURL.Scheme
		}
		baseURL = fmt.Sprintf("%s://%s", scheme, c.Request.Host)
	}
	return
}
