package http

import (
	"datwire/pkg/shared"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// AuthGateway is a transport layer to decode/encode JSON request
// and response from the service REST server over http protocol.
type AuthGateway struct {
	*mux.Router
	Logger        *log.Logger
	ServiceConfig *shared.ServiceConfig
}

// NewAuthGateway returns a new instance of AuthGateway
func NewAuthGateway() *AuthGateway {
	g := &AuthGateway{
		Router:        mux.NewRouter(),
		Logger:        log.New(os.Stderr, "", log.LstdFlags),
		ServiceConfig: shared.GetEnvironmentVariables("datwire-auth"),
	}

	g.Handle("/auth", http.HandlerFunc(g.handleAuthorization)).Methods("POST")

	return g
}

func (g *AuthGateway) handleAuthorization(w http.ResponseWriter, r *http.Request) {
	gatewayHandlerFactory(
		"POST",
		g.ServiceConfig.Address+":"+g.ServiceConfig.Port,
		"/auth",
		w, r, g.Logger,
	)
}
