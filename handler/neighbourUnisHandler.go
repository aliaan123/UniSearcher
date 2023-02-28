package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// TODO : 1. fix limit i url
// 		  2. error handling hvis man feks skriver by i name istedet for country, eller hvis man typer feil
//		  3. skriv README
//        4. country API should not include the country that was searched for.  ## FIXED
//        5. Gjør diagHandler endpoint
//        6. fix language og maps i UNIVERSITY	##FIXED

func NeighbourUnisHandler(w http.ResponseWriter, r *http.Request) {

	// splits the path into components.
	pathComponents := strings.Split(r.URL.Path, "/")
	expectedLength := len(strings.Split(NEIGHBOR_UNIS_PATH, "/")) + 1

	// Expects the path to have exactly 5 segments
	if len(pathComponents) != 6 || len(pathComponents) < expectedLength || len(pathComponents) > expectedLength {
		status := http.StatusBadRequest
		http.Error(w, "Expecting format .../name/university", status)
		return
	}

	// The country name that is searched for is at the second last index of the URL.
	searchCountryName := strings.ReplaceAll(pathComponents[expectedLength-2], " ", "%20")
	// The university name that is searched for is at the last index of the URL.
	searchUniName := strings.ReplaceAll(pathComponents[expectedLength-1], " ", "%20")

	if searchCountryName == "" || searchUniName == "" {
		http.Error(w, "Country name and partial or complete name of university must be given", http.StatusBadRequest)
	}

	var countrySearchedFor []Country
	// adds country searched for in the slice
	countrySearchedFor = append(countrySearchedFor, getCountryInfo(w, searchCountryName)...)
	//removeCountrySearchedFor()
	// uses the getBorderOfCountry-function to get the bordering countries of the searched for country, and then uses
	// the BorderCountries-function to translate the alpha-two-code of the bordering countries returned from getBorders function
	// into countries. Then it adds these bordering countries to the same slice.
	countrySearchedFor = append(countrySearchedFor, BorderCountries(w, getBordersOfCountry(w, searchCountryName))...)

	// gets all the universities that matches with the searchUniName
	var uni = getUniversityInfo(searchUniName, w)
	// filters out universities that are in the given country.
	var uniFiltered = removeCountrySearchedFor(uni, searchCountryName)
	// combines the data of the country that was searched for and the data of the uni that was searched for.
	var combinedData = combineData(countrySearchedFor, dataOfUnis(uniFiltered))

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

// function that combines the data from the university and the country
func combineData(countryInfo []Country, uniInfo []University) []CombinedStruct {
	var combinedData []CombinedStruct
	for i := range uniInfo {
		for j := range countryInfo {
			// matches the alpha-two-code of the university to the cca2 code of the country to match the country to the university.
			if uniInfo[i].Isocode == countryInfo[j].Cca2 {
				combinedData = append(combinedData, CombinedStruct{
					Name:      uniInfo[i].Name,
					Country:   uniInfo[i].Country,
					TwoCode:   uniInfo[i].Isocode,
					WebPages:  uniInfo[i].Webpages,
					Languages: countryInfo[j].Languages,
					Map:       countryInfo[j].Maps.OpenStreetMaps,
				})
			}
		}
	}

	return combinedData
}

// function for filtering out the country that was searched for, so that universities from the given country is not included
func removeCountrySearchedFor(universities []UniFromHipo, countrySearchedFor string) []UniFromHipo {
	var filteredUnis []UniFromHipo
	for i := range universities {
		if universities[i].Country != countrySearchedFor {
			filteredUnis = append(filteredUnis, universities[i])
		}
	}
	return filteredUnis
}

// PUTTA INN []UNIVERSITY ISTEDET FOR UNIFROMHIPO FOR Å TESTE: FJERN OM IKKEFUNKER.
// checks if a country is in the slice. If the country is not in the slice, it gets added. If it already is the check is set to true.
func checkCountries(w http.ResponseWriter, countries []Country, universities []University) []Country {
	var check = false
	for i := range universities {
		check = false
		for j := range countries {
			// checks if the alpha-two-code of the university matches the cca2 of the country
			if universities[i].Isocode == countries[j].Cca2 {
				j = len(countries)
				check = true
			}
		}
		if !check {
			countries = append(countries, getCountryInfo(w, universities[i].Country)...)
		}
	}
	return countries
}

// function for getting the information about a country from the api
func getCountryInfo(w http.ResponseWriter, countryName string) []Country {

	// builds the url to search for the requested country in the api
	requestedCountry := "https://restcountries.com/v3.1/name/" + countryName + "?fields=name,languages,maps,cca2"
	// makes an HTTP GET request to an external API endpoint
	responseCountry, err := http.Get(requestedCountry)
	if err != nil {
		http.Error(w, "Error in creating Get-request. Cannot reach service.", http.StatusServiceUnavailable)
		return nil
	}

	//defer responseCountry.Body.Close()

	// reads the response from the API endpoint
	bodyOfCountry, err := ioutil.ReadAll(responseCountry.Body)
	if err != nil {
		http.Error(w, "Unexpected format", http.StatusServiceUnavailable)
		return nil
	}

	// unmarshal the JSON data into a slice of struct of type Country
	var response []Country
	err = json.Unmarshal(bodyOfCountry, &response)
	if err != nil {
		log.Fatal("Error when unmarshalling:", err)
	}

	// returns that struct as the function output
	return response

}

// functions that finds the borders of a country in iso code format
func getBordersOfCountry(w http.ResponseWriter, countryName string) []Borders {

	// builds the url to get the bordering countries of the requested country from the api
	requestedCountry := "https://restcountries.com/v3.1/name/" + countryName + "?fields=borders"

	// makes an HTTP GET request to an external API endpoint
	resp, err := http.Get(requestedCountry)
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

	// builds the url to get the country that matches the alpha-two-code from the api
	requestedCountry := "https://restcountries.com/v3.1/alpha/" + alpha2Code + "?fields=name,languages,maps,cca2"

	// makes an HTTP GET request to an external API endpoint
	resp, err := http.Get(requestedCountry)
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
