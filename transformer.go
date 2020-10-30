package scrapers

import "net/url"

type Transformer interface {
	Transform(*url.URL) (*url.URL, error)
}
