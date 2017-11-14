package http

import (
	"net/http"
	"strings"
)

// GatewayHandler represents the API gateway, in which lies between the client and the services.
type GatewayHandler struct {
	UserGateway *UserGateway
}

// NewGatewayHandler returns a new instance of GatewayHandler
func NewGatewayHandler() *GatewayHandler {
	return &GatewayHandler{
		UserGateway: NewUserGateway(),
	}
}

func (g *GatewayHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch true {
	case strings.HasPrefix(r.URL.Path, "/users"):
		g.UserGateway.ServeHTTP(w, r)
	}
}
