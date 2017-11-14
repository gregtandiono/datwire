package shared

import (
	"encoding/json"
	"log"
	"net/http"
)

// ResponseTemplate represents generic json response
type ResponseTemplate struct {
	Message string      `json:"message"`
	Error   string      `json:"error"`
	Data    interface{} `json:"data"`
}

// EncodeError represents JSON response builder for error state
func EncodeError(w http.ResponseWriter, err error, code int, logger *log.Logger) {
	logger.Printf("http error: %s (code=%d)", err, code)

	// Write generic error response.
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(&ResponseTemplate{Message: "fail", Error: err.Error()})
}

// GatewayErrorEncoder is an api gateway wrapper around the EncodeError method and
// takes an interface as a param for flexibility
func GatewayErrorEncoder(w http.ResponseWriter, v interface{}, code int, logger *log.Logger) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}

// EncodeJSON represents a generic JSON response builder
func EncodeJSON(w http.ResponseWriter, v interface{}, logger *log.Logger) {
	if err := json.NewEncoder(w).Encode(v); err != nil {
		EncodeError(w, err, http.StatusInternalServerError, logger)
	}
}
