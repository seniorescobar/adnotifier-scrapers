package scrapers

import "context"

type Scraper interface {
	Scrape(ctx context.Context, url string) ([]*Item, error)
}
