package handler

import (
	"Assignment1/time"
	"encoding/json"
	"net/http"
)

// DiagHandler is a function handles HTTP requests, combines the responses into a struct and encodes the struct as a JSON and writes to the body
func DiagHandler(w http.ResponseWriter, r *http.Request) {

	// makes call to uniAPIstatus function to retrieve HTTP status code
	var uniStatus = uniAPIstatus(w)
	// makes call to countryAPIstatus function to retrieve HTTP status code
	var countryStatus = countryAPIstatus(w)

	// calls the diagResponse function, which creates a struct and combines the information
	combinedResponses := diagResponse(uniStatus, countryStatus, time.TimeSinceStart())

	// sets the content-type header of the response to JSON
	w.Header().Add("content-type", "application/json")
	// encodes the struct as JSON and writes to the body
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "\t")
	err := encoder.Encode(combinedResponses)
	if err != nil {
		http.Error(w, "ERROR when encoding.", http.StatusInternalServerError)
	}

}

// function that creates a new DiagStruct and populates it with the status codes, version and uptime of the service, and then returns the struct
func diagResponse(uniAPI string, countryAPI string, time float64) DiagStruct {

	diagResponses := DiagStruct{
		UniversitiesApi: uniAPI,
		CountriesApi:    countryAPI,
		Version:         VERSION,
		Uptime:          time,
	}
	return diagResponses

}

// function for getting the status code from the university api
func uniAPIstatus(w http.ResponseWriter) string {

	// makes an HTTP GET request to an external API endpoint
	response, err := http.Get("http://universities.hipolabs.com/")
	// handles errors if request fails.
	if err != nil {
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
	}
	// returns the HTTP status code of the response
	return response.Status
}

// function for getting the status code from the country api
func countryAPIstatus(w http.ResponseWriter) string {

	// makes an HTTP GET request to an external API endpoint
	response, err := http.Get("https://restcountries.com/")
	// handles errors if request fails.
	if err != nil {
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
	}
	// returns the HTTP status code of the response
	return response.Status
}
