package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

func main() {
	c := colly.NewCollector()

	// Find and visit all links
	c.OnHTML("div.property-card-data", func(e *colly.HTMLElement) {
		fmt.Println("Name: " + e.ChildText("div.cWiizR") + "\nPrice: " + e.ChildText("div.bqsBln") + "\nDetails: " + e.ChildText("div.gxlfal") + "\nAddress: " + e.ChildText("address") + "\n" + e.ChildText("div.jXNpbs") + "\n")
	})

	c.Visit("https://www.zillow.com/boston-ma")
}
