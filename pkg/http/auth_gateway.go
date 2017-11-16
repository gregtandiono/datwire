package http

import (
	"datwire/pkg/shared"
	"encoding/json"
	"errors"
	"fmt"
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

func (g *AuthGateway) handleAuthorization(w http.ResponseWriter, r *http.Request) {
	userID, _, err := g.handleFindUser("gtandiono")
	if err != nil {
		fmt.Println("HERE????")
		fmt.Println(err)
	}
	fmt.Println("user id", userID)
	// gatewayHandlerFactory(
	// 	"POST",
	// 	g.ServiceConfig.Address+":"+g.ServiceConfig.Port,
	// 	"/auth",
	// 	w, r, g.Logger,
	// )
}

type findUserResponse struct{}

func (g *AuthGateway) handleFindUser(username string) (userID uuid.UUID, hashedPassword string, err error) {
	var responseBody *shared.ResponseTemplate
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
		// type assertion required without a supporting pointer for
		// response body, since `Data` is an undefined interface
		if str, ok := responseBody.Data.(string); ok {
			userID = uuid.FromStringOrNil(str)
		} else {
			err = errors.New("bad data from auth request")
		}
	} else {
		err = errors.New(responseBody.Error)
	}

	return
}
