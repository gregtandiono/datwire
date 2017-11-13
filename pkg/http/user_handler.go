package http

import (
	"datwire/pkg/bolt"
	"datwire/pkg/shared"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

// UserHandler represents User REST APIs
type UserHandler struct {
	*mux.Router
	UserService *bolt.UserService
	Logger      *log.Logger
}

// NewUserHandler returns a new instance of UserHandler
func NewUserHandler() *UserHandler {
	h := &UserHandler{
		Router:      mux.NewRouter(),
		Logger:      log.New(os.Stderr, "", log.LstdFlags),
		UserService: &bolt.UserService{},
	}
	return h
}

func (h *UserHandler) handleGetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := uuid.FromStringOrNil(vars["id"])

	u, err := h.UserService.User(userID)
	if err != nil {
		shared.EncodeError(w, err, 400, h.Logger)
	} else {
		shared.EncodeJSON(
			w,
			&shared.ResponseTemplate{Message: "success", Data: u},
			h.Logger,
		)
	}
}

func (h *UserHandler) handleCreateUser(w http.ResponseWriter, r *http.Request) {
}

func (h *UserHandler) handleSetName(w http.ResponseWriter, r *http.Request) {
}

func (h *UserHandler) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
}
