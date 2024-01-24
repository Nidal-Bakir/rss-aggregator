package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/Nidal-Bakir/rss-aggregator/internal/auth"
	"github.com/Nidal-Bakir/rss-aggregator/internal/database"
)

type authHandlerFunc func(w http.ResponseWriter, r *http.Request, dbUser database.User)

func (apiConfig *apiConfig) authMiddleware(handler authHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetApiKey(r.Header)

		if err != nil {
			responseWithError(w, 401, fmt.Sprint(err), err)
			return
		}

		dbUser, err := apiConfig.DB.GetUserByApiKey(r.Context(), apiKey)

		if err != nil {
			if err == sql.ErrNoRows {
				responseWithError(w, 404, "User not found", err)
				return
			}

			responseWithError(w, 500, "Can not get the user", err)
			return
		}

		handler(w, r, dbUser)
	}
}
