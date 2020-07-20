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

func openLocalFile(filename string) (io.ReadCloser, error) {
	return os.Open(filename)
}

func openRemoteSource(u string) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	client := util.GetHttpClient()

	r, err := url.Parse(u)
	if err == nil {
		req.Header.Set("Referer", fmt.Sprintf("%s://%s", r.Scheme, r.Host))
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
