package main

import (
	"net/http"
)

func handleRediness(w http.ResponseWriter, r *http.Request) {
	respondWithJson(w, 200, struct{}{})
}
