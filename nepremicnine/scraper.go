package nepremicnine

import (
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	scrapers "github.com/seniorescobar/adnotifier-scrapers"
)

const (
	// query selector for ads.
	itemSelector = `div.seznam > div.oglas_container > div > meta[itemprop=mainEntityOfPage]`

	// HTML attribute "content".
	attrContent = `content`
)

// headers added to the request.
var headers = map[string]string{
	"User-Agent":      "Mozilla/5.0 (X11; Linux x86_64; rv:84.0) Gecko/20100101 Firefox/84.0",
	"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
	"Accept-Language": "en-US,en;q=0.5",
	"Accept-Encoding": "gzip",
}

type (
	// Scraper represents an Nepremicnine Scraper.
	Scraper struct {
		httpClient httpClient
	}

	// httpClient abstracts a http client.
	httpClient interface {
		Do(*http.Request) (*http.Response, error)
	}

	// ErrAttributeDoesNotExist occurs when a HTML attribute does not exist.
	ErrAttributeDoesNotExist struct {
		attr string
	}
)

func (e ErrAttributeDoesNotExist) Error() string {
	return fmt.Sprintf("attribute %q does not exist", e.attr)
}

// NewScraper creates an instance of Nepremicnine Scraper.
func NewScraper(httpClient httpClient) *Scraper {
	return &Scraper{
		httpClient: httpClient,
	}
}

// Scrape scrapes the content of the given url and returns ads found.
func (s *Scraper) Scrape(ctx context.Context, url string) ([]scrapers.Item, error) {
	r, err := s.fetch(ctx, url)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	return s.processItems(r)
}

// Scrape scrapes the content of the given url and returns ads found.
func (s *Scraper) fetch(ctx context.Context, url string) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating a new request: %w", err)
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	res, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error doing a request: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid status code (%d)", res.StatusCode)
	}

	gzipR, err := gzip.NewReader(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error creating a gzip reader from response body: %w", err)
	}

	return gzipR, nil
}

func (s *Scraper) processItems(body io.Reader) ([]scrapers.Item, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, fmt.Errorf("error creating document from reader: %w", err)
	}

	var (
		items     = make([]scrapers.Item, 0)
		itemNodes = doc.Find(itemSelector)
	)

	if len(itemNodes.Nodes) == 0 {
		return nil, scrapers.ErrNoMatches
	}

	itemNodes.Each(func(_ int, sel *goquery.Selection) {
		path, ok := sel.Attr(attrContent)
		if !ok {
			err = ErrAttributeDoesNotExist{attrContent}
			return
		}

		item := scrapers.Item(path)
		items = append(items, &item)
	})
	if err != nil {
		return nil, fmt.Errorf("error iterating through item nodes: %w", err)
	}

	return items, nil
}
