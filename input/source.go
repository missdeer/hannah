package input

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/missdeer/hannah/util"
)

var (
	httpClient = util.GetHttpClient()
)

func openLocalFile(filename string) (io.ReadCloser, error) {
	return os.Open(filename)
}

func openRemoteSource(u string) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	if r, err := url.Parse(u); err == nil {
		req.Header.Set("Referer", fmt.Sprintf("%s://%s", r.Scheme, r.Hostname()))
	}

	resp, err := httpClient.Do(req)
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
