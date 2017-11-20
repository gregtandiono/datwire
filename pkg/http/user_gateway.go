package http

import (
	"datwire/pkg/consul"
	"datwire/pkg/shared"
	"log"
	"net/http"
	"os"

	"github.com/urfave/negroni"

	"github.com/gorilla/mux"
)

// UserGateway is a transport layer to decode/encode JSON request
// and response from the service REST server over http protocol.
type UserGateway struct {
	*mux.Router
	Logger        *log.Logger
	ServiceConfig *shared.ServiceConfig
}

// NewUserGateway returns a new instance of UserGateway.
func NewUserGateway() *UserGateway {
	g := &UserGateway{
		Router:        mux.NewRouter(),
		Logger:        log.New(os.Stderr, "", log.LstdFlags),
		ServiceConfig: shared.GetEnvironmentVariables("datwire-users"),
	}

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

	g.Handle("/users", http.HandlerFunc(g.handleCreateUser)).Methods("POST")

	g.Handle("/users/{id}", negroni.New(
		negroni.HandlerFunc(shared.JWTMiddleware(hashString).HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(g.handleGetUser)),
	)).Methods("GET")

	g.Handle("/users", negroni.New(
		negroni.HandlerFunc(shared.JWTMiddleware(hashString).HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(g.handleGetUsers)),
	)).Methods("GET")

	g.Handle("/users/{id}", negroni.New(
		negroni.HandlerFunc(shared.JWTMiddleware(hashString).HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(g.handleSetName)),
	)).Methods("PUT")

	g.Handle("/users/{id}", negroni.New(
		negroni.HandlerFunc(shared.JWTMiddleware(hashString).HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(g.handleDeleteUser)),
	)).Methods("DELETE")

	return g
}

func (g *UserGateway) handleGetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	gatewayHandlerFactory(
		"GET",
		g.ServiceConfig.Address+":"+g.ServiceConfig.Port,
		"/users/"+userID,
		w, r, g.Logger,
	)

}

func (g *UserGateway) handleGetUsers(w http.ResponseWriter, r *http.Request) {
	vals := r.FormValue("username")

	if vals != "" {
		gatewayHandlerFactory(
			"GET",
			g.ServiceConfig.Address+":"+g.ServiceConfig.Port,
			"/users?username="+vals,
			w, r, g.Logger,
		)
	} else {
		gatewayHandlerFactory(
			"GET",
			g.ServiceConfig.Address+":"+g.ServiceConfig.Port,
			"/users",
			w, r, g.Logger,
		)
	}
}

func (g *UserGateway) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	gatewayHandlerFactory(
		"POST",
		g.ServiceConfig.Address+":"+g.ServiceConfig.Port,
		"/users",
		w, r, g.Logger,
	)
}

func (g *UserGateway) handleSetName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	gatewayHandlerFactory(
		"PUT",
		g.ServiceConfig.Address+":"+g.ServiceConfig.Port,
		"/users/"+userID,
		w, r, g.Logger,
	)
}

func (g *UserGateway) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	gatewayHandlerFactory(
		"DELETE",
		g.ServiceConfig.Address+":"+g.ServiceConfig.Port,
		"/users/"+userID,
		w, r, g.Logger,
	)
}
