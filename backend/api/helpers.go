package api

import (
	"encoding/json"
	"net/http"
)

func parseJsonBody[T any](r *http.Request) (T, error) {
	var zero T
	if err := json.NewDecoder(r.Body).Decode(&zero); err != nil {
		return zero, err
	}
	return zero, nil
}

func sendJsonBody[T any](w http.ResponseWriter, data T) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func sendJsonError(w http.ResponseWriter, message string, status int) {
	type res struct {
		Error string `json:"error"`
	}

	data := res{
		Error: message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func ptrTo[T any](value T) *T {
	return &value
}
