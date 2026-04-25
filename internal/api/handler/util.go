package handler

import (
	"encoding/json"
	"net/http"
)

func RespondJson(w http.ResponseWriter, r *http.Request, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
