package http

import (
	"datwire/pkg/shared"
	"log"
	"os"

	"github.com/gorilla/mux"
)

type AuthGateway struct {
	*mux.Router
	Logger        *log.Logger
	ServiceConfig *shared.ServiceConfig
}

func NewAuthGateway() *AuthGateway {
	g := &AuthGateway{
		Router:        mux.NewRouter(),
		Logger:        log.New(os.Stderr, "", log.LstdFlags),
		ServiceConfig: shared.GetEnvironmentVariables("datwire-auth"),
	}

	return g
}
