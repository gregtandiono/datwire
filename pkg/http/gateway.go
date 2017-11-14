package http

import (
	"datwire/pkg/shared"
	"encoding/json"
	"errors"
	"log"
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

func gatewayHandlerFactory(reqMethod, targetBaseURL, endpoint string, w http.ResponseWriter, r *http.Request, l *log.Logger) {
	var responseBody *shared.ResponseTemplate
	client := &http.Client{}

	request, err := http.NewRequest(reqMethod, targetBaseURL+endpoint, r.Body)
	if err != nil {
		shared.EncodeError(w, err, 500, l)
	}

	response, err := client.Do(request)
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
