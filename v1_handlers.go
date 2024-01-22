package main

import (
	"encoding/json"
	"errors"

	"net/http"
	"time"

	"github.com/Nidal-Bakir/rss-aggregator/internal/database"
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

func (conf *apiConfig) GetUserByApiKeyHandler(w http.ResponseWriter, r *http.Request) {

}
