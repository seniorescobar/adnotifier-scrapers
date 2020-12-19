package ssrs

import (
	"net/url"
)

type Transformer struct{}

func (t *Transformer) Transform(uo *url.URL) (*url.URL, error) {
	ut := *uo

	ut.Path = "/api/stanovanjefilter"

	return &ut, nil
}
