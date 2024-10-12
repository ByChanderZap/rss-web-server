package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ByChanderZap/rss-web-server/internal/database"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handleCreateFeedFollow(w http.ResponseWriter, r *http.Request, usr database.User) {
	type parameters struct {
		FeedId uuid.UUID `json:"feed_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	_, err = apiCfg.DB.GetFeedById(r.Context(), params.FeedId)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			respondWithError(w, 400, "Feed not exists")
			return
		}
		respondWithError(w, 400, fmt.Sprintf("Could not get feed: %v", err))
		return
	}

	feedFollowsDb, err := apiCfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    usr.ID,
		FeedID:    params.FeedId,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Could not create feed follows: %v", err))
		return
	}

	respondWithJson(w, 201, dbFFToFF(feedFollowsDb))
}

func (apiCfg *apiConfig) handleGetFeedsFollows(w http.ResponseWriter, r *http.Request, usr database.User) {
	ffDb, err := apiCfg.DB.GetFeedsFollows(r.Context(), usr.ID)
	if err != nil {
		respondWithError(w, 400, err.Error())
	}

	respondWithJson(w, 200, dbFeedfsToFeedf(ffDb))
}

func (apiCfg *apiConfig) handleDeleteFeedsFollows(w http.ResponseWriter, r *http.Request, usr database.User) {
	ffIdString := chi.URLParam(r, "feedFollowId")
	ffUuid, err := uuid.Parse(ffIdString)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}

	err = apiCfg.DB.DeleteFeedFollows(r.Context(), database.DeleteFeedFollowsParams{
		ID:     ffUuid,
		UserID: usr.ID,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Could not delete follow: %v", err))
		return
	}

	type pld struct {
		Message string `json:"message"`
	}

	respondWithJson(w, 200, pld{Message: "Follow deleted."})
}
