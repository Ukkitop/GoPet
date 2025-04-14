package apiUtils

import (
	"encoding/json"
	"net/http"
)

type Error struct {
	Code    int
	Message string
}

func WriteError(w http.ResponseWriter, message string, code int) {
	resp := Error{
		Code:    code,
		Message: message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	json.NewEncoder(w).Encode(resp)
}

func WriteSuccess[T any](w http.ResponseWriter, message T) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(message)
}

var (
	RequestErrorHandler = func(w http.ResponseWriter, err error) {
		WriteError(w, err.Error(), http.StatusBadRequest)
	}
	InternalErrorHandler = func(w http.ResponseWriter) {
		WriteError(w, "Internal server error", http.StatusInternalServerError)
	}
	SuccessResponseHandler = func(w http.ResponseWriter) {
		WriteSuccess(w, "Success")
	}
)
