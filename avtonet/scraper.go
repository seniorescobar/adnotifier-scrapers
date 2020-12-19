package avtonet

import (
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"

	scrapers "github.com/seniorescobar/adnotifier-scrapers"

	log "github.com/sirupsen/logrus"
)

type Scraper struct{}

func (s *Scraper) Scrape(ctx context.Context, url string) ([]*scrapers.Item, error) {
	r, err := s.fetch(ctx, url)
	if err != nil {
		return nil, err
	}

	return s.processItems(r)
}

func (s *Scraper) processItems(body io.ReadCloser) ([]*scrapers.Item, error) {
	bodyBytes, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}
	body.Close()

	r := regexp.MustCompile(`/Ads/details\.asp\?id=\d+`)

	matches := r.FindAll(bodyBytes, -1)

	if matches == nil {
		return nil, errors.New("no matches")
	}

	locMap := make(map[string]struct{})
	for _, loc := range matches {
		locMap[string(loc)] = struct{}{}
	}

	items := make([]*scrapers.Item, 0)
	for loc := range locMap {
		url := "https://www.avto.net" + loc

		item := scrapers.Item(url)
		items = append(items, &item)
	}

	return items, nil
}

func (s *Scraper) fetch(ctx context.Context, url string) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", " Mozilla/5.0 (X11; Linux x86_64; rv:84.0) Gecko/20100101 Firefox/84.0")
	req.Header.Set("Accept", " text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", " en-US,en;q=0.5")
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("DNT", " 1")
	req.Header.Set("Connection", " keep-alive")
	req.Header.Set("Cookie", " ogledov=; CookieConsent=-2")
	req.Header.Set("Upgrade-Insecure-Requests", " 1")
	req.Header.Set("Cache-Control", " max-age=0")

	proxyURL, _ := http.ProxyFromEnvironment(req)
	if proxyURL != nil {
		log.WithField("url", proxyURL).Debug("using proxy")
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("invalid status code (%d)", res.StatusCode)
	}

	return gzip.NewReader(res.Body)
}
