package avtonet_test

import (
	"net/url"
	"testing"

	"github.com/seniorescobar/adnotifier-scrapers/avtonet"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransfomer(t *testing.T) {
	tests := []struct {
		scenario                  string
		url, urlTransformed, site string
		err                       error
	}{
		{
			scenario:       "no change",
			url:            "https://www.avto.net/abc",
			urlTransformed: "https://www.avto.net/abc",
			site:           "avtonet",
		},
		{
			scenario:       "http to https",
			url:            "http://www.avto.net/abc",
			urlTransformed: "https://www.avto.net/abc",
			site:           "avtonet",
		},
		{
			scenario:       "prefix www",
			url:            "https://avto.net/abc",
			urlTransformed: "https://www.avto.net/abc",
			site:           "avtonet",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.scenario, func(t *testing.T) {
			urlParsed, err := url.Parse(tt.url)
			require.NoError(t, err)

			tr := avtonet.NewTransformer()
			ut, err := tr.Transform(urlParsed)
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.urlTransformed, ut.String())
		})
	}
}
