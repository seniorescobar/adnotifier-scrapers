package mobilede

import (
	"net/url"
	"strings"
)

type Transformer struct{}

func (t *Transformer) Transform(uo *url.URL) (*url.URL, error) {
	ut := *uo

	ut.Host = "m.mobile.de"
	ut.Path = "/svc/s/"

	// semicolons in query are interpreted as separators
	// replace them with "%3B"
	ut.RawQuery = strings.ReplaceAll(ut.RawQuery, ";", "%3B")

	query := ut.Query()

	query.Set("sb", "doc")
	query.Set("od", "down")

	// tmp bug fix
	if query.Get("pw") == "KW:KW" {
		query.Set("pw", "")
	}

	ut.RawQuery = query.Encode()

	return &ut, nil
}
