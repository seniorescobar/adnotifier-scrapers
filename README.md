Here's an example source code for scraping mobile.de.

    package main

    import (
        "context"
        "fmt"
        "net/url"

        "github.com/seniorescobar/adnotifier-scrapers/mobilede"
    )

    const (
        u1 = `https://suchen.mobile.de/fahrzeuge/search.html?dam=0&fr=2016&isSearchRequest=true&ms=17200;60&p=:30000&sfmr=false&vc=Car`
    )

    func main() {
        // parse url
        uo, _ := url.Parse(u1)

        // transform url
        t := new(mobilede.Transformer)
        ut, _ := t.Transform(uo)

        // scrape
        s := new(mobilede.Scraper)
        items, _ := s.Scrape(context.TODO(), ut.String())

        // print
        for _, i := range items {
            fmt.Println(*i)
        }
    }