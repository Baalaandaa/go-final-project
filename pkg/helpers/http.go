package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type Error struct {
	Message string `json:"message" example:"error message"`
}

func WriteResponse(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	_, _ = w.Write([]byte(message))
	_, _ = w.Write([]byte("\n"))
}

func WriteJSONResponse(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		WriteResponse(w, http.StatusInternalServerError, fmt.Sprintf("can't marshal data: %s", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")

	WriteResponse(w, status, string(response))
}

func WriteError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrForbidden):
		WriteJSONResponse(w, http.StatusForbidden, Error{Message: err.Error()})
	case errors.Is(err, ErrNotFound):
		WriteJSONResponse(w, http.StatusNotFound, Error{Message: err.Error()})
	default:
		WriteJSONResponse(w, http.StatusInternalServerError, Error{Message: err.Error()})
	}
}
