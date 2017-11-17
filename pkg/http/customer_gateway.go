package http

import (
	"datwire/pkg/shared"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// CustomerGateway is a transport layer to decode/encode JSON request
// and response from the service REST server over http protocol.
type CustomerGateway struct {
	*mux.Router
	Logger        *log.Logger
	ServiceConfig *shared.ServiceConfig
}

// NewCustomerGateway returns a new instance of CustomerGateway.
func NewCustomerGateway() *CustomerGateway {
	g := &CustomerGateway{
		Router:        mux.NewRouter(),
		Logger:        log.New(os.Stderr, "", log.LstdFlags),
		ServiceConfig: shared.GetEnvironmentVariables("datwire-customers"),
	}

	g.Handle("/customers/{id}", http.HandlerFunc(g.handleGetCustomer)).Methods("GET")
	g.Handle("/customers", http.HandlerFunc(g.handleGetCustomers)).Methods("GET")
	g.Handle("/customers", http.HandlerFunc(g.handleCreateCustomer)).Methods("POST")
	g.Handle("/customers/{id}", http.HandlerFunc(g.handleUpdateCustomer)).Methods("PUT")
	g.Handle("/customers/{id}", http.HandlerFunc(g.handleDeleteCustomer)).Methods("DELETE")

	return g
}

func (g *CustomerGateway) handleGetCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	custID := vars["id"]

	gatewayHandlerFactory(
		"GET",
		g.ServiceConfig.Address+":"+g.ServiceConfig.Port,
		"/customers/"+custID,
		w, r, g.Logger,
	)
}

func (g *CustomerGateway) handleGetCustomers(w http.ResponseWriter, r *http.Request) {
	gatewayHandlerFactory(
		"GET",
		g.ServiceConfig.Address+":"+g.ServiceConfig.Port,
		"/customers",
		w, r, g.Logger,
	)
}

func (g *CustomerGateway) handleCreateCustomer(w http.ResponseWriter, r *http.Request) {
	gatewayHandlerFactory(
		"POST",
		g.ServiceConfig.Address+":"+g.ServiceConfig.Port,
		"/customers",
		w, r, g.Logger,
	)
}

func (g *CustomerGateway) handleUpdateCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	custID := vars["id"]

	gatewayHandlerFactory(
		"PUT",
		g.ServiceConfig.Address+":"+g.ServiceConfig.Port,
		"/customers/"+custID,
		w, r, g.Logger,
	)
}

func (g *CustomerGateway) handleDeleteCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	custID := vars["id"]

	gatewayHandlerFactory(
		"DELETE",
		g.ServiceConfig.Address+":"+g.ServiceConfig.Port,
		"/customers/"+custID,
		w, r, g.Logger,
	)
}
