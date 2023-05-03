package main

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/gocolly/colly"
)

type Houses struct {
	seller  string
	price   int
	address string
	details string
}

func newHouse(name string, price int, address string, details string) Houses {

	house := Houses{seller: name}
	house.price = price
	house.address = address
	house.details = details
	return house

}

func main() {
	var houseListings [18]Houses
	c := colly.NewCollector()

	counter := 0
	var priceRegex = regexp.MustCompile(`[^0-9]`)

	c.OnHTML("div.property-card-data", func(e *colly.HTMLElement) {

		var seller = e.ChildText("div.cWiizR")
		i, err := strconv.Atoi(priceRegex.ReplaceAllString(e.ChildText("div.bqsBln"), ""))
		_ = err
		var price = i
		var details = e.ChildText("div.gxlfal")
		var address = e.ChildText("div.jXNpbs")
		houseListings[counter] = newHouse(seller, price, address, details)
		counter = counter + 1
	})
	c.Visit("https://www.zillow.com/boston-ma")
	var averagePrice = 0
	for n := 0; n <= 8; n++ {
		averagePrice += houseListings[n].price
	}
	fmt.Print(averagePrice / 9)
}
