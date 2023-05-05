package main

import (
	"compress/gzip"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html"
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

	url := "https://www.zillow.com/homes/Boston,-MA_rb/"
	client := http.Client{Timeout: time.Second * 5}

	request, err := http.NewRequest(http.MethodGet, url, nil)

	request.Header.Set("User-Agent", "User")
	request.Header.Set("Accept", "*/")
	request.Header.Set("Accept-Encoding", "gzip")
	request.Header.Set("Accept-Language", "en-US,en;q=0.9")
	request.Header.Set("cache-control", "no-cache")
	request.Header.Set("pragma", "no-cache")
	request.Header.Set("sec-fetch-dest", "empty")
	request.Header.Set("sec-fetch-mode", "cors")
	request.Header.Set("sec-ch-ua", "'Google Chrome';v='113', 'Chromium';v='113', 'Not-A.Brand';v='24'")
	request.Header.Set("sec-fetch-site", "same-origin")
	request.Header.Set("sec-ch-ua-mobile", "?0")

	if err != nil {
		log.Fatal(err)
	}

	response, resErr := client.Do(request)
	if resErr != nil {
		log.Fatal(err)
	}

	reader, err := gzip.NewReader(response.Body)
	if err != nil {
		panic(err)
	}
	defer reader.Close()

	htmlString := html.NewTokenizer(reader)
	for {
		tt := htmlString.Next()
		if tt == html.ErrorToken {
			// ...
			fmt.Print("error")
			break
		}
		tagAttr := htmlString.Token().Attr
		if len(tagAttr) > 0 {
			if strings.Contains(tagAttr[0].Val, "cWiizR") {
				fmt.Println(tagAttr[0])
			}
		}
	}
}
