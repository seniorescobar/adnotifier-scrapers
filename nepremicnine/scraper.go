package nepremicnine

import (
	"fmt"
	"io"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	scrapers "github.com/seniorescobar/adnotifier-scrapers"

	log "github.com/sirupsen/logrus"
)

const (
	itemSelector = `div.seznam > div.oglas_container > div > a.slika`
)

type Scraper struct{}

func (s *Scraper) Scrape(url string) ([]*scrapers.Item, error) {
	r, err := fetch(url)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	return processItems(r)
}

func processItems(body io.ReadCloser) ([]*scrapers.Item, error) {
	log.Debug("processing items")

	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		log.Fatal(err)
	}

	items := make([]*scrapers.Item, 0)
	itemNodes := doc.Find(itemSelector)

	if len(itemNodes.Nodes) == 0 {
		log.Error("no item nodes found")
		return nil, nil
	}

	itemNodes.Each(func(i int, s *goquery.Selection) {
		path, ok := s.Attr("href")
		if !ok {
			log.Error(`attribute "href" does not exist`)
			return
		}

		item := scrapers.Item("https://www.nepremicnine.net/" + path)
		items = append(items, &item)
	})

	return items, nil
}

func fetch(url string) (io.ReadCloser, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("invalid status code (%d)", res.StatusCode)
	}

	return res.Body, nil
}
