package main

import "github.com/go-chi/chi/v5"

func initV1Router(apiConfig apiConfig) (v1Router *chi.Mux) {
	r := chi.NewRouter()

	r.Get("/err", errEndpointHandler)

	r.Post("/users", apiConfig.createUserHandler)
	r.Get("/me", apiConfig.authMiddleware(apiConfig.GetUserByApiKeyHandler))

	r.Post("/feed", apiConfig.authMiddleware(apiConfig.CreateFeedHandler))
	r.Get("/all-feeds", apiConfig.GetAllFeeds)

	r.Get("/feed-follows", apiConfig.authMiddleware(apiConfig.FeedFollowsHandler))

	r.Post("/follow-feed", apiConfig.authMiddleware(apiConfig.FollowFeedHandler))
	r.Delete("/follow-feed/{id}", apiConfig.authMiddleware(apiConfig.UnfollowFeedHandler))

	return r
}
