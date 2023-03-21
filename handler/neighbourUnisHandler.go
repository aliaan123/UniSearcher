package handler

import (
	"Assignment1/constants"
	"Assignment1/structs"
	"Assignment1/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// NeighbourUnisHandler is a function, an HTTP handler that makes a request to a given API endpoint, and returns a JSON response.
func NeighbourUnisHandler(w http.ResponseWriter, r *http.Request) {

	// splits the path into components.
	pathComponents := strings.Split(r.URL.Path, "/")
	// find the expected length of the URL
	expectedLength := len(strings.Split(constants.NEIGHBOR_UNIS_PATH, "/")) + 1

	// expects the path to have exactly 5 segments, if it doesn't we send an error
	if len(pathComponents) != 6 || len(pathComponents) < expectedLength || len(pathComponents) > expectedLength {
		status := http.StatusBadRequest
		http.Error(w, "Expecting format .../name/university", status)
		return
	}

	// The university name that is searched for is the last segment of the URL. We replace all spaces with %20.
	searchUniName := strings.ReplaceAll(pathComponents[expectedLength-1], " ", "%20")
	// The country name that is searched for is the second-to-last segment of the URL. We also replace all spaces with %20.
	searchCountryName := strings.ReplaceAll(pathComponents[expectedLength-2], " ", "%20")

	// Checks if the search values are empty. If they are we send an error.
	if searchCountryName == "" || searchUniName == "" {
		http.Error(w, "partial or complete name of university and country name must be given.", http.StatusBadRequest)
	}

	var countrySearchedFor []structs.Country
	// adds country searched for in the slice
	countrySearchedFor = append(countrySearchedFor, utils.GetCountryInfo(w, searchCountryName)...)
	// uses the getBorderOfCountry-function to get the bordering countries of the searched for country, and then uses
	// the BorderCountries-function to translate the alpha-two-code of the bordering countries returned from getBorders function
	// into countries. Then it adds these bordering countries to the same slice.
	countrySearchedFor = append(countrySearchedFor, BorderCountries(w, getBordersOfCountry(w, searchCountryName))...)

	// gets all the universities that matches with the searchUniName
	var uni = utils.GetUniversityInfo(searchUniName, w)

	// filters out universities that are in the given searched for country.
	var uniFiltered = filterOutCountrySearchedFor(uni, searchCountryName)

	// combines the data of the country that was searched for and the data of the uni that was searched for.
	var combinedData = utils.CombineData(countrySearchedFor, utils.DataOfUnis(uniFiltered))

	var limit int

	if r.URL.RawQuery != "" {
		partsOfQuery := strings.Split(r.URL.RawQuery, "=")
		if len(partsOfQuery) == 2 && partsOfQuery[0] == "limit" {
			limit, _ = strconv.Atoi(partsOfQuery[1])
			if limit <= 0 {
				http.Error(w, "Limit should be a positive integer.", http.StatusBadRequest)
				return
			} else if limit > 100 {
				http.Error(w, "Limit should be less than or equal to 100.", http.StatusBadRequest)
				return
			}
		}
	} else {
		limit = 0
	}

	if limit == 0 {
		limit = len(combinedData)
	}

	if len(combinedData) >= limit {
		combinedData = append(combinedData[:limit], combinedData[len(combinedData):]...)
	}

	w.Header().Add("content-type", "application/json")
	// encode the combined data in json
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "\t")
	err := encoder.Encode(combinedData)
	if err != nil {
		http.Error(w, "ERROR encoding.", http.StatusInternalServerError)
		return
	}

	http.Error(w, "", http.StatusOK)
	fmt.Println(r.URL.Path)

}

// function for filtering out the country that was searched for, so that universities from the given country is not included
func filterOutCountrySearchedFor(universities []structs.UniFromHipo, countrySearchedFor string) []structs.UniFromHipo {
	var filteredUnis []structs.UniFromHipo

	// removes the spaces from the search string
	countrySearchedFor = strings.ReplaceAll(countrySearchedFor, "%20", " ")
	for i := range universities {
		if strings.ToLower(universities[i].Country) != strings.ToLower(countrySearchedFor) {
			filteredUnis = append(filteredUnis, universities[i])
		}
	}
	return filteredUnis
}

// functions that finds the borders of a country in iso code format
func getBordersOfCountry(w http.ResponseWriter, countryName string) []structs.Borders {

	// removes the spaces from the search string
	countryName = strings.ReplaceAll(countryName, " ", "%20")
	// builds the url to get the bordering countries of the requested country from the api
	requestedCountry := constants.COUNTRY_NAME + countryName + constants.COUNTRY_FIELD_BORDER
	// makes an HTTP GET request to an external API endpoint
	resp, err := http.Get(requestedCountry)
	// handles errors if request fails.
	if err != nil {
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
		return nil
	}

	// reads the response body from the API endpoint, a byte array representing the response body is returned.
	respBody, err := ioutil.ReadAll(resp.Body)
	// handles errors if reading the response fails
	if err != nil {
		http.Error(w, "Unexpected format", http.StatusServiceUnavailable)
	}

	// creates an empty slice of "Borders" objects
	var response []structs.Borders
	// unmarshal the JSON data (the byte array response body) into a slice of type Borders
	err2 := json.Unmarshal(respBody, &response)
	// handle errors if the unmarshal fails
	if err2 != nil {
		http.Error(w, "Error when unmarshalling. Unexpected format", http.StatusServiceUnavailable)
	}

	// returns that slice as the function output
	return response

}

// BorderCountries is a function that turns the alpha-two-code of countries into countries.
func BorderCountries(w http.ResponseWriter, borders []structs.Borders) []structs.Country {
	// empty slice of Country struct, for storing all the bordering countries of a country
	var borderingCountries []structs.Country

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
func findCountryByAlpha2Code(w http.ResponseWriter, alpha2Code string) structs.Country {

	// builds the url to get the country that matches the alpha-two-code from the api
	//requestedCountry := "https://restcountries.com/v3.1/alpha/" + alpha2Code + "?fields=name,languages,maps,cca2"
	requestedCountry := constants.COUNTRY_ALPHACODE + alpha2Code + constants.COUNTRY_FIELDS
	// makes an HTTP GET request to an external API endpoint
	resp, err := http.Get(requestedCountry)
	// handles errors if request fails.
	if err != nil {
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
		return structs.Country{}
	}

	// reads the response body from the API endpoint, a byte array representing the response body is returned.
	respBody, err := ioutil.ReadAll(resp.Body)
	// handle errors if reading the response body fails
	if err != nil {
		http.Error(w, "Error when reading response body. Unexpected format", http.StatusServiceUnavailable)
		return structs.Country{}
	}

	// creates an empty slice of "Country" objects
	var response structs.Country
	// unmarshal the JSON data (the byte array response body) into a slice of type Country
	err2 := json.Unmarshal(respBody, &response)
	// handles errors if the unmarshal fails.
	if err2 != nil {
		http.Error(w, "Error when unmarshalling. Unexpected format", http.StatusServiceUnavailable)
		return structs.Country{}
	}

	// returns that slice as the function output
	return response

}
