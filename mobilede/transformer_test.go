package mobilede

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
		{"https://suchen.mobile.de/fahrzeuge/search.html?dam=0&fr=2016%3A&isSearchRequest=true&ms=17200%3B60%3B%3B%3B&p=%3A30000&pw=KW%3AKW&s=Car&sfmr=false&vc=Car", "https://m.mobile.de/svc/s/?dam=0&fr=2016%3A&isSearchRequest=true&ms=17200%3B60%3B%3B%3B&od=down&p=%3A30000&pw=&s=Car&sb=doc&sfmr=false&vc=Car", "mobilede", nil},
		{"https://suchen.mobile.de/fahrzeuge/search.html?dam=0&fr=2016&isSearchRequest=true&ms=17200;60&p=:30000&sfmr=false&vc=Car", "https://m.mobile.de/svc/s/?dam=0&fr=2016&isSearchRequest=true&ms=17200%3B60&od=down&p=%3A30000&sb=doc&sfmr=false&vc=Car", "mobilede", nil},
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
