package handlers

import (
	"encoding/json"
	"io"
	"net/http"
)

func decodeRequestBody[T any](w http.ResponseWriter, body io.Reader) (*T, bool) {
	var value T

	if err := json.NewDecoder(body).Decode(&value); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return nil, false
	}

	return &value, true
}

func encodeJSONResponse[T any](w http.ResponseWriter, status int, result T) bool {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return false
	}

	return true
}
