package main

import (
	"net/http"
	"fmt"
	"github.com/Puneet-Pal-Singh/go-rssfeed/internal/database"
	"github.com/Puneet-Pal-Singh/go-rssfeed/internal/auth"
)

type authehandler func (http.ResponseWriter, *http.Request, database.User)

func (apiCfg *apiConfig) middlewareAuth(handler authehandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil{
			respondWithError(w, 403, fmt.Sprintf("Auth Error: %v", err))
			return
		}

		user, err := apiCfg.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil{
			respondWithError(w, 400, fmt.Sprintf("Couldn't get user: %v", err))
			return
		}

		handler(w, r, user)
	}
}