package handler

import (
	"Assignment1/structs"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// getUniversityInfo is a function which makes HTTP request to the given API endpoint and, returns the unmarshalled JSON data as a slice of type "UniFromHipo".
func getUniversityInfo(searchName string, w http.ResponseWriter) []structs.UniFromHipo {

	// removes the spaces from the search string
	searchName = strings.ReplaceAll(searchName, " ", "%20")
	// builds the url to the get information about the universities in the requested country from the api
	requestedUni := "http://universities.hipolabs.com/search?name_contains=" + searchName
	// makes an HTTP GET request to an external API endpoint
	responseUni, err := http.Get(requestedUni)
	// handles errors if request fails.
	if err != nil {
		http.Error(w, "Error in creating Get-request. Cannot reach service.", http.StatusServiceUnavailable)
	}

	// reads the response body from the API endpoint, a byte array representing the response body is returned.
	bodyOfUniversity, err := ioutil.ReadAll(responseUni.Body)
	// handles errors if reading the response fails
	if err != nil {
		http.Error(w, "Unexpected format", http.StatusServiceUnavailable)
	}

	// creates an empty slice of "UniFromHipo" objects
	var universities []structs.UniFromHipo
	// unmarshal the JSON data (the byte array response body) into a slice of type UniFromHipo
	err2 := json.Unmarshal(bodyOfUniversity, &universities)
	// handle errors if the unmarshal fails
	if err2 != nil {
		log.Fatal("Error when unmarshalling:", err)
	}

	// returns a populated slice of UniFromHipo objects as the function output
	return universities

}

// The function dataOfUnis takes a slice of UniFromHipo struct, and converts it into a slice of University struct.
func dataOfUnis(unis []structs.UniFromHipo) []structs.University {
	// create empty slice of type University struct
	var dataOfUnis []structs.University
	// create a temporary variable of type University, objects which are appended to the slice
	var currentUni structs.University
	// loops through the slice of UniFromHipo which was taken in as parameter.
	// For each element in the unis slice, a new variable of type University is created,
	// and its name, alphaTwoCode, webpages and country fields are set to the corresponding fields of the uni element.
	for _, uni := range unis {
		currentUni.Name = uni.Name
		currentUni.Isocode = uni.AlphaTwoCode
		currentUni.Webpages = uni.WebPages
		currentUni.Country = uni.Country
		// appends the temporary object to the dataOfUnis slice.
		dataOfUnis = append(dataOfUnis, currentUni)
	}

	return dataOfUnis
}

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

	unisInfo := getUniversityInfo(searchName, w)
	if unisInfo == nil || len(unisInfo) == 0 {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	var countries []structs.Country
	countries = append(countries, getCountryInfo(w, unisInfo[0].Country)...)

	dataOfUniversity := dataOfUnis(unisInfo)
	if dataOfUniversity == nil {
		http.Error(w, "Could not put together data", http.StatusInternalServerError)
		return
	}

	countries = checkCountries(w, countries, dataOfUniversity)
	if countries == nil {
		http.Error(w, "Could not put together data", http.StatusInternalServerError)
		return
	}

	// combines the data of the country that was searched for and the data of the uni that was searched for.
	var combinedData = combineData(countries, dataOfUnis(unisInfo))
	if combinedData == nil {
		http.Error(w, "Could not put together data", http.StatusInternalServerError)
		return
	}

	w.Header().Add("content-type", "application/json")
	encoder := json.NewEncoder(w)

	err := encoder.Encode(combinedData)
	if err != nil {
		http.Error(w, "Error during encoding", http.StatusInternalServerError)
		return
	}

	http.Error(w, "", http.StatusOK)
	fmt.Println(r.URL.Path)

}
