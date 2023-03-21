package utils

import (
	"Assignment1/constants"
	"Assignment1/structs"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// Utils file that holds functions both used in uniInfoHandler and neighbourUnisHandler

// GetUniversityInfo is a function which makes HTTP request to the given API endpoint and, returns the unmarshalled JSON data as a slice of type "UniFromHipo".
func GetUniversityInfo(searchName string, w http.ResponseWriter) []structs.UniFromHipo {

	// removes the spaces from the search string
	searchName = strings.ReplaceAll(searchName, " ", "%20")
	// builds the url to the get information about the universities in the requested country from the api
	//requestedUni := "http://universities.hipolabs.com/search?name_contains=" + searchName
	requestedUni := constants.UNI_NAME + searchName
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

// CombineData function that combines the data-struct from the university and the country
func CombineData(countryInfo []structs.Country, uniInfo []structs.University) []structs.CombinedStruct {
	var combinedStruct []structs.CombinedStruct
	for i := range uniInfo {
		for j := range countryInfo {
			// matches the alpha-two-code of the university to the cca2 code of the country to match the country to the university.
			if uniInfo[i].Isocode == countryInfo[j].Cca2 {
				combinedStruct = append(combinedStruct, structs.CombinedStruct{
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

	return combinedStruct
}

// CheckCountries function that makes sure that all countries associated with unis in the universities slice are included in the countries slice
func CheckCountries(w http.ResponseWriter, countries []structs.Country, universities []structs.University) []structs.Country {
	var check = false
	// loops through universities in the University slice
	for i := range universities {
		check = false
		// loops through countries in the Country slice
		for j := range countries {
			// checks if the alpha-two-code of the university matches the cca2 of the country in the slice.
			if universities[i].Isocode == countries[j].Cca2 {
				j = len(countries)
				// sets check to true if there is a match.
				check = true
			}
		}
		if !check {
			// if there is no match, we first get information about the country the universities resides in, and appends them to the country slice
			countries = append(countries, GetCountryInfo(w, universities[i].Country)...)
		}
	}
	// returns updated slice of countries with info about all the countries where the universities reside in , with no duplicates.
	return countries
}

// GetCountryInfo function which makes HTTP request to given API endpoint to retrieve information about a specific country based on a given string, and returns the unmarshalled JSON data.
func GetCountryInfo(w http.ResponseWriter, countryName string) []structs.Country {

	// removes the spaces from the search string
	countryName = strings.ReplaceAll(countryName, " ", "%20")
	// builds the url to search for the requested country in the api
	requestedCountry := constants.COUNTRY_NAME + countryName + constants.COUNTRY_FIELDS
	// makes an HTTP GET request to an external API endpoint
	responseCountry, err := http.Get(requestedCountry)
	// handles errors if request fails.
	if err != nil {
		http.Error(w, "Error in creating Get-request. Cannot reach service.", http.StatusServiceUnavailable)
		return nil
	}

	//defer responseCountry.Body.Close()

	// reads the response body from the API endpoint, a byte array representing the response body is returned.
	bodyOfCountry, err := ioutil.ReadAll(responseCountry.Body)
	// handle errors if reading the response body fails
	if err != nil {
		http.Error(w, "Unexpected format", http.StatusServiceUnavailable)
		return nil
	}

	// unmarshal the JSON data (the byte array response body) into a slice of type Country
	var response []structs.Country
	err = json.Unmarshal(bodyOfCountry, &response)
	// handle errors if the unmarshal fails
	if err != nil {
		http.Error(w, "Error when unmarshalling. Unexpected format", http.StatusServiceUnavailable)
		//log.Fatal("Error when unmarshalling3:", err)
		return nil
	}

	// returns that struct as the function output
	return response
}

// DataOfUnis function takes a slice of UniFromHipo struct, and converts it into a slice of University struct.
func DataOfUnis(unis []structs.UniFromHipo) []structs.University {
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
