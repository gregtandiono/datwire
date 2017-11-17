package http

import (
	"datwire/pkg/bolt"
	"datwire/pkg/shared"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

// CustomerHandler represents Customer REST APIs
type CustomerHandler struct {
	*mux.Router
	CustomerService *bolt.CustomerService
	Logger          *log.Logger
}

// NewCustomerHandler returns a new instance of CustomerHandler
func NewCustomerHandler() *CustomerHandler {
	h := &CustomerHandler{
		Router:          mux.NewRouter(),
		Logger:          log.New(os.Stderr, "", log.LstdFlags),
		CustomerService: &bolt.CustomerService{},
	}

	return h
}

func (h *CustomerHandler) handleGetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	custID := uuid.FromStringOrNil(vars["id"])
	c, err := h.CustomerService.Customer(custID)
	if err != nil {
		shared.EncodeError(w, err, 400, h.Logger)
	} else {
		shared.EncodeJSON(
			w,
			&shared.ResponseTemplate{Message: "success", Data: c},
			h.Logger,
		)
	}
}

func (h *CustomerHandler) handleGetUsers(w http.ResponseWriter, r *http.Request) {
	customers, err := h.CustomerService.Customers()
	if err != nil {
		shared.EncodeError(w, err, 400, h.Logger)
	} else {
		shared.EncodeJSON(
			w,
			&shared.ResponseTemplate{Message: "success", Data: customers},
			h.Logger,
		)
	}
}
