package rest

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type Validator interface {
	Validate() url.Values
}

func DecodeAndValidateJSONPayload(w http.ResponseWriter, body io.ReadCloser, v Validator) bool {
	if err := json.NewDecoder(body).Decode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}

	validationMap := make(map[string]interface{})
	if validErrs := v.Validate(); len(validErrs) > 0 {
		validationMap["validationErrors"] = validErrs
		w.Header().Set("Content-type", "applciation/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(validationMap)
	}
	return len(validationMap) == 0
}
