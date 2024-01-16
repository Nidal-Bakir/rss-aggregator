package main

import "net/http"

func errEndpointHandler(w http.ResponseWriter, r *http.Request) {
	responseWithError(w, 444, "Error you dum dum!")
}
