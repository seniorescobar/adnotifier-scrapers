package avtonet

import (
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	scrapers "github.com/seniorescobar/adnotifier-scrapers"
)

var (
	// adPattern is a pattern for matching new ad urls.
	adPattern = regexp.MustCompile(`/Ads/details\.asp\?id=\d+"`)

	// headers are added to the request.
	headers = map[string]string{
		"User-Agent":                "Mozilla/5.0 (X11; Linux x86_64; rv:84.0) Gecko/20100101 Firefox/84.0",
		"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
		"Accept-Language":           "en-US,en;q=0.5",
		"Accept-Encoding":           "gzip",
		"DNT":                       "1",
		"Connection":                "keep-alive",
		"Cookie":                    "ogledov=; CookieConsent=-2",
		"Upgrade-Insecure-Requests": "1",
		"Cache-Control":             "max-age=0",
	}
)

// httpClient abstracts a http client.
type httpClient interface {
	Do(*http.Request) (*http.Response, error)
}

// Scraper represents an AvtoNet Scraper.
type Scraper struct {
	httpClient httpClient
}

// NewScraper creates an instance of AvtoNet Scraper.
func NewScraper(httpClient httpClient) *Scraper {
	return &Scraper{
		httpClient: httpClient,
	}
}

// Scrape scrapes the content of the given url and returns ads found.
func (s *Scraper) Scrape(ctx context.Context, url string) ([]scrapers.Item, error) {
	r, err := s.fetch(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("error fetching url: %w", err)
	}
	defer r.Close()

	items, err := s.processItems(r)
	if err != nil {
		return nil, fmt.Errorf("error processing items: %w", err)
	}

	return items, nil
}

func (s *Scraper) processItems(body io.ReadCloser) ([]scrapers.Item, error) {
	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	matches := adPattern.FindAll(bodyBytes, -1)
	if matches == nil {
		return nil, scrapers.ErrNoMatches
	}

	items := make([]scrapers.Item, len(matches))
	for i, m := range matches {
		url := string(m)
		url = strings.TrimRight(url, `"`)
		url = "https://www.avto.net" + url

		items[i] = scrapers.Item{
			URL: url,
		}
	}

	return items, nil
}

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
