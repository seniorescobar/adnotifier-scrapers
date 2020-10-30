package bolha

import "net/url"

type Transformer struct{}

func (t *Transformer) Transform(uo *url.URL) (*url.URL, error) {
	ut := *uo

	ut.Scheme = "https"
	ut.Host = "www.bolha.com"

	qVals := ut.Query()
	qVals.Set("sort", "new")
	ut.RawQuery = qVals.Encode()

	return &ut, nil
}
