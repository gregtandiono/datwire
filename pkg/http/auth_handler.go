package http

import (
	"datwire/pkg/bolt"
	"log"
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
