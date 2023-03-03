package main

import (
	"Assignment1/constants"
	"Assignment1/handler"
	"log"
	"net/http"
	"os"
)

func main() {

	// Handle port assignment (either based on environment variable, or local override)
	log.Println("Starting main function.")
	port := os.Getenv("PORT")
	if port == "" {
		log.Println("$PORT has not been set. Default: 8080")
		port = "8080"
	}

	// Set up handler endpoints
	http.HandleFunc("/", handler.DefaultHandler)
	http.HandleFunc(constants.UNI_INFO_PATH, handler.UniInfoHandler)
	http.HandleFunc(constants.NEIGHBOR_UNIS_PATH, handler.NeighbourUnisHandler)
	http.HandleFunc(constants.DIAG_PATH, handler.DiagHandler)

	// Start server
	log.Println("Starting server on port " + port + " ...")
	port = ":" + port
	log.Fatal(http.ListenAndServe(port, nil))

}
