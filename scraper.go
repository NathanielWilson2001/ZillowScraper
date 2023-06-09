package main

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
				ImgSrc  string `json:"imgSrc"`
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

type SendData struct {
	CityDataList []ListItem `json:"cityDataList"`
}

type ListItem struct {
	Location     string       `json:"location"`
	Calculations Calculations `json:"data"`
}

type Calculations struct {
	RunningTotalEntries       int     `json:"runningTotalEntries"`
	AveragePriceSum           float64 `json:"averagePriceSum"`
	AverageSquareFootSum      float64 `json:"averageSquareFootSum"`
	AveragePricePerSquareFoot float64 `json:"averagePricePerSquareFoot"`
	AverageZestimate          float64 `json:"averageZestimate"`
	AverageRentZestimate      float64 `json:"averageRentZestimate"`
	MultiFamily               int     `json:"multiFamily"`
	SingleFamily              int     `json:"singleFamily"`
	Condo                     int     `json:"condo"`
	Date                      string  `json:"timeOfCompletion"`
}

type Data struct {
	NorthEast struct {
		Boston struct {
			Calculations []Calculations `json:"data"`
		} `json:"Boston"`
		Brooklyn struct {
			Calculations []Calculations `json:"data"`
		} `json:"Brooklyn"`
		Philadelphia struct {
			Calculations []Calculations `json:"data"`
		} `json:"Philadelphia"`
		WashingtonDC struct {
			Calculations []Calculations `json:"data"`
		} `json:"WashingtonDC"`
		Baltimore struct {
			Calculations []Calculations `json:"data"`
		} `json:"Baltimore"`
	} `json:"NorthEast"`
	South struct {
		Dallas struct {
			Calculations []Calculations `json:"data"`
		} `json:"Dallas"`
		Jacksonville struct {
			Calculations []Calculations `json:"data"`
		} `json:"Jacksonville"`
		Charlotte struct {
			Calculations []Calculations `json:"data"`
		} `json:"Charlotte"`
		Nashville struct {
			Calculations []Calculations `json:"data"`
		} `json:"Nashville"`
		Memphis struct {
			Calculations []Calculations `json:"data"`
		} `json:"Memphis"`
	} `json:"South"`
	WestCoast struct {
		LosAngeles struct {
			Calculations []Calculations `json:"data"`
		} `json:"LosAngeles"`
		SanFrancisco struct {
			Calculations []Calculations `json:"data"`
		} `json:"SanFrancisco"`
		Seattle struct {
			Calculations []Calculations `json:"data"`
		} `json:"Seattle"`
		Portland struct {
			Calculations []Calculations `json:"data"`
		} `json:"Portland"`
		SanDiego struct {
			Calculations []Calculations `json:"data"`
		} `json:"SanDiego"`
		Anchorage struct {
			Calculations []Calculations `json:"data"`
		} `json:"Anchorage"`
	} `json:"WestCoast"`
}

/*
*	The purpose of responseToString(responseStructure Response) is to take in a JSON struct response from the request
*	And output the results to console in a more clear, easy to read JSON repsonse
 */
func responseToString(responseStructure Response) {

	b, err := json.MarshalIndent(responseStructure, "", "\t")
	_ = err

	fmt.Printf(string(b))
	fmt.Println()
	return
}

/*
*	The purpose of calculate(responseStructure Response) is to take in a JSON struct response from the request
*	and calculate the information from the data set by looping through all entries:
*	Gets the Average Price, Average Sq Footage, Average Zestimate, Average Rent Zestimate
*	Outputs them to the console as a String
 */
