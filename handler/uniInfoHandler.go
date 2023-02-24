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
The function takes in two parameters, a string variable "searchName" and an HTTP response writer "w",
and returns a slice of objects of the type "UniFromHipo".

Inside the function, it first constructs a URL string using the "searchName" parameter and
calls the "http.Get()" function to send an HTTP GET request to the API endpoint at the constructed URL.

If the request fails, it returns an HTTP error message using the "http.Error()" function,
with an HTTP status code of "http.StatusServiceUnavailable". If the request is successful,
it reads the response body using "ioutil.ReadAll()", which returns a byte array representing the response body.

The function then creates an empty slice of "UniFromHipo" objects and
unmarshal the byte array response body into this slice using the "json.Unmarshal()" function.
If there are any errors during the unmarshalling process, the function will panic.

Finally, the function returns the populated slice of "UniFromHipo" objects.
This function essentially makes an HTTP request to the given API endpoint and
returns the unmarshalled JSON data as a slice of objects of type "UniFromHipo".
*/
func getUniversityInfo(searchName string, w http.ResponseWriter) []UniFromHipo {

	// builds the url to the get information about the universities in the requested country from the api
	requestedUni := "http://universities.hipolabs.com/search?name=" + searchName
	responseUni, err := http.Get(requestedUni)
	if err != nil {
		http.Error(w, "Error in creating Get-request. Cannot reach service.", http.StatusServiceUnavailable)
	}

	bodyOfUniversity, err := ioutil.ReadAll(responseUni.Body)
	if err != nil {
		http.Error(w, "Unexpected format", http.StatusServiceUnavailable)
	}

	var universities []UniFromHipo
	err2 := json.Unmarshal(bodyOfUniversity, &universities)
	if err2 != nil {
		log.Fatal("Error when unmarshalling:", err)
	}

	return universities

}

/*
The getCountries function takes in a slice of UniFromHipo structs as input and
returns a slice of strings that represents the countries in which those universities are located.
It does this by iterating over the input slice of universities and checking if the country of each university is already in the list of countries.
If the country is not in the list, it is added to the list of countries.
The function uses a boolean variable isFound to keep track of whether a country has been found in the list of countries or not.
The isFound variable is reset to false at the end of each iteration to ensure that the function works correctly.
*/
func getCountries(unis []UniFromHipo) []string {
	var countries []string
	var isFound bool

	for _, u := range unis {
		country := u.Country
		for _, c := range countries {
			if country == c {
				isFound = true
			}
		}
		if !isFound {
			countries = append(countries, country)
		}
		isFound = false
	}
	return countries
}

/*
The dataOfUnis function takes a slice of UniFromHipo struct and converts it into a slice of University struct,
which is a custom struct defined somewhere in the code that's not shown here. It creates an empty slice dataOfUnis of type University.
(Slice is similar to an array, but unlike arrays, slices are dynamically sized, meaning their size can grow or shrink as needed)

For each element in the unis slice, it creates a new tempUni variable of type University,
sets its Name, Iso-code, Webpages and Country fields based on the corresponding fields of the uni element.
Then it appends the tempUni to the dataOfUnis slice.
*/
func dataOfUnis(unis []UniFromHipo) []University {
	var dataOfUnis []University
	var currentUni University
	for _, uni := range unis {
		currentUni.Name = uni.Name
		currentUni.Isocode = uni.AlphaTwoCode
		currentUni.Webpages = uni.WebPages
		currentUni.Country = uni.Country
		/*
			for _, s := range restCountries {
				if s.Country == r.Country {
					tempU.Languages = s.Languages
					tempU.Maps = s.Maps
				}
			}
		*/
		dataOfUnis = append(dataOfUnis, currentUni)
	}

	return dataOfUnis
}

/*
This function is an HTTP handler that takes a request and returns a JSON response with information about universities that match a search name.
The handler first checks that the URL path contains exactly 5 segments and extracts the search name from the last segment.
It then uses helper functions to retrieve and process university data, including finding the countries in which
the universities are located and formatting the data as a slice of University structs.

Finally, the function sets the content type header to application/json, encodes the university data as JSON,
and writes the response to the client. If an error occurs at any stage, the function returns an error status code with a corresponding error message.
*/
func UniInfoHandler(w http.ResponseWriter, r *http.Request) {

	// splits the path into components.
	pathComponents := strings.Split(r.URL.Path, "/")
	// Expects the path to have exactly 5 segments
	if len(pathComponents) != 5 {
		http.Error(w, "Malformed URL. Expecting format ../name", http.StatusBadRequest)
		log.Println("Malformed URL in request.")
		return
	}

	// searchName is the index of the last element in the slice.
	searchName := pathComponents[len(pathComponents)-1]

	unisInfo := getUniversityInfo(searchName, w)
	if unisInfo == nil {
		http.Error(w, "No data found", http.StatusInternalServerError)
		return
	}

	countries := getCountries(unisInfo)
	if countries == nil {
		http.Error(w, "Could not put together data", http.StatusInternalServerError)
		return
	}

	dataOfUnis := dataOfUnis(unisInfo)
	if dataOfUnis == nil {
		http.Error(w, "Could not put together data", http.StatusInternalServerError)
		return
	}

	w.Header().Add("content-type", "application/json")
	encoder := json.NewEncoder(w)

	err := encoder.Encode(dataOfUnis)
	if err != nil {
		http.Error(w, "Error during encoding", http.StatusInternalServerError)
		return
	}

	http.Error(w, "", http.StatusOK)
	fmt.Println(r.URL.Path)

}
