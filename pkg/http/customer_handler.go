package http

import (
	"datwire/pkg/apps/customer"
	"datwire/pkg/bolt"
	"datwire/pkg/shared"
	"encoding/json"
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

func (h *CustomerHandler) handleGetCustomer(w http.ResponseWriter, r *http.Request) {
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

func (h *CustomerHandler) handleGetCustomers(w http.ResponseWriter, r *http.Request) {
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

func (h *CustomerHandler) handleCreateCustomer(w http.ResponseWriter, r *http.Request) {
	var cust *customer.Customer
	if err := json.NewDecoder(r.Body).Decode(&cust); err != nil {
		shared.EncodeError(w, err, 400, h.Logger)
	} else if err := h.CustomerService.CreateCustomer(cust); err != nil {
		shared.EncodeError(w, err, 400, h.Logger)
	} else {
		shared.EncodeJSON(w, &shared.ResponseTemplate{Message: "success"}, h.Logger)
	}
}

func (h *CustomerHandler) handleUpdateCustomer(w http.ResponseWriter, r *http.Request) {
	var updateReqBody map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&updateReqBody)
	if err != nil {
		shared.EncodeError(w, err, 400, h.Logger)
	}
}

func (h *CustomerHandler) handleDeleteCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	custID := uuid.FromStringOrNil(vars["id"])

	if err := h.CustomerService.DeleteCustomer(custID); err != nil {
		shared.EncodeError(w, err, 400, h.Logger)
	} else {
		shared.EncodeJSON(w, &shared.ResponseTemplate{Message: "success"}, h.Logger)
	}
}
