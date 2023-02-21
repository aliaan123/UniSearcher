package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

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

func getCountryInfo(w http.ResponseWriter, countryName string) []Country {

	// builds the url to search for the requested country
	requestedCountry := "https://restcountries.com/v3.1/name/" + countryName + "?fields=name,languages,maps,cca2"
	// makes an HTTP GET request to an external API endpoint
	responseCountry, err := http.Get(requestedCountry)
	if err != nil {
		http.Error(w, "Error in creating Get-request. Cannot reach service.", http.StatusServiceUnavailable)
	}

	// reads the response from the API endpoint
	bodyOfCountry, err := ioutil.ReadAll(responseCountry.Body)
	if err != nil {
		http.Error(w, "Unexpected format", http.StatusServiceUnavailable)
	}

	// unmarshal the JSON data into a slice of struct of type Country
	var response []Country
	err2 := json.Unmarshal(bodyOfCountry, &response)
	if err2 != nil {
		log.Fatal("Error when unmarshalling:", err)
	}

	// returns that struct as the function output
	return response

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

	// can be done like this instead
	//requestedCountry := "http://universities.hipolabs.com/search?name=" + countryName + "?fields=borders"

	// makes an HTTP GET request to an external API endpoint
	//resp2, err := http.Get(requestedCountry)
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

// function that finds all universities in the countries

// BorderCountries is a function that turns the alpha-two-code of countries into countries.
func BorderCountries(w http.ResponseWriter, borders []Borders) []Country {
	// empty slice of Country struct, for storing all the bordering countries of a country
	var borderingCountries []Country

	// borders, which is a slice of structs of type "Borders" is a slice
	// containing alpha-2-codes for bordering countries of a particular country

	// loops through all the bordering countries in the given borders slice of structs of type Borders,
	// which contains info about the bordering countries of a particular country
	for i := range borders {
		// for each specific bordering country to a country it will search for the bordering countries by using their isocode.
		// It searches for the countries by the isocode and appends them to the bordering countries variable.
		for j := range borders[i].Borders {
			// uses the findCountryByAlpha2Code function to find the bordering countries based on the alpha-two-code.
			// for each alpha-2-code in the borders slice, the findCountryByAlpha2Code is called and
			// passes the alpha-2-code in as argument in the findCountryByAlpha2Code function, which returns
			// a Country struct that is appended to the borderingCountries slice, creating
			borderingCountries = append(borderingCountries, findCountryByAlpha2Code(w, borders[i].Borders[j]))
		}
	}

	// returns borderingCountries, a slice of Country structs.
	return borderingCountries
}

// Function that finds countries based on the alpha-two-code of the country.
func findCountryByAlpha2Code(w http.ResponseWriter, alpha2Code string) Country {

	url1 := "https://restcountries.com/v3.1/alpha/"
	url2 := "?fields=name,languages,maps,cca2"
	var buildString strings.Builder

	// building the url with the country name
	buildString.WriteString(url1)
	buildString.WriteString(alpha2Code)
	buildString.WriteString(url2)

	//requestedCountry := "https://restcountries.com/v3.1/alpha/" + alpha2Code + "?fields=name,languages,maps,cca2"

	//resp2, err := http.Get(requestedCountry)
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
