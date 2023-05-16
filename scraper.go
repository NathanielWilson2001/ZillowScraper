package main

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"image/color"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
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
func calculate(responses []Response, numberOfPages int, output *widget.Label) {

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

	outputString := fmt.Sprintf("Count: %d \nAverage Price = %.2f \nAverage Square Foot = %.2f\nAverage Price per Square Foot = %.2f\nAverage Zestimate = %.2f\nAverage ZRentEstimate = %.2f\nHousing Type Breakdown:\n\tMultifamily: %d\n\tSingle Family: %d\n\tCondo: %d\n",
		runningTotalEntries,
		averagePriceSum,
		averageSquareFootSum,
		averagePriceSum/averageSquareFootSum,
		averageZestimate,
		averageZestimate,
		housingType[0],
		housingType[1],
		housingType[2])

	output.SetText(outputString)
}

/*
*	The purpose of makeRequest is to take in 4 GPS coordinate values that are used to set a boundary for a given search location,
*	as well as the Page number, and generate a request URL to Zillow. The request is then sent and the response is stored
*	in the response object. Response calculate is called to output the calculations from the given data set.
 */
func makeRequest(north float64, south float64, east float64, west float64, numPages int, output *widget.Label) {

	var responses []Response
	for pageNumber := 1; pageNumber <= numPages; pageNumber++ {
		url := "https://www.zillow.com/search/GetSearchPageState.htm?"
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

		//responseToString(*responseStructure)
		if err != nil {
			panic(err)
		}

		responses = append(responses, *responseStructure)
	}

	calculate(responses, numPages, output)

}

/*
*	The purpose of displayApp() is to display the GUI interface for the application
*
*
 */
func displayApp() {

	app := app.New()
	w := app.NewWindow("Zillow Scraper")

	rectangle := canvas.NewRectangle(color.White)
	w.SetContent(rectangle)

	label1 := widget.NewLabel("Results will output here:")
	label1.TextStyle.Bold = true
	value1 := widget.NewLabel("")

	inputNorth := widget.NewEntry()
	inputNorth.SetPlaceHolder("Enter North Coordinate...")
	inputSouth := widget.NewEntry()
	inputSouth.SetPlaceHolder("Enter South Coordinate...")
	inputEast := widget.NewEntry()
	inputEast.SetPlaceHolder("Enter East Coordinate...")
	inputWest := widget.NewEntry()
	inputWest.SetPlaceHolder("Enter West Coordinate...")
	inputPages := widget.NewEntry()
	inputPages.SetPlaceHolder("Enter number of pages")
	grid := container.NewVBox(inputNorth, inputSouth, inputEast, inputWest, inputPages, widget.NewButton("Submit data", func() {
		north, _ := strconv.ParseFloat(inputNorth.Text, 64)
		south, _ := strconv.ParseFloat(inputSouth.Text, 64)
		east, _ := strconv.ParseFloat(inputEast.Text, 64)
		west, _ := strconv.ParseFloat(inputWest.Text, 64)
		pageNum, _ := strconv.ParseInt(inputPages.Text, 10, 32)
		makeRequest(north, south, east, west, int(pageNum), value1)

	}), widget.NewCard("Reslts will output here", "", nil), value1)

	w.SetContent(grid)
	w.CenterOnScreen()
	w.Resize(fyne.NewSize(1080, 720))

	w.ShowAndRun()
}
func main() {

	displayApp()
	//makeRequest(42.28936961396935, 41.86942883946931, -70.90983246728517, -71.16732453271486, 10)

	//makeRequest(43.1847127353123, 41.512239125025815, -70.65074389648439, -71.68071215820314, 1)

}
