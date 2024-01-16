package main

import "github.com/go-chi/chi/v5"

func initV1Router() (v1Router *chi.Mux) {
	r := chi.NewRouter()
	
	r.Get("/err",errEndpointHandler)


	return r
}