func calculate(responses []Response, numberOfPages int) Calculations {

	// Initilize Base values to 0
	averagePriceSum := 0.0
	averageSquareFootSum := 0.0
	averageZestimate := 0.0
	averageZRentEstimate := 0.0
	var housingType [3]int
	runningTotalEntries := 0
	for pageNumber := 0; pageNumber < numberOfPages; pageNumber++ {

		responseStructure := responses[pageNumber]
		numOfEntries := len(responseStructure.Cat1.SearchResults.ListResults)

		runningTotalEntries += numOfEntries
		// Loop through all entries in result
		for i := 0; i < numOfEntries; i++ {

			currentEntry := responseStructure.Cat1.SearchResults.ListResults[i].HdpData.HomeInfo

			averagePriceSum += float64(currentEntry.Price)
			averageSquareFootSum += float64(currentEntry.LivingArea)
			averageZRentEstimate += float64(currentEntry.RentZestimate)
			averageZestimate += float64(currentEntry.Zestimate)

			if strings.Contains(currentEntry.HomeType, "MULTI") {
				housingType[0]++
			}
			if strings.Contains(currentEntry.HomeType, "SINGLE") {
				housingType[1]++
			}
			if strings.Contains(currentEntry.HomeType, "CONDO") {
				housingType[2]++
			}

		}
	}
	// Calculate the average
	averagePriceSum = float64(averagePriceSum / float64(runningTotalEntries))
	averageSquareFootSum = float64(averageSquareFootSum / float64(runningTotalEntries))
	averageZRentEstimate = float64(averageZRentEstimate / float64(runningTotalEntries))
	averageZestimate = float64(averageZestimate / float64(runningTotalEntries))

	currentTime := fmt.Sprint(time.Now().Format("2006-01-02"))
	calculations := Calculations{
		RunningTotalEntries:       runningTotalEntries,
		AveragePriceSum:           averagePriceSum,
		AverageSquareFootSum:      averageSquareFootSum,
		AveragePricePerSquareFoot: averagePriceSum / averageSquareFootSum,
		AverageZestimate:          averageZestimate,
		AverageRentZestimate:      averageZRentEstimate,
		MultiFamily:               housingType[0],
		SingleFamily:              housingType[1],
		Condo:                     housingType[2],
		Date:                      currentTime,
	}
	return calculations
}

/*
*
*
*
 */
func calculateRegion(w http.ResponseWriter, r *http.Request) {

}

/*
*	The purpose of makeRequest is to take in 4 GPS coordinate values that are used to set a boundary for a given search location,
*	as well as the Page number, and generate a request URL to Zillow. The request is then sent and the response is stored
*	in the response object. Response calculate is called to output the calculations from the given data set.
 */
func makeRequest(north float64, south float64, east float64, west float64, numPages int) Calculations {

	time.Sleep(15 * time.Second)
	var responses []Response
	for pageNumber := 1; pageNumber <= numPages; pageNumber++ {
		url := "https://www.zillow.com/search/GetSearchPageState.htm?requestID=2?"
		client := http.Client{Timeout: time.Second * 5}

		request, err := http.NewRequest(http.MethodGet, url, nil)
		query := request.URL.Query()

		// URL parameter struct created to form the 'searchQueryState' data necessary for the request.
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
			Pagination    struct {
				CurrentPage int `json:"currentPage"`
			} `json:"pagination"`
		}

		// URL 'wants' parameter that attaches to the end of the request link.
		type urlWants struct {
			Cat1          [1]string `json:"cat1"`
			Cat2          [1]string `json:"cat2"`
			RegionResults [1]string `json:"regionResults"`
		}

		params := &urlParamters{}

		// Set values for URL Parameters
		params.MapBounds.North = north
		params.MapBounds.South = south
		params.MapBounds.East = east
		params.MapBounds.West = west
		params.MapZoom = 9
		params.IsMapVisible = false
		params.FilterState.IsAllHomes.Value = true
		params.FilterState.SortSelection.Value = "globalrelevanceex"
		params.IsListVisible = true
		params.Pagination.CurrentPage = pageNumber

		// Set want values
		wants := &urlWants{}
		wants.Cat1[0] = "listResults"
		wants.Cat2[0] = "total"
		wants.RegionResults[0] = "regionResults"

		// Attach structs to the query and encode
		b, _ := json.Marshal(params)
		query.Add("searchQueryState", string(b))

		b, _ = json.Marshal(wants)
		query.Add("wants", string(b))

		request.URL.RawQuery = query.Encode()

		// Set default header values and encoding type
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

		fmt.Println(request.URL)
		if err != nil {
			log.Fatal(err)
		}

		// Perform request and decode response
		response, resErr := client.Do(request)
		if resErr != nil {
			log.Fatal(err)
		}
		reader, err := gzip.NewReader(response.Body)
		if err != nil {
			panic(err)
		}
		defer reader.Close()

		// Store response into a Response struct
		body, err := io.ReadAll(reader)
		responseStructure := &Response{}
		json.Unmarshal(body, &responseStructure)

		if err != nil {
			panic(err)
		}

		responses = append(responses, *responseStructure)
		//	Print Out Each request results in json
		//responseToString(*responseStructure)
	}

	return calculate(responses, numPages)

}

