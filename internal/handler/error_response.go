package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

func respondWithValidationError(w http.ResponseWriter, err error) {
	for _, err := range err.(validator.ValidationErrors) {
		switch err.Tag() {
		case "required":
			respondWithError(w, http.StatusBadRequest, err.Field()+" is missing")
		case "datetime":
			respondWithError(w, http.StatusBadRequest, "wrong date format for "+err.Field())
		case "gte":
			respondWithError(w, http.StatusBadRequest, err.Field()+" must be at least "+err.Param())
		}
	}
}
