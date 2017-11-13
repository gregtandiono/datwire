package http

import (
	"datwire/pkg/bolt"
	"datwire/pkg/shared"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// AuthHandler represents Auth REST APIs
type AuthHandler struct {
	*mux.Router
	AuthService *bolt.AuthService
	Logger      *log.Logger
}

// NewAuthHandler returns a new instance of AuthHandler
func NewAuthHandler() *AuthHandler {
	h := &AuthHandler{
		Router:      mux.NewRouter(),
		Logger:      log.New(os.Stderr, "", log.LstdFlags),
		AuthService: &bolt.AuthService{},
	}

	return h
}

type authRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *AuthHandler) handleAuthorization(w http.ResponseWriter, r *http.Request) {
	var authReqBody *authRequestBody
	if err := json.NewDecoder(r.Body).Decode(&authReqBody); err != nil {
		shared.EncodeError(w, err, 400, h.Logger)
	}
}
