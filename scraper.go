package main

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type Response struct {
	RequestId int `json:"requestId"`
	Cat1      struct {
		SearchResults struct {
			ListResults []struct {
				Zpid    string `json:"zpid"`
				Id      string `json:"id"`
				HdpData struct {
					HomeInfo struct {
						StreetAddress string  `json:"streetAddress"`
						Zipcode       string  `json:"zipcode"`
						City          string  `json:"city"`
						State         string  `json:"state"`
						Price         float64 `json:"price"`
						Bathrooms     float64 `json:"bathrooms"`
						Bedrooms      float64 `json:"bedrooms"`
						LivingArea    float64 `json:"livingArea"`
						HomeType      string  `json:"homeType"`
						HomeStatus    string  `json:"homeStatus"`
						Zestimate     int     `json:"zestimate"`
						RentZestimate int     `json:"rentZestimate"`
					} `json:"homeInfo"`
				} `json:"hdpData"`
			} `json:"listResults"`
		} `json:"searchResults"`
	} `json:"cat1"`
}

func responseToString(responseStructure Response) {

	b, err := json.MarshalIndent(responseStructure, "", "\t")
	_ = err

	fmt.Printf(string(b))
	fmt.Println()
	return
}

func calculate(responseStructure Response) {

	averagePriceSum := 0.0
	averageSquareFootSum := 0.0
	averageZestimate := 0.0
	averageZRentEstimate := 0.0
	var housingType [3]int

	for i := 0; i < len(responseStructure.Cat1.SearchResults.ListResults); i++ {

		averagePriceSum += float64(responseStructure.Cat1.SearchResults.ListResults[i].HdpData.HomeInfo.Price)
		averageSquareFootSum += float64(responseStructure.Cat1.SearchResults.ListResults[i].HdpData.HomeInfo.LivingArea)
		averageZRentEstimate += float64(responseStructure.Cat1.SearchResults.ListResults[i].HdpData.HomeInfo.RentZestimate)
		averageZestimate += float64(responseStructure.Cat1.SearchResults.ListResults[i].HdpData.HomeInfo.Zestimate)

		if strings.Contains(responseStructure.Cat1.SearchResults.ListResults[i].HdpData.HomeInfo.HomeType, "MULTI") {
			housingType[0]++
		}
		if strings.Contains(responseStructure.Cat1.SearchResults.ListResults[i].HdpData.HomeInfo.HomeType, "SINGLE") {
			housingType[1]++
		}
		if strings.Contains(responseStructure.Cat1.SearchResults.ListResults[i].HdpData.HomeInfo.HomeType, "CONDO") {
			housingType[2]++
		}

	}
	averagePriceSum = float64(averagePriceSum / float64(len(responseStructure.Cat1.SearchResults.ListResults)))
	averageSquareFootSum = float64(averageSquareFootSum / float64(len(responseStructure.Cat1.SearchResults.ListResults)))
	averageZRentEstimate = float64(averageZRentEstimate / float64(len(responseStructure.Cat1.SearchResults.ListResults)))
	averageZestimate = float64(averageZestimate / float64(len(responseStructure.Cat1.SearchResults.ListResults)))

	fmt.Printf("Count: %d \nAverage Price = %.2f \nAverage Square Foot = %.2f\nAverage Price per Square Foot = %.2f\nAverage Zestimate = %.2f\nAverage ZRentEstimate = %.2f\nHousing Type Breakdown:\n\tMultifamily: %d\n\tSingle Family: %d\n\tCondo: %d", len(responseStructure.Cat1.SearchResults.ListResults),
		averagePriceSum,
		averageSquareFootSum,
		averagePriceSum/averageSquareFootSum,
		averageZestimate,
		averageZestimate,
		housingType[0],
		housingType[1],
		housingType[2])
}

func main() {

	url := "https://www.zillow.com/search/GetSearchPageState.htm?searchQueryState=%7B%22mapBounds%22%3A%7B%22north%22%3A42.28936961396935%2C%22east%22%3A-70.90983246728517%2C%22south%22%3A41.86942883946931%2C%22west%22%3A-71.16732453271486%7D%2C%22isMapVisible%22%3Atrue%2C%22filterState%22%3A%7B%22isAllHomes%22%3A%7B%22value%22%3Atrue%7D%2C%22sortSelection%22%3A%7B%22value%22%3A%22globalrelevanceex%22%7D%2C%22isRecentlySold%22%3A%7B%22value%22%3Atrue%7D%2C%22isForSaleByAgent%22%3A%7B%22value%22%3Afalse%7D%2C%22isForSaleByOwner%22%3A%7B%22value%22%3Afalse%7D%2C%22isNewConstruction%22%3A%7B%22value%22%3Afalse%7D%2C%22isComingSoon%22%3A%7B%22value%22%3Afalse%7D%2C%22isAuction%22%3A%7B%22value%22%3Afalse%7D%2C%22isForSaleForeclosure%22%3A%7B%22value%22%3Afalse%7D%7D%2C%22isListVisible%22%3Atrue%2C%22regionSelection%22%3A%5B%7B%22regionId%22%3A58710%2C%22regionType%22%3A7%7D%5D%2C%22mapZoom%22%3A11%2C%22pagination%22%3A%7B%7D%7D&wants=%7B%22cat1%22:[%22listResults%22,%22mapResults%22]%7D&requestId=3"
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

	body, err := io.ReadAll(reader)

	responseStructure := &Response{}
	json.Unmarshal(body, &responseStructure)
	//fmt.Println(responseStructure)
	responseToString(*responseStructure)
	if err != nil {
		panic(err)
	}

	calculate(*responseStructure)
}
