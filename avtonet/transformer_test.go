package avtonet

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
		{"https://www.avto.net/abc", "https://www.avto.net/abc", "avtonet", nil},
		{"http://www.avto.net/abc", "https://www.avto.net/abc", "avtonet", nil},
		{"https://avto.net/abc", "https://www.avto.net/abc", "avtonet", nil},
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
