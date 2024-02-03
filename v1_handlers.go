package main

import (
	"encoding/json"
	"errors"
	"strconv"

	"strings"

	"net/http"
	"time"

	"github.com/Nidal-Bakir/rss-aggregator/internal/database"
	"github.com/Nidal-Bakir/rss-aggregator/internal/utils"
	"github.com/go-chi/chi/v5"
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

func (apiConfig *apiConfig) FeedFollowsHandler(w http.ResponseWriter, r *http.Request, dbUser database.User) {

	dbFeeds, err := apiConfig.DB.GetFeedFollows(r.Context(), dbUser.ID)

	if err != nil {
		responseWithError(w, 500, "Can not get the feed", err)
		return
	}

	feedFollowsSlice := make([]FeedFollowsModel, len(dbFeeds))

	for i, v := range dbFeeds {
		feedFollowsSlice[i] = FeedFollowsModel{
			ID:        v.ID,
			FeedId:    v.FeedID,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
			Name:      v.Name,
			Url:       v.Url,
		}
	}

	responseWithJson(w, 200, feedFollowsSlice)
}

func (apiConfig *apiConfig) UnfollowFeedHandler(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	strId := chi.URLParam(r, "id")
	id, err := uuid.Parse(strId)
	if err != nil {
		responseWithError(w, 400, "malformed id", err)
		return
	}

	sqlResult, err := apiConfig.DB.UnfollowFeed(r.Context(),
		database.UnfollowFeedParams{
			ID:     id,
			UserID: dbUser.ID,
		})
	if err != nil {
		responseWithError(w, 500, "Can not unfollow feed", err)
		return
	}

	rowsAffected, err := sqlResult.RowsAffected()
	if err != nil {
		responseWithError(w, 500, "Can not unfollow feed", err)
		return
	}
	if rowsAffected == 0 {
		responseWithError(w, 404, "No feed to unfollow", err)
		return
	}

	responseWithJson(w, 200, struct {
		Message string `json:"message"`
	}{Message: "Done!"})

}

func (apiConfig *apiConfig) PostsForFollowedFeedsHandler(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	params := r.URL.Query().Get

	var pageSize int
	if val, err := strconv.Atoi(params("page_size")); err == nil {
		pageSize = val
	}
	pageSize = utils.Clamp(pageSize, 1, 40)

	var page int
	if val, err := strconv.Atoi(params("page")); err == nil {
		page = val
	}

	dbPosts, err := apiConfig.DB.GetPostsForFollowedFeed(
		r.Context(),
		database.GetPostsForFollowedFeedParams{
			UserID: dbUser.ID,
			Offset: int32(page * pageSize),
			Limit:  int32(pageSize + 1),
		},
	)

	if err != nil {
		responseWithError(w, 500, "Error while gating posts", err)
		return
	}

	publicPostsSlices := make([]PostModel, len(dbPosts))
	for i, post := range dbPosts {
		publicPostsSlices[i] = toPostModel(post)
	}

	type Payload struct {
		Data        []PostModel `json:"data"`
		CanLoadMore bool        `json:"can_load_more"`
	}

	canLoadMore := false
	if len(publicPostsSlices) > pageSize {
		publicPostsSlices = publicPostsSlices[:len(publicPostsSlices)-1]
		canLoadMore = true
	}

	responseWithJson(w, 200, Payload{Data: publicPostsSlices, CanLoadMore: canLoadMore})

}
