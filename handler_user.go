package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ByChanderZap/rss-web-server/internal/database"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	u, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error creating user: %v", err))
		return
	}

	respondWithJson(w, 201, databaseUserToUser(u))
}

func (apiCfg *apiConfig) handleGetUserByApiKey(w http.ResponseWriter, r *http.Request, usr database.User) {
	respondWithJson(w, 200, databaseUserToUser(usr))
}

func (apiCfg *apiConfig) handleGetPostsForUser(w http.ResponseWriter, r *http.Request, usr database.User) {
	p, err := apiCfg.DB.GetPostsForUser(r.Context(), database.GetPostsForUserParams{
		UserID: usr.ID,
		Limit:  10,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Unable to get posts: %v", err))
		return
	}
	respondWithJson(w, 200, dbPostsToPosts(p))
}
