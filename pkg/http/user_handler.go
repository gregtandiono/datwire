package http

import (
	"datwire/pkg/apps/user"
	"datwire/pkg/bolt"
	"datwire/pkg/shared"
	"encoding/json"
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

	h.Handle("/users", http.HandlerFunc(h.handleCreateUser)).Methods("POST")
	h.Handle("/users/{id}", http.HandlerFunc(h.handleGetUser)).Methods("GET")
	h.Handle("/users", http.HandlerFunc(h.handleGetUsers)).Methods("GET")
	h.Handle("/users/{id}", http.HandlerFunc(h.handleSetName)).Methods("PUT")
	h.Handle("/users/{id}", http.HandlerFunc(h.handleDeleteUser)).Methods("DELETE")

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

func (h *UserHandler) handleGetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.UserService.Users()
	if err != nil {
		shared.EncodeError(w, err, 400, h.Logger)
	} else {
		shared.EncodeJSON(
			w,
			&shared.ResponseTemplate{Message: "success", Data: users},
			h.Logger,
		)
	}
}

func (h *UserHandler) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var user *user.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		shared.EncodeError(w, err, 400, h.Logger)
		return
	} else if err := h.UserService.CreateUser(user); err != nil {
		shared.EncodeError(w, err, 400, h.Logger)
	} else {
		shared.EncodeJSON(w, &shared.ResponseTemplate{Message: "success"}, h.Logger)
	}
}

type nameSetterReq struct {
	Name string `json:"name"`
}

func (h *UserHandler) handleSetName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := uuid.FromStringOrNil(vars["id"])

	var reqBody *nameSetterReq
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		shared.EncodeError(w, err, 400, h.Logger)
		return
	}

	err = h.UserService.SetName(userID, reqBody.Name)
	if err != nil {
		shared.EncodeError(w, err, 400, h.Logger)
	} else {
		shared.EncodeJSON(w, &shared.ResponseTemplate{Message: "success"}, h.Logger)
	}
}

func (h *UserHandler) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := uuid.FromStringOrNil(vars["id"])

	if err := h.UserService.DeleteUser(userID); err != nil {
		shared.EncodeError(w, err, 400, h.Logger)
	} else {
		shared.EncodeJSON(w, &shared.ResponseTemplate{Message: "success"}, h.Logger)
	}
}
