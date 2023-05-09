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

	fmt.Printf("Count: %d \nAverage Price = %.2f \nAverage Square Foot = %.2f\nAverage Price per Square Foot = %.2f\nAverage Zestimate = %.2f\nAverage ZRentEstimate = %.2f\nHousing Type Breakdown:\n\tMultifamily: %d\n\tSingle Family: %d\n\tCondo: %d\n", len(responseStructure.Cat1.SearchResults.ListResults),
		averagePriceSum,
		averageSquareFootSum,
		averagePriceSum/averageSquareFootSum,
		averageZestimate,
		averageZestimate,
		housingType[0],
		housingType[1],
		housingType[2])

}
func makeRequest(north float64, south float64, east float64, west float64) {
	url := "https://www.zillow.com/search/GetSearchPageState.htm?"
	client := http.Client{Timeout: time.Second * 5}

	request, err := http.NewRequest(http.MethodGet, url, nil)
	query := request.URL.Query()

	type urlParamters struct {
		MapBounds struct {
			North float64 `json:"north"`
			South float64 `json:"south"`
			East  float64 `json:"east"`
			West  float64 `json:"west"`
		} `json:"mapBounds"`
		MapZoom      int  `json:"mapZoom"`
		IsMapVisible bool `json:"isMapVisible"`
		FilterState  struct {
			IsAllHomes struct {
				Value bool `json:"value"`
			} `json:"isAllHomes"`
			SortSelection struct {
				Value string `json:"value"`
			} `json:"sortSelection"`
		} `json:"filterState"`
		IsListVisible bool `json:"isListVisible"`
	}

	type urlWants struct {
		Cat1          [1]string `json:"cat1"`
		Cat2          [1]string `json:"cat2"`
		RegionResults [1]string `json:"regionResults"`
	}

	params := &urlParamters{}

	params.MapBounds.North = north
	params.MapBounds.South = south
	params.MapBounds.East = east
	params.MapBounds.West = west
	params.MapZoom = 9
	params.IsMapVisible = false
	params.FilterState.IsAllHomes.Value = true
	params.FilterState.SortSelection.Value = "globalrelevanceex"
	params.IsListVisible = true
	wants := &urlWants{}
	wants.Cat1[0] = "listResults"
	wants.Cat2[0] = "total"
	wants.RegionResults[0] = "regionResults"

	b, _ := json.Marshal(params)
	query.Add("searchQueryState", string(b))

	b, _ = json.Marshal(wants)
	query.Add("wants", string(b))

	request.URL.RawQuery = query.Encode()

	fmt.Println(request.URL.String())

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

	//responseToString(*responseStructure)
	if err != nil {
		panic(err)
	}

	calculate(*responseStructure)
}
func main() {

	makeRequest(42.28936961396935, 41.86942883946931, -70.90983246728517, -71.16732453271486)
	makeRequest(43.1847127353123, 41.512239125025815, -70.65074389648439, -71.68071215820314)

}
