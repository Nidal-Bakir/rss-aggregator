package main

import (
	"encoding/json"
	"errors"
	"strings"

	"net/http"
	"time"

	"github.com/Nidal-Bakir/rss-aggregator/internal/database"
	"github.com/Nidal-Bakir/rss-aggregator/internal/utils"
	"github.com/google/uuid"
)

func errEndpointHandler(w http.ResponseWriter, r *http.Request) {
	responseWithError(w, 444, "Error you dum dum!", errors.New("error"))
}

func (conf *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Name string `json:"name"`
	}

	var p params

	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		responseWithError(w, 400, "Error parsing the json body", err)
		return
	}

	if len(p.Name) <= 2 {
		responseWithError(w, 422, "The field name must be at least 3 char long", errors.New("constraint Error"))
		return
	}

	userParams := database.CreateNewUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      p.Name,
	}

	dbUser, err := conf.DB.CreateNewUser(r.Context(), userParams)

	if err != nil {
		responseWithError(w, 500, "Internal server error, can not create a new user", err)
		return
	}

	responseWithJson(w, 201, toUserModel(dbUser))

}

func (conf *apiConfig) GetUserByApiKeyHandler(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	responseWithJson(w, 200, toUserModel(dbUser))
}

func (apiConfig *apiConfig) CreateFeedHandler(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	type feedParams struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}

	var reqParams feedParams
	err := json.NewDecoder(r.Body).Decode(&reqParams)
	if err != nil {
		responseWithError(w, 400, "malformed json", err)
		return
	}

	name := strings.TrimSpace(reqParams.Name)
	if name == "" {
		responseWithError(w, 400, "name is required", err)
		return
	}

	url := strings.TrimSpace(reqParams.URL)
	if url == "" {
		responseWithError(w, 400, "url is required", err)
		return
	}
	if !utils.IsValidUrl(url) {
		responseWithError(w, 400, "url valid", err)
		return
	}

	dbFeed, err := apiConfig.DB.CreateFeed(r.Context(),
		database.CreateFeedParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Name:      name,
			Url:       url,
			UserID:    dbUser.ID,
		})

	if err != nil {
		responseWithError(w, 500, "con not create the feed", err)
		return
	}

	responseWithJson(w, 201, toFeedModel(dbFeed))
}

func (apiConfig *apiConfig) GetAllFeeds(w http.ResponseWriter, r *http.Request) {
	dbFeeds, err := apiConfig.DB.GetAllFeeds(r.Context())

	if err != nil {
		responseWithError(w, 500, "Can not get the rss feeds,", err)
		return
	}

	publicFeeds := make([]PublicFeedModel, len(dbFeeds))

	for i, v := range dbFeeds {
		publicFeeds[i] = PublicFeedModel{
			ID:        v.ID,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
			Name:      v.Name,
			Url:       v.Url,
		}
	}

	responseWithJson(w, 200, publicFeeds)
}

func (apiConfig *apiConfig) FollowFeedHandler(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	type followFeedParams struct {
		FeedId uuid.UUID `json:"feed_id"`
	}

	var reqParams followFeedParams
	err := json.NewDecoder(r.Body).Decode(&reqParams)
	if err != nil {
		responseWithError(w, 400, "malformed json", err)
		return
	}

	err = apiConfig.DB.FollowFeed(r.Context(),
		database.FollowFeedParams{
			ID:     uuid.New(),
			UserID: dbUser.ID,
			FeedID: reqParams.FeedId,
		})

	if err != nil {
		responseWithError(w, 500, "Can follow feed,", err)
		return
	}

	responseWithJson(w, 201, struct {
		Message string `json:"message"`
	}{Message: "Done!"})
}
