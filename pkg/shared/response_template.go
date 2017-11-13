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

func encodeError(w http.ResponseWriter, err error, code int, logger *log.Logger) {
	logger.Printf("http error: %s (code=%d)", err, code)

	// Write generic error response.
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(&ResponseTemplate{Message: "fail", Error: err.Error()})
}

func encodeJSON(w http.ResponseWriter, v interface{}, logger *log.Logger) {
	if err := json.NewEncoder(w).Encode(v); err != nil {
		encodeError(w, err, http.StatusInternalServerError, logger)
	}
}
