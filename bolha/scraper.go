package bolha

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	scrapers "github.com/seniorescobar/adnotifier-scrapers"

	log "github.com/sirupsen/logrus"
)

const (
	baseURL      = "https://www.bolha.com"
	itemSelector = `div.EntityList > ul.EntityList-items > li.EntityList-item--Regular`
)

type Scraper struct{}

func (s *Scraper) Scrape(ctx context.Context, url string) ([]*scrapers.Item, error) {
	r, err := fetch(ctx, url)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	return extractItems(r)
}

func extractItems(body io.Reader) ([]*scrapers.Item, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, err
	}

	items := make([]*scrapers.Item, 0)
	itemNodes := doc.Find(itemSelector)

	if len(itemNodes.Nodes) == 0 {
		log.Info("no item nodes found")
		return nil, nil
	}

	itemNodes.Each(func(i int, s *goquery.Selection) {
		path, ok := s.Attr("data-href")
		if !ok {
			log.Error(`attribute "data-href" does not exist`)
			return
		}

		item := scrapers.Item(baseURL + path)
		items = append(items, &item)
	})

	return items, nil
}

func fetch(ctx context.Context, url string) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid status code (%d)", res.StatusCode)
	}

	return res.Body, nil
}
