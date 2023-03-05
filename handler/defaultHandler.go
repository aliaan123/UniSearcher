package handler

import (
	"Assignment1/constants"
	"fmt"
	"net/http"
)

// DefaultHandler function, which just prints info about the service
func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi! Here is some info in order to get started using this service:\n\t")
	fmt.Fprintf(w, "\n\t For information about the country a specific university recides in; "+constants.UNI_INFO_PATH+"{your_search_name}")
	fmt.Fprintf(w, "\n\t For information about universities in neightbouring countries with the same name; "+constants.NEIGHBOR_UNIS_PATH+"{country_search_name}"+"/{your_search_name")
	fmt.Fprintf(w, "\n\t For information about the diagnosis of the API; "+constants.DIAG_PATH+"")
}
