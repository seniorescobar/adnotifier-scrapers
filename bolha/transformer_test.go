package bolha

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransfomer(t *testing.T) {
	tr := new(Transformer)

	for _, tc := range []struct {
		url, urlTransformed, site string
		e                         error
	}{
		{"https://www.bolha.com/abc", "https://www.bolha.com/abc?sort=new", "bolha", nil},
		{"http://www.bolha.com/abc", "https://www.bolha.com/abc?sort=new", "bolha", nil},
		{"https://www.bolha.com/abc?a=1", "https://www.bolha.com/abc?a=1&sort=new", "bolha", nil},
	} {
		ut, err := tr.Transform(strToURL(tc.url))

		assert.Equal(t, tc.e, err)
		assert.Equal(t, tc.urlTransformed, ut.String())
	}
}

func strToURL(s string) *url.URL {
	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}
	return u
}
