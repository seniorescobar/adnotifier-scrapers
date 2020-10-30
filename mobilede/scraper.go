package mobilede

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	scrapers "github.com/seniorescobar/adnotifier-scrapers"
)

type Scraper struct{}

func (s *Scraper) Scrape(ctx context.Context, url string) ([]*scrapers.Item, error) {
	body, err := fetch(url)
	if err != nil {
		return nil, err
	}

	items, err := process(body)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func fetch(url string) (io.ReadCloser, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Host", " m.mobile.de")
	req.Header.Add("User-Agent", " Mozilla/5.0 (iPhone; CPU iPhone OS 12_0 like Mac OS X) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/12.0 Mobile/15A372 Safari/604.1")
	req.Header.Add("Accept", " application/json")
	req.Header.Add("Accept-Language", " en-US,en;q=0.5")
	req.Header.Add("Accept-Encoding", " gzip, deflate, br")
	req.Header.Add("X-Requested-With", " XMLHttpRequest")
	req.Header.Add("X-Mobile-Client", " de.mobile.mportal.app/606/4e02613b-2430-4dcf-bfac-f274f7181b0d")
	req.Header.Add("X-Mobile-Api-Version", " 7")
	req.Header.Add("X-Override-Language", " en")
	req.Header.Add("X-Mobile-Feature-Variant", " cd175-equal-choice:NO_VARIANT,cd199-mweb-vip-integration-fin:NO_VARIANT,cd173-ovk-interscroller-test:NO_VARIANT,web2app-banner-overlay:NO_VARIANT")
	req.Header.Add("cache-control", " no-cache")
	req.Header.Add("X-Mobile-Vi", " eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJjaWQiOiJlOTViOGI2MS1iMzFiLTQ3NDgtYmIxZS03MTZiMjk5MTZhMGQiLCJhdWQiOltdLCJpYXQiOjE1OTA5NDE3MzYsImRudCI6dHJ1ZX0.seazK2ywS1dxWV7aDZxFuWE0XJ6gK8brBSU206Hso_w")
	req.Header.Add("DNT", " 1")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("invalid status code (%d)", res.StatusCode)
	}

	return gzip.NewReader(res.Body)
}

func process(body io.ReadCloser) ([]*scrapers.Item, error) {
	type (
		item struct {
			ID  int    `json:"id"`
			URL string `json:"url"`
		}
		response struct {
			Items []item `json:"items"`
		}
	)

	bodyBytes, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}
	body.Close()

	var r response
	if err := json.Unmarshal(bodyBytes, &r); err != nil {
		return nil, err
	}

	items := make([]*scrapers.Item, len(r.Items))
	for i, it := range r.Items {
		item := scrapers.Item(it.URL)
		items[i] = &item
	}

	return items, nil
}
