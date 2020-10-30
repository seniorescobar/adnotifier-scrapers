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
)

type Scraper struct{}

func (s *Scraper) Scrape(ctx context.Context, url string) ([]*scrapers.Item, error) {
	r, err := fetch(ctx, url)
	if err != nil {
		return nil, err
	}

	return processItems(r)
}

func processItems(body io.ReadCloser) ([]*scrapers.Item, error) {
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

func fetch(ctx context.Context, url string) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept-Encoding", "gzip")

	client := &http.Client{
		Timeout: time.Second * 30,
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
