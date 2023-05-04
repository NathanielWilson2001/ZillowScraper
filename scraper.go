package main

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
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

	searchQueryState := "{'pagination':{},'usersSearchTerm':'Boston, MA','mapBounds':{'west':-71.30031054687501,'east':-70.79493945312501,'south':42.111529685321216,'north':42.514687824439775},'regionSelection':[{'regionId':44269,'regionType':6}],'isMapVisible':true,'filterState':{'sortSelection':{'value':'globalrelevanceex'},'isForSaleByAgent':{'value':false},'isForSaleByOwner':{'value':false},'isNewConstruction':{'value':false},'isForSaleForeclosure':{'value':false},'isComingSoon':{'value':false},'isAuction':{'value':false},'isRecentlySold':{'value':true},'isAllHomes':{'value':true}},'isListVisible':true,'mapZoom':11}"
	wants := "{'cat1':['mapResults']}"
	var stringSearch string = ""
	var stringWants string = ""
	errStringSearch := json.Unmarshal([]byte(searchQueryState), &stringSearch)
	errStringWants := json.Unmarshal([]byte(wants), &stringWants)
	if errStringSearch != nil && errStringWants != nil {
		fmt.Println("error: ", errStringSearch)
	}
	url := "https://www.zillow.com/homes/Boston,-MA_rb/"
	client := http.Client{Timeout: time.Second * 5}

	request, err := http.NewRequest(http.MethodGet, url, nil)

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
	// Read the decompressed response body
	body, err := io.ReadAll(reader)
	if err != nil {
		panic(err)
	}
	// Do something with the response body
	fmt.Println(string(body))
}
