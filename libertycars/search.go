package libertycars

import (
	"../common"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type SearchRequest struct {
	Location string
	Model    string
	MinYear  int
	MaxYear  int
	MinPrice int
	MaxPrice int
	MinMiles int
	MaxMiles int
}

type ListingRequest struct {
	Link string
}

func HandleSearch(w http.ResponseWriter, r *http.Request) {
	if !AuthenticateLibertyCarsRequest(w, r) {
		return
	}

	switch r.Method {
	case http.MethodPost:
		handleSearchPost(w, r)
	default:
		common.ErrorMethodNotAllowed(w, r)
	}
}

func HandleListing(w http.ResponseWriter, r *http.Request) {
	if !AuthenticateLibertyCarsRequest(w, r) {
		return
	}

	switch r.Method {
	case http.MethodPost:
		handleListingPost(w, r)
	default:
		common.ErrorMethodNotAllowed(w, r)
	}
}

func handleListingPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		common.ErrorMethodNotAllowed(w, r)
		return
	}

	var listingRequest ListingRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&listingRequest)
	if err != nil {
		fmt.Printf("Listing request decode error: %s\n", err.Error())
		common.ErrorBadRequest(w, r, err)
		return
	}

	url := listingRequest.Link

	req, _ := http.NewRequest("GET", url, nil)
	client := http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Printf("Listing request error: %s\n", err.Error())
		common.ErrorBadRequest(w, r, err)
		return
	}

	respBody, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		errMsg := fmt.Sprintf("something went wrong (%d)", resp.StatusCode)
		common.ErrorBadRequest(w, r, fmt.Errorf(errMsg))
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(respBody)
}

func handleSearchPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		common.ErrorMethodNotAllowed(w, r)
		return
	}

	var searchRequest SearchRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&searchRequest)
	if err != nil {
		fmt.Printf("Search request decode error: %s\n", err.Error())
		common.ErrorBadRequest(w, r, err)
		return
	}

	url := fmt.Sprintf("https://%s.craigslist.org/search/cto?auto_transmission=2&hasPic=1&auto_title_status=1"+
		"&max_auto_miles=%d&max_price=%d&auto_make_model=%s&min_auto_miles=%d&max_auto_year=%d&min_auto_year=%d&min_price=%d",
		searchRequest.Location, searchRequest.MaxMiles, searchRequest.MaxPrice, searchRequest.Model, searchRequest.MinMiles,
		searchRequest.MaxYear, searchRequest.MinYear, searchRequest.MinPrice)

	req, _ := http.NewRequest("GET", url, nil)
	client := http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Printf("Search request error: %s\n", err.Error())
		common.ErrorBadRequest(w, r, err)
		return
	}

	respBody, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		errMsg := fmt.Sprintf("something went wrong (%d)", resp.StatusCode)
		common.ErrorBadRequest(w, r, fmt.Errorf(errMsg))
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(respBody)
}
