package http

import (
	"datwire/pkg/bolt"
	"datwire/pkg/consul"
	"datwire/pkg/shared"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

// AuthHandler represents Auth REST APIs
type AuthHandler struct {
	*mux.Router
	AuthService *bolt.AuthService
	Logger      *log.Logger
}

// NewAuthHandler returns a new instance of AuthHandler
func NewAuthHandler() *AuthHandler {
	var hashString string
	if env := os.Getenv("ENV"); env != "TEST" {
		consuld := consul.NewConsuld(nil)
		hash, err := consuld.GetKV("datwire/config/hashString", nil)
		hashString = hash
		if err != nil {
			log.Fatal(err)
		}
	} else {
		hashString = "869826e158da8666906ec2681b19b96b729665fd2fae1328ace29171a1e8b3e2" // just for testing purposes
	}
	h := &AuthHandler{
		Router: mux.NewRouter(),
		Logger: log.New(os.Stderr, "", log.LstdFlags),
		AuthService: &bolt.AuthService{
			Hash: hashString,
		},
	}

	h.Handle("/auth", http.HandlerFunc(h.handleAuthorization)).Methods("POST")

	return h
}

type authRequestBody struct {
	Password       string    `json:"password"`
	HashedPassword string    `json:"hashed_password"`
	UserID         uuid.UUID `json:"user_id"`
}

func (h *AuthHandler) handleAuthorization(w http.ResponseWriter, r *http.Request) {
	var authReqBody *authRequestBody
	if err := json.NewDecoder(r.Body).Decode(&authReqBody); err != nil {
		shared.EncodeError(w, err, 400, h.Logger)
		return
	}
	if a, err := h.AuthService.Authorize(
		authReqBody.Password,
		authReqBody.HashedPassword,
		authReqBody.UserID,
	); err != nil {
		shared.EncodeError(w, err, 400, h.Logger)
	} else {
		shared.EncodeJSON(w, &shared.ResponseTemplate{Message: "success", Data: a}, h.Logger)
	}
}
