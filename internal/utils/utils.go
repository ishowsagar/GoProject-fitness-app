package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

//! type declaration
//! Envelope --> wrapper for JSON responses, allows flexible key-value pairs
type Envelope map[string]interface{}

//! WriteJson --> standardized JSON response writer used across all handlers
func WriteJson(w http.ResponseWriter, status int, data Envelope) error {
	//* using MarshalIndent for pretty formatted JSON output (easier to read in browser/postman)
	json,err := json.MarshalIndent(data,""," ")
	
	if err != nil {
		return err
	}

	json = append(json, '\n') //* adding newline at end for cleaner terminal output
	w.Header().Set("Content-type","application/json") //* setting response content type
	w.WriteHeader(status) //* HTTP status code (200, 400, 500, etc.)
	w.Write(json) //* writing JSON to response
	return nil
}

//! ReadIDParam --> extracts and validates ID from URL path parameter
//! Used by GET/PUT/DELETE endpoints like /workouts/{id}
func ReadIDParam(r *http.Request) (int64,error) {
	idParam := chi.URLParam(r,"id") //* reads {id} slug from URL path
	
	//? if no id was passed or empty string
	if idParam == "" {
		return 0, errors.New("Invalid id parameter")
	}
	//* convert string to int64 (base 10, 64-bit)
	id,err := strconv.ParseInt(idParam,10,64)
	if err!= nil {
		return 0, errors.New("Invalid id parameter type")
	}
	
	return id,err
}