package scrapers

import (
	"context"
	"errors"
)

// ErrNoMatches occurs when there are no matches.
var ErrNoMatches = errors.New("no matches")

type Scraper interface {
	Scrape(ctx context.Context, url string) ([]Item, error)
}
