package nepremicnine

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
		{"https://www.nepremicnine.net/abc", "https://www.nepremicnine.net/abc?s=16", "nepremicnine", nil},
		{"http://www.nepremicnine.net/abc", "https://www.nepremicnine.net/abc?s=16", "nepremicnine", nil},
		{"https://www.nepremicnine.net/abc?a=1", "https://www.nepremicnine.net/abc?a=1&s=16", "nepremicnine", nil},
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
