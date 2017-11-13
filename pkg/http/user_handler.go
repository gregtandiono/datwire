package http

import (
	"datwire/pkg/bolt"
	"log"

	"github.com/gorilla/mux"
)

// UserHandler represents User REST APIs
type UserHandler struct {
	*mux.Router
	UserService *bolt.UserService
	Logger      *log.Logger
}

// NewUserHandler returns a new instance of UserHandler
func NewUserHandler() *UserHandler {
	h := &UserHandler{}
	return h
}
