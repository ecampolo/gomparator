package fetcher

import (
	"encoding/json"
	"github.com/emacampolo/gomparator/internal/http"
	"log"
	"net/url"
)

type Response struct {
	URL        *url.URL
	JSON       interface{}
	StatusCode int
}

func (r Response) IsOk() bool {
	return r.StatusCode == 200
}

type Fetcher interface {
	Fetch(host string, relPath string, headers map[string]string) (*Response, error)
}

func New() Fetcher {
	return fetcher{}
}

type fetcher struct{}

func (fetcher) Fetch(host string, relPath string, headers map[string]string) (*Response, error) {
	u, err := url.Parse(relPath)
	if err != nil {
		log.Fatal(err)
	}

	queryString := u.Query()
	u.RawQuery = queryString.Encode()

	base, err := url.Parse(host)
	if err != nil {
		log.Fatal(err)
	}

	u = base.ResolveReference(u)
	response, err := http.Get(u.String(), headers)
	if err != nil {
		return nil, err
	}

	return &Response{
		URL:        u,
		StatusCode: response.StatusCode,
		JSON:       toJson(response.Body),
	}, nil
}

func toJson(b []byte) interface{} {
	var j interface{}
	err := json.Unmarshal(b, &j)
	if err != nil {
		return nil
	}
	return j
}