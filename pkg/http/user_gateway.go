package http

import (
	"datwire/pkg/shared"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// UserGateway is a transport layer to decode/encode JSON request
// and response from the service REST server over http protocol.
type UserGateway struct {
	*mux.Router
	Logger *log.Logger
}

// NewUserGateway returns a new instance of UserGateway.
func NewUserGateway() *UserGateway {
	g := &UserGateway{
		Router: mux.NewRouter(),
		Logger: log.New(os.Stderr, "", log.LstdFlags),
	}

	g.Handle("/users/{id}", http.HandlerFunc(g.handleGetUser)).Methods("GET")
	g.Handle("/users", http.HandlerFunc(g.handleCreateUser)).Methods("POST")

	return g
}

func (g *UserGateway) handleGetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	gatewayHandlerFactory(
		"GET",
		"http://localhost:1337",
		"/users/"+userID,
		w, r, g.Logger,
	)

}

// func (g *UserGateway) handleGetUsers(w http.ResponseWriter, r *http.Request) {
// 	vals := r.FormValue("username")
// 	if vals != "" {
// 		if userID, err := h.UserService.CheckIfUserExists(vals); err != nil {
// 			shared.EncodeError(w, err, 400, h.Logger)
// 		} else {
// 			shared.EncodeJSON(w, &shared.ResponseTemplate{Message: "success", Data: userID}, h.Logger)
// 		}

// 	} else {
// 		if users, err := h.UserService.Users(); err != nil {
// 			shared.EncodeError(w, err, 400, h.Logger)
// 		} else {
// 			shared.EncodeJSON(
// 				w,
// 				&shared.ResponseTemplate{Message: "success", Data: users},
// 				h.Logger,
// 			)
// 		}
// 	}
// }

func (g *UserGateway) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	gatewayHandlerFactory(
		"POST",
		"http://localhost:1337",
		"/users",
		w, r, g.Logger,
	)
}

func gatewayHandlerFactory(reqMethod, targetBaseURL, endpoint string, w http.ResponseWriter, r *http.Request, l *log.Logger) {
	var responseBody *shared.ResponseTemplate
	switch reqMethod {
	case "POST":
		response, err := http.Post(targetBaseURL+endpoint, "application/json", r.Body)
		if err != nil {
			shared.EncodeError(w, err, 500, l)
		}
		err = json.NewDecoder(response.Body).Decode(&responseBody)
		if responseBody.Message == "fail" {
			shared.EncodeError(w, errors.New(responseBody.Error), 400, l)
		} else {
			shared.EncodeJSON(w, responseBody, l)
		}
	case "GET":
		response, err := http.Get(targetBaseURL + endpoint)
		if err != nil {
			shared.EncodeError(w, err, 500, l)
		}
		err = json.NewDecoder(response.Body).Decode(&responseBody)
		if responseBody.Message == "fail" {
			shared.EncodeError(w, errors.New(responseBody.Error), 400, l)
		} else {
			shared.EncodeJSON(w, responseBody, l)
		}
	}
}

// type nameSetterReq struct {
// 	Name string `json:"name"`
// }

// func (g *UserGateway) handleSetName(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	userID := uuid.FromStringOrNil(vars["id"])

// 	var reqBody *nameSetterReq
// 	err := json.NewDecoder(r.Body).Decode(&reqBody)
// 	if err != nil {
// 		shared.EncodeError(w, err, 400, h.Logger)
// 		return
// 	}

// 	err = h.UserService.SetName(userID, reqBody.Name)
// 	if err != nil {
// 		shared.EncodeError(w, err, 400, h.Logger)
// 	} else {
// 		shared.EncodeJSON(w, &shared.ResponseTemplate{Message: "success"}, h.Logger)
// 	}
// }

// func (g *UserGateway) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	userID := uuid.FromStringOrNil(vars["id"])

// 	if err := h.UserService.DeleteUser(userID); err != nil {
// 		shared.EncodeError(w, err, 400, h.Logger)
// 	} else {
// 		shared.EncodeJSON(w, &shared.ResponseTemplate{Message: "success"}, h.Logger)
// 	}
// }
