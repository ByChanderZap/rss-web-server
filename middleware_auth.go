package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ByChanderZap/rss-web-server/internal/auth"
	"github.com/ByChanderZap/rss-web-server/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiCfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetApiKey(r.Header)
		if err != nil {
			respondWithError(w, 403, fmt.Sprintf("Auth error: %v", err))
			return
		}

		usr, err := apiCfg.DB.GetUserByApi(r.Context(), apiKey)
		if err != nil {
			if strings.Contains(err.Error(), "sql: no rows in result set") {
				respondWithError(w, 400, fmt.Sprintf("User not found: %v", err))
				return
			}
			respondWithError(w, 400, fmt.Sprintf("Something went wrong %v", err))
			return
		}

		handler(w, r, usr)
	}
}
