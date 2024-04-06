package main

import (
	"fmt"
	"net/http"

	"github.com/intaek-h/rss/internal/auth"
	"github.com/intaek-h/rss/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (api *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GrabAPIKey(r.Header)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Failed to authenticate: %v", err))
			return
		}

		user, err := api.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get user: %v", err))
			return
		}

		handler(w, r, user)
	}
}
