package rest

import (
	"encoding/json"
	"net/http"
)

// JSON encodes data to rw in JSON format
func JSON(rw http.ResponseWriter, data interface{}) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)

	err := json.NewEncoder(rw).Encode(data)
	if err != nil {
		http.Error(rw, HTTPInternalServerError, http.StatusInternalServerError)
	}
}
