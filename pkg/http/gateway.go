package http

import (
	"log"
	"os"

	"github.com/gorilla/mux"
)

// GatewayHandler represents the API gateway, in which lies between the client and the services.
type GatewayHandler struct {
	*mux.Router
	Logger *log.Logger
}

// NewGatewayHandler returns a new instance of GatewayHandler
func NewGatewayHandler() *GatewayHandler {
	h := &GatewayHandler{
		Router: mux.NewRouter(),
		Logger: log.New(os.Stderr, "", log.LstdFlags),
	}

	return h
}
