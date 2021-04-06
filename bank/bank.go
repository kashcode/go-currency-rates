package bank

import (
	"log"
	"strconv"
	"strings"

	"example.com/currency-rates/entities"
	"github.com/mmcdole/gofeed"
)

type Bank struct {
	Url string
}

func parseCurrencyRates(items []*gofeed.Item) ([]entities.CurrencyRate, error) {
	log.Println("Parse rss feed into list of currency rates")

	var l []entities.CurrencyRate

	for _, item := range items {
		data := strings.Split(item.Description, " ")

		for i := 0; i < len(data); i += 2 {
			c := data[i]
			if c == "" {
				continue
			}

			n := i + 1
			r, err := strconv.ParseFloat(data[n], 32)
			if err != nil {
				log.Println(err.Error())
				continue
			}

			l = append(l, entities.CurrencyRate{
				Currency: c,
				Rate:     r,
				Date:     item.PublishedParsed,
			})
		}
	}

	return l, nil
}

// Get currency rates from RSS feed
func (b *Bank) GetCurrencyRates() ([]entities.CurrencyRate, error) {
	log.Println("Getting currency rates")

	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(b.Url)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	rates, err := parseCurrencyRates(feed.Items)

	return rates, err
}
