package http

import (
	"bytes"
	"datwire/pkg/shared"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

// AuthGateway is a transport layer to decode/encode JSON request
// and response from the service REST server over http protocol.
type AuthGateway struct {
	*mux.Router
	Logger            *log.Logger
	ServiceConfig     *shared.ServiceConfig
	UserServiceConfig *shared.ServiceConfig
}

// NewAuthGateway returns a new instance of AuthGateway
func NewAuthGateway() *AuthGateway {
	g := &AuthGateway{
		Router:            mux.NewRouter(),
		Logger:            log.New(os.Stderr, "", log.LstdFlags),
		ServiceConfig:     shared.GetEnvironmentVariables("datwire-auth"),
		UserServiceConfig: shared.GetEnvironmentVariables("datwire-users"),
	}

	g.Handle("/auth", http.HandlerFunc(g.handleAuthorization)).Methods("POST")

	return g
}

type loginRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (g *AuthGateway) handleAuthorization(w http.ResponseWriter, r *http.Request) {
	var reqBody *loginRequestBody
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		shared.EncodeError(w, err, 400, g.Logger)
	}

	userID, hashedPassword, err := g.handleCheckIfUserExists("gtandiono")
	if err != nil {
		shared.EncodeError(w, err, 400, g.Logger)
	}

	var responseBody *shared.ResponseTemplate
	client := &http.Client{}

	request, err := http.NewRequest(
		"POST",
		g.ServiceConfig.Address+":"+g.ServiceConfig.Port+"/auth",
		bytes.NewBufferString(`{
			"user_id": "`+userID.String()+`",
			"password": "`+reqBody.Password+`",
			"hashed_password": "`+hashedPassword+`"
		}`),
	)
	if err != nil {
		shared.EncodeError(w, err, 400, g.Logger)
	}

	response, err := client.Do(request)
	if err != nil {
		shared.EncodeError(w, err, 500, g.Logger)
	}

	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if responseBody.Message == "fail" {
		shared.EncodeError(w, errors.New(responseBody.Error), 400, g.Logger)
	} else {
		shared.EncodeJSON(w, responseBody, g.Logger)
	}
}

type checkIfUserExistsResponse struct {
	shared.ResponseTemplate
	Data struct { // override Data struct
		ID             uuid.UUID
		HashedPassword string
	} `json:"data"`
}

func (g *AuthGateway) handleCheckIfUserExists(username string) (userID uuid.UUID, hashedPassword string, err error) {
	var responseBody *checkIfUserExistsResponse
	client := &http.Client{}

	request, err := http.NewRequest(
		"GET",
		g.UserServiceConfig.Address+":"+g.UserServiceConfig.Port+"/users?username="+username,
		nil,
	)

	response, err := client.Do(request)
	if err != nil {
		return
	}

	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if responseBody.Message == "success" {
		userID = responseBody.Data.ID
		hashedPassword = responseBody.Data.HashedPassword
	} else {
		err = errors.New(responseBody.Error)
	}

	return
}
