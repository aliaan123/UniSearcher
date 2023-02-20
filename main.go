package main

import (
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
	log.Println("Starting defaultHandler function.")
	http.HandleFunc("/", handler.DefaultHandler)
	log.Println("Starting uniInfoHandler function.")
	http.HandleFunc(handler.UNI_INFO_PATH, handler.UniInfoHandler)
	//http.HandleFunc(handler.NEIGHBOR_UNIS_PATH, handler.NeighbourUnisHandler)
	//http.HandleFunc(handler.DIAG_PATH, handler.diagHandler)

	// Start server
	log.Println("Starting server on port " + port + " ...")
	log.Fatal(http.ListenAndServe(":"+port, nil))

}

/*

På slutten så må dette leveres
 - gitlab (NTNU) link
 - github private repo link
 - render project link

*/
