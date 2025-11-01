package api

import (
	"encoding/json"
	"net/http"
)

func HandlePing(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string {
		"message":"pong",
	})
}

func HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}