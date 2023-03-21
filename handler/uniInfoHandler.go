package handler

import (
	"Assignment1/structs"
	"Assignment1/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// UniInfoHandler is a function, an HTTP handler that makes a request to a given API endpoint, and returns a JSON response of info about unis.
func UniInfoHandler(w http.ResponseWriter, r *http.Request) {

	// splits the path into components.
	pathComponents := strings.Split(r.URL.Path, "/")
	// Expects the path to have exactly 5 segments
	if len(pathComponents) < 5 || len(pathComponents) > 5 {
		http.Error(w, "Malformed URL. Expecting format ../name", http.StatusBadRequest)
		log.Println("Malformed URL in request.")
		return
	}

	// gets the searchName from the URL, searchName is the last segment of the URL. Removes the spaces from the search string.
	searchName := strings.ReplaceAll(pathComponents[len(pathComponents)-1], " ", "%20")

	// retrieves information about university based on search name
	unisInfo := utils.GetUniversityInfo(searchName, w)
	if unisInfo == nil || len(unisInfo) == 0 {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	var countries []structs.Country
	// appends countries and information about countries that were found to the empty slice above.
	countries = append(countries, utils.GetCountryInfo(w, unisInfo[0].Country)...)

	// converts retrieved information about universities into slice of University objects.
	// slice of UniFromHipo struct converted to slice of University objects
	dataOfUniversity := utils.DataOfUnis(unisInfo)
	if dataOfUniversity == nil {
		http.Error(w, "Internal server error. Could not put together data", http.StatusInternalServerError)
		return
	}

	// makes sure that all countries associated with unis in the universities slice are included in the countries slice
	countries = utils.CheckCountries(w, countries, dataOfUniversity)
	if countries == nil {
		http.Error(w, "Internal server error. Could not put together data", http.StatusInternalServerError)
		return
	}

	// combines the data of the country that was searched for and the data of the uni that was searched for.
	var combinedData = utils.CombineData(countries, utils.DataOfUnis(unisInfo))
	if combinedData == nil {
		http.Error(w, "Internal server error. Could not put together data", http.StatusInternalServerError)
		return
	}

	// sets content type header to json
	w.Header().Add("content-type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "\t")
	// encodes data to Json and sends the encoded data in the response body
	err := encoder.Encode(combinedData)
	if err != nil {
		http.Error(w, "Error during encoding", http.StatusInternalServerError)
		return
	}

	http.Error(w, "", http.StatusOK)
	fmt.Println(r.URL.Path)

}
