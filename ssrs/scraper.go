package ssrs

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	netURL "net/url"
	"strings"

	scrapers "github.com/seniorescobar/adnotifier-scrapers"
)

type Scraper struct{}

func (s *Scraper) Scrape(ctx context.Context, url string) ([]*scrapers.Item, error) {
	body, err := fetch(url)
	if err != nil {
		return nil, err
	}

	items, err := process(body)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func fetch(url string) (io.ReadCloser, error) {
	u, err := netURL.Parse(url)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(u.Fragment))
	if err != nil {
		return nil, err
	}

	// headers
	req.Header.Set("Host", " www.najem.stanovanjskisklad-rs.si")
	req.Header.Set("User-Agent", " Mozilla/5.0 (X11; Linux x86_64; rv:84.0) Gecko/20100101 Firefox/84.0")
	req.Header.Set("Accept", " */*")
	req.Header.Set("Accept-Language", " en-US,en;q=0.5")
	req.Header.Set("Accept-Encoding", " gzip, deflate")
	req.Header.Set("Content-Type", " application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("X-Requested-With", " XMLHttpRequest")
	req.Header.Set("Content-Length", " 84")
	req.Header.Set("Origin", " http://www.najem.stanovanjskisklad-rs.si")
	req.Header.Set("DNT", " 1")
	req.Header.Set("Connection", " keep-alive")
	req.Header.Set("Referer", " http://www.najem.stanovanjskisklad-rs.si/iskanje")
	req.Header.Set("Cookie", " cookie_choice=acc=1")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("invalid status code (%d)", res.StatusCode)
	}

	return gzip.NewReader(res.Body)
}

func process(body io.ReadCloser) ([]*scrapers.Item, error) {
	defer body.Close()

	type (
		item struct {
			ID string `json:"oznaka"`
		}
		response struct {
			List []item `json:"stanovanjeList"`
		}
	)

	var r response
	if err := json.NewDecoder(body).Decode(&r); err != nil {
		return nil, err
	}

	items := make([]*scrapers.Item, len(r.List))
	for i, it := range r.List {
		item := scrapers.Item("http://www.najem.stanovanjskisklad-rs.si/stanovanje/" + it.ID)
		items[i] = &item
	}

	return items, nil
}
