package main

import "github.com/go-chi/chi/v5"

func initV1Router(apiConfig apiConfig) (v1Router *chi.Mux) {
	r := chi.NewRouter()

	r.Get("/err", errEndpointHandler)
	
	r.Post("/users", apiConfig.createUserHandler)
	r.Get("/me", apiConfig.authMiddleware(apiConfig.GetUserByApiKeyHandler))

	r.Post("/feed", apiConfig.authMiddleware(apiConfig.CreateFeedHandler))
	r.Get("/all-feeds", apiConfig.GetAllFeeds)
	r.Get("/follow-feed", apiConfig.authMiddleware(apiConfig.FollowFeedHandler))

	return r
}