/* Send Fetch Data Request */
func fetchData(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	jsonFile, _ := os.Open("src/dataStorage.json")
	jsonData, _ := ioutil.ReadAll(jsonFile)
	dataStructure := &Data{}

	json.Unmarshal(jsonData, &dataStructure)
	responseToWeb := &SendData{}

	var listItem = &ListItem{}
	switch string(body) {
	case "NorthEast":
		listItem.Location = "Boston"
		listItem.Calculations = dataStructure.NorthEast.Boston.Calculations[len(dataStructure.NorthEast.Boston.Calculations)-1]
		responseToWeb.CityDataList = append(responseToWeb.CityDataList, *listItem)

		listItem.Location = "Brooklyn"
		listItem.Calculations = dataStructure.NorthEast.Brooklyn.Calculations[len(dataStructure.NorthEast.Brooklyn.Calculations)-1]
		responseToWeb.CityDataList = append(responseToWeb.CityDataList, *listItem)

		listItem.Location = "Philadelphia"
		listItem.Calculations = dataStructure.NorthEast.Philadelphia.Calculations[len(dataStructure.NorthEast.Philadelphia.Calculations)-1]
		responseToWeb.CityDataList = append(responseToWeb.CityDataList, *listItem)

		listItem.Location = "Washington DC"
		listItem.Calculations = dataStructure.NorthEast.WashingtonDC.Calculations[len(dataStructure.NorthEast.WashingtonDC.Calculations)-1]
		responseToWeb.CityDataList = append(responseToWeb.CityDataList, *listItem)

		listItem.Location = "Baltimore"
		listItem.Calculations = dataStructure.NorthEast.Baltimore.Calculations[len(dataStructure.NorthEast.Baltimore.Calculations)-1]
		responseToWeb.CityDataList = append(responseToWeb.CityDataList, *listItem)

	case "South":

		listItem.Location = "Dallas"
		listItem.Calculations = dataStructure.South.Dallas.Calculations[len(dataStructure.South.Dallas.Calculations)-1]
		responseToWeb.CityDataList = append(responseToWeb.CityDataList, *listItem)

		listItem.Location = "Jacksonville"
		listItem.Calculations = dataStructure.South.Jacksonville.Calculations[len(dataStructure.South.Jacksonville.Calculations)-1]
		responseToWeb.CityDataList = append(responseToWeb.CityDataList, *listItem)

		listItem.Location = "Charlotte"
		listItem.Calculations = dataStructure.South.Charlotte.Calculations[len(dataStructure.South.Charlotte.Calculations)-1]
		responseToWeb.CityDataList = append(responseToWeb.CityDataList, *listItem)

		listItem.Location = "Nashville"
		listItem.Calculations = dataStructure.South.Nashville.Calculations[len(dataStructure.South.Nashville.Calculations)-1]
		responseToWeb.CityDataList = append(responseToWeb.CityDataList, *listItem)

		listItem.Location = "Memphis"
		listItem.Calculations = dataStructure.South.Memphis.Calculations[len(dataStructure.South.Memphis.Calculations)-1]
		responseToWeb.CityDataList = append(responseToWeb.CityDataList, *listItem)

	case "WestCoast":
		listItem.Location = "Los Angeles"
		listItem.Calculations = dataStructure.WestCoast.LosAngeles.Calculations[len(dataStructure.WestCoast.LosAngeles.Calculations)-1]
		responseToWeb.CityDataList = append(responseToWeb.CityDataList, *listItem)

		listItem.Location = "San Francisco"
		listItem.Calculations = dataStructure.WestCoast.SanFrancisco.Calculations[len(dataStructure.WestCoast.SanFrancisco.Calculations)-1]
		responseToWeb.CityDataList = append(responseToWeb.CityDataList, *listItem)

		listItem.Location = "Seattle"
		listItem.Calculations = dataStructure.WestCoast.Seattle.Calculations[len(dataStructure.WestCoast.Seattle.Calculations)-1]
		responseToWeb.CityDataList = append(responseToWeb.CityDataList, *listItem)

		listItem.Location = "Portland"
		listItem.Calculations = dataStructure.WestCoast.Portland.Calculations[len(dataStructure.WestCoast.Portland.Calculations)-1]
		responseToWeb.CityDataList = append(responseToWeb.CityDataList, *listItem)

		listItem.Location = "San Diego"
		listItem.Calculations = dataStructure.WestCoast.SanDiego.Calculations[len(dataStructure.WestCoast.SanDiego.Calculations)-1]
		responseToWeb.CityDataList = append(responseToWeb.CityDataList, *listItem)

		listItem.Location = "Anchorage"
		listItem.Calculations = dataStructure.WestCoast.Anchorage.Calculations[len(dataStructure.WestCoast.Anchorage.Calculations)-1]
		responseToWeb.CityDataList = append(responseToWeb.CityDataList, *listItem)
	}

	toSend, _ := json.Marshal(responseToWeb)
	w.Write(toSend)
}

