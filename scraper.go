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

type Response struct {
	User struct {
		IsLoggedIn                    bool   `json:"isLoggedin"`
		Email                         string `json:"email"`
		DisplayName                   string `json:"displayName"`
		HasHousingConnectorPermission bool   `json:"hasHousingConnectorPermission"`
		SavedHomesCount               int    `json:"savedHomesCount"`
		PersonalizedSearchTraceID     string `json:"personalizedSearchTraceID"`
		Guid                          string `json:"guid"`
		Zuid                          string `json:"zuid"`
		IsBot                         bool   `json:"isBot"`
		UserSpecializedSEORegion      bool   `json:"userSpecializedSEORegion"`
	} `json:"user"`
	RequestId int `json:"requestId"`
	Cat1      struct {
		SearchResults struct {
			ListResults []struct {
				Zpid              string     `json:"zpid"`
				Id                string     `json:"id"`
				ProviderListingId string     `json:"providerListingId"`
				ImgSrc            string     `json:"imgSrc"`
				HasImage          bool       `json:"hasImage"`
				CarouselPhotos    []struct{} `json:"carouselPhotos"`
				HdpData           struct {
					HomeInfo struct {
						Zpid string `json:"zpid"`
					} `json:"homeInfo"`
				} `json:"hdpData"`
			} `json:"listResults"`
		} `json:"searchResults"`
	} `json:"cat1"`
}

func main() {

	url := "https://www.zillow.com/search/GetSearchPageState.htm?searchQueryState=%7B%22mapBounds%22%3A%7B%22north%22%3A43.1847127353123%2C%22south%22%3A41.512239125025815%2C%22east%22%3A-70.65074389648439%2C%22west%22%3A-71.68071215820314%7D%2C%22mapZoom%22%3A9%2C%22isMapVisible%22%3Afalse%2C%22filterState%22%3A%7B%22isAllHomes%22%3A%7B%22value%22%3Atrue%7D%2C%22sortSelection%22%3A%7B%22value%22%3A%22globalrelevanceex%22%7D%7D%2C%22isListVisible%22%3Atrue%7D&wants={%22cat1%22:[%22listResults%22],%22cat2%22:[%22total%22],%22regionResults%22:[%22regionResults%22]}&requestId=2"
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

	jsonResponse, err := json.Marshal(string(body))
	_ = jsonResponse
	responseStructure := &Response{}
	json.Unmarshal(body, &responseStructure)
	fmt.Println(responseStructure)
	if err != nil {
		panic(err)
	}

	/*htmlString := html.NewTokenizer(reader)
	htmlString.Next()
	json.
	fmt.Println(htmlString.Token()) */

}
