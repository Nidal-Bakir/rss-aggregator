package main

import "github.com/go-chi/chi/v5"

func initV1Router(apiConfig apiConfig) (v1Router *chi.Mux) {
	r := chi.NewRouter()

	r.Get("/err", errEndpointHandler)
	r.Post("/users", apiConfig.createUserHandler)
	r.Get("/me", apiConfig.GetUserByApiKeyHandler)

	return r
}
