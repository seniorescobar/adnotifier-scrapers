Here's an example source code for scraping mobile.de.

    package main

    import (
        "context"
        "fmt"
        "net/url"

        "github.com/seniorescobar/adnotifier-scrapers/mobilede"
    )

    func main() {
        // parse url
        uo, _ := url.Parse(`https://suchen.mobile.de/fahrzeuge/search.html?dam=0&fr=2016&isSearchRequest=true&ms=17200;60&p=:30000&sfmr=false&vc=Car`)

        // transform url
        ut, _ := new(mobilede.Transformer).Transform(uo)

        // scrape
        items, _ := new(mobilede.Scraper).Scrape(context.TODO(), ut.String())

        // print
        for _, i := range items {
            fmt.Println(*i)
        }
    }