func fill() {
	jsonFile, _ := os.Open("src/dataStorage.json")
	jsonData, _ := ioutil.ReadAll(jsonFile)
	dataStructure := &Data{}
	json.Unmarshal(jsonData, &dataStructure)
	fmt.Println("Fill Start")
	/* North East */
	dataStructure.NorthEast.Boston.Calculations = append(dataStructure.NorthEast.Boston.Calculations, makeRequest(42.414140424834216, 42.21256130186242, -70.8934730834961, -71.20177691650392, 10))
	dataStructure.NorthEast.Brooklyn.Calculations = append(dataStructure.NorthEast.Brooklyn.Calculations, makeRequest(40.758364183254, 40.551558686334644, -73.73687118896483, -74.13855881103514, 10))
	dataStructure.NorthEast.Philadelphia.Calculations = append(dataStructure.NorthEast.Philadelphia.Calculations, makeRequest(40.202093630241286, 39.78440980977288, -74.81724916699217, -75.4338568330078, 10))
	dataStructure.NorthEast.WashingtonDC.Calculations = append(dataStructure.NorthEast.WashingtonDC.Calculations, makeRequest(38.995548, 38.791645, -76.909393, -77.119759, 10))
	dataStructure.NorthEast.Baltimore.Calculations = append(dataStructure.NorthEast.Baltimore.Calculations, makeRequest(39.396749, 39.197207, -76.529453, -76.711519, 10))

	/* South */
	dataStructure.South.Dallas.Calculations = append(dataStructure.South.Dallas.Calculations, makeRequest(33.04672230257906, 32.588541447873666, -96.3757438779297, -97.17911912207032, 10))
	dataStructure.South.Jacksonville.Calculations = append(dataStructure.South.Jacksonville.Calculations, makeRequest(30.586232, 30.098988, -81.328404, -82.049502, 10))
	dataStructure.South.Charlotte.Calculations = append(dataStructure.South.Charlotte.Calculations, makeRequest(35.431860761322845, 34.98640854430645, -80.42966737792969, -81.23304262207031, 10))
	dataStructure.South.Nashville.Calculations = append(dataStructure.South.Nashville.Calculations, makeRequest(36.405496, 35.989226, -86.515588, -87.054903, 10))
	dataStructure.South.Memphis.Calculations = append(dataStructure.South.Memphis.Calculations, makeRequest(35.264187, 34.994185, -89.637081, -90.304493, 10))

	/* MidWest */

	/* West Coast */
	dataStructure.WestCoast.LosAngeles.Calculations = append(dataStructure.WestCoast.LosAngeles.Calculations, makeRequest(34.4717411923813, 33.567993762094694, -117.60835725585937, -119.21510774414062, 10))
	dataStructure.WestCoast.SanFrancisco.Calculations = append(dataStructure.WestCoast.SanFrancisco.Calculations, makeRequest(37.842914, 37.707608, -122.32992, -122.536739, 10))
	dataStructure.WestCoast.Seattle.Calculations = append(dataStructure.WestCoast.Seattle.Calculations, makeRequest(47.734145, 47.491912, -122.224433, -122.465159, 10))
	dataStructure.WestCoast.Portland.Calculations = append(dataStructure.WestCoast.Portland.Calculations, makeRequest(45.714497, 45.395871, -122.471849, -122.919539, 10))
	dataStructure.WestCoast.SanDiego.Calculations = append(dataStructure.WestCoast.SanDiego.Calculations, makeRequest(33.114249, 32.534175, -116.90816, -117.309797, 10))
	dataStructure.WestCoast.Anchorage.Calculations = append(dataStructure.WestCoast.Anchorage.Calculations, makeRequest(61.326922, 60.733791, -148.473475, -150.420615, 10))

	updated, _ := json.Marshal((dataStructure))
	_ = ioutil.WriteFile("src/dataStorage.json", updated, 0777)
	fmt.Println("Fill Complete")
}

func main() {

	fileServer := http.FileServer(http.Dir("src"))
	http.Handle("/", http.StripPrefix("/", fileServer))
	http.HandleFunc("/fetchData", fetchData)
	http.HandleFunc("/hello", fetchData)
	http.HandleFunc("/newpage", handleNewPage)

	port := ":5505"
	fmt.Printf("Go backend listening on %s...\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	htmlBytes, err := ioutil.ReadFile("index.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	_, err = w.Write(htmlBytes)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func handleHello(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/city.html", http.StatusFound)
}

func handleNewPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("location", "Boston")
	htmlBytes, err := ioutil.ReadFile("city.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	_, err = w.Write(htmlBytes)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
