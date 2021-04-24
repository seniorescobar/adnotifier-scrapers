package avtonet_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"testing"

	scrapers "github.com/seniorescobar/adnotifier-scrapers"
	"github.com/seniorescobar/adnotifier-scrapers/avtonet"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

const testResources = "./test_resources"

var errFailed = errors.New("failed")

type httpClient struct {
	mock.Mock
}

func (m *httpClient) Do(r *http.Request) (*http.Response, error) {
	args := m.Called(r)

	return args.Get(0).(*http.Response), args.Error(1)
}

type ScraperSuite struct {
	suite.Suite

	httpClient *httpClient
	scraper    *avtonet.Scraper
}

func (s *ScraperSuite) SetupTest() {
	s.httpClient = new(httpClient)
	s.scraper = avtonet.NewScraper(s.httpClient)
}

func TestMigrator(t *testing.T) {
	suite.Run(t, new(ScraperSuite))
}

func (s *ScraperSuite) TestScrape_Success() {
	var (
		url      = "https://www.avto.net/Ads/results.asp"
		expItems = []scrapers.Item{
			{URL: "https://www.avto.net/Ads/details.asp?id=16290198"},
			{URL: "https://www.avto.net/Ads/details.asp?id=15711734"},
			{URL: "https://www.avto.net/Ads/details.asp?id=16296258"},
		}
	)

	f, err := os.Open(path.Join(testResources, "benz-c.html.gz"))
	s.Require().NoError(err)
	defer f.Close()

	s.httpClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return assert.Equal(s.T(), http.MethodGet, req.Method) &&
			assert.Equal(s.T(), url, req.URL.String()) &&

			// assert headers
			assert.Equal(s.T(), "Mozilla/5.0 (X11; Linux x86_64; rv:84.0) Gecko/20100101 Firefox/84.0", req.Header.Get("User-Agent")) &&
			assert.Equal(s.T(), "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8", req.Header.Get("Accept")) &&
			assert.Equal(s.T(), "en-US,en;q=0.5", req.Header.Get("Accept-Language")) &&
			assert.Equal(s.T(), "gzip", req.Header.Get("Accept-Encoding")) &&
			assert.Equal(s.T(), "1", req.Header.Get("DNT")) &&
			assert.Equal(s.T(), "keep-alive", req.Header.Get("Connection")) &&
			assert.Equal(s.T(), "ogledov=; CookieConsent=-2", req.Header.Get("Cookie")) &&
			assert.Equal(s.T(), "1", req.Header.Get("Upgrade-Insecure-Requests")) &&
			assert.Equal(s.T(), "max-age=0", req.Header.Get("Cache-Control"))
	})).Return(
		response(http.StatusOK, f),
		nil,
	)

	items, err := s.scraper.Scrape(context.Background(), url)
	s.NoError(err)
	s.ElementsMatch(expItems, items)
}

func (s *ScraperSuite) TestScrape_ErrNewRequest() {
	items, err := s.scraper.Scrape(nil, "")
	s.Error(err)
	s.Nil(items)
}

func (s *ScraperSuite) TestScrape_ErrDo() {
	s.httpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(&http.Response{}, errFailed)

	items, err := s.scraper.Scrape(context.Background(), "")
	s.Equal(errFailed, unwrap(err))
	s.Nil(items)
}

func (s *ScraperSuite) TestScrape_ErrHttpStatus() {
	s.httpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(
		&http.Response{
			StatusCode: http.StatusInternalServerError,
		},
		nil,
	)

	items, err := s.scraper.Scrape(context.Background(), "")
	s.EqualError(unwrap(err), "invalid status code (500)")
	s.Empty(items)
}

func (s *ScraperSuite) TestScrape_ErrGzip() {
	s.httpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(
		response(http.StatusOK, strings.NewReader("abc")),
		nil,
	)

	items, err := s.scraper.Scrape(context.Background(), "")
	s.Equal(io.ErrUnexpectedEOF, unwrap(err))
	s.Nil(items)
}

func response(statusCode int, r io.Reader) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(r),
	}
}

func (s *ScraperSuite) TestScrape_ErrNoMatches() {
	f, err := os.Open(path.Join(testResources, "no-matches.html.gz"))
	s.Require().NoError(err)
	defer f.Close()

	s.httpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(
		response(http.StatusOK, f),
		nil,
	)

	items, err := s.scraper.Scrape(context.Background(), "")
	s.Equal(scrapers.ErrNoMatches, unwrap(err))
	s.Nil(items)
}

func unwrap(err error) error {
	for errors.Unwrap(err) != nil {
		err = errors.Unwrap(err)
	}

	return err
}
