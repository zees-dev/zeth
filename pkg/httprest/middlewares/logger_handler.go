package middlewares

import (
	"encoding/json"
	"net/http"

	"github.com/zees-dev/zeth/pkg/httprest/rest"
)

type (
	// LoggerHandler defines a HTTP handler that includes a HandlerError return pointer
	LoggerHandler func(http.ResponseWriter, *http.Request) *rest.HandlerError

	// errorResponse represents an error response returned to the client
	errorResponse struct {
		Message string `json:"message,omitempty"`
		Details string `json:"details,omitempty"`
	}
)

func (handler LoggerHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	err := handler(rw, r)
	if err != nil {
		writeErrorResponse(rw, err)
	}
}

func writeErrorResponse(rw http.ResponseWriter, err *rest.HandlerError) {
	// log.Debug("http error: %s (err=%s) (code=%d)\n", err.Message, err.Err, err.StatusCode)
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(err.StatusCode)
	json.NewEncoder(rw).Encode(&errorResponse{Message: err.Message, Details: err.Err.Error()})
}
