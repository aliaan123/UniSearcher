package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

/*
func getCountryInfo(searchName string, w http.ResponseWriter) []UniFromHipo {

}
*/

func NeighbourUnisHandler(w http.ResponseWriter, r *http.Request) {

	// splits the path into components.
	pathComponents := strings.Split(r.URL.Path, "/")
	expectedLength := len(strings.Split(NEIGHBOR_UNIS_PATH, "/")) + 1

	// Expects the path to have exactly 5 segments
	if len(pathComponents) != 6 || len(pathComponents) < expectedLength {
		status := http.StatusBadRequest
		http.Error(w, "Expecting format .../name/university", status)
		return
	}

	searchCountryName := pathComponents[expectedLength-2]
	searchUniName := pathComponents[expectedLength-1]

	if searchCountryName == "" || searchUniName == "" {
		http.Error(w, "Country name and partial or complete name of university must be given", http.StatusBadRequest)
	}

	// searchName is the index of the last element in the slice.
	//searchName := pathComponents[len(pathComponents)-1]

	http.Error(w, "", http.StatusOK)
	fmt.Println(r.URL.Path)

}

// functions that finds the borders of a country in iso code format
func getBordersOfCountry(w http.ResponseWriter, countryName string) []Borders {

	url1 := "https://restcountries.com/v3.1/name/"
	url2 := "?fields=borders"
	var buildString strings.Builder

	// building the url with the country name
	buildString.WriteString(url1)
	buildString.WriteString(countryName)
	buildString.WriteString(url2)

	// makes an HTTP GET request to an external API endpoint
	resp, err := http.Get(buildString.String())
	if err != nil {
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
		return nil
	}

	//  reads the response from the API endpoint
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Unexpected format", http.StatusServiceUnavailable)
	}

	// unmarshal the JSON data into a slice of structs of type Borders
	var response []Borders
	err2 := json.Unmarshal(respBody, &response)

	if err2 != nil {
		log.Fatal("Error when unmarshalling:", err)
	}

	// returns that slice as the function output
	return response

}

/*
func BorderCountries(w http.ResponseWriter, borders []Borders) []Country {


}
*/

// Function that finds countries based on the alpha-two-code of the country.
func findCountryByAlpha2Code(w http.ResponseWriter, alpha2Code string) Country {

	url1 := "https://restcountries.com/v3.1/alpha/"
	url2 := "?fields=name,languages,maps,cca2"
	var buildString strings.Builder

	// building the url with the country name
	buildString.WriteString(url1)
	buildString.WriteString(alpha2Code)
	buildString.WriteString(url2)

	// makes an HTTP GET request to an external API endpoint
	resp, err := http.Get(buildString.String())
	// handling errors
	if err != nil {
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
		return Country{}
	}

	// reads the response from the API endpoint
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Unexpected format", http.StatusServiceUnavailable)
		return Country{}
	}

	// unmarshal the JSON data into a struct of type Country
	var response Country
	err2 := json.Unmarshal(respBody, &response)

	if err2 != nil {
		log.Fatal("Error when unmarshalling:", err)
		return Country{}
	}

	// returns that struct as the function output
	return response

}
