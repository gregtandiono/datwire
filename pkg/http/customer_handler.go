package http

import (
	"datwire/pkg/apps/customer"
	"datwire/pkg/bolt"
	"datwire/pkg/shared"
	"encoding/json"
	"errors"
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

	h.Handle("/customers", http.HandlerFunc(h.handleCreateCustomer)).Methods("POST")
	h.Handle("/customers/{id}", http.HandlerFunc(h.handleGetCustomer)).Methods("GET")
	h.Handle("/customers", http.HandlerFunc(h.handleGetCustomers)).Methods("GET")
	h.Handle("/customers/{id}", http.HandlerFunc(h.handleUpdateCustomer)).Methods("PUT")
	h.Handle("/customers/{id}", http.HandlerFunc(h.handleDeleteCustomer)).Methods("DELETE")

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
	vars := mux.Vars(r)
	custID := uuid.FromStringOrNil(vars["id"])

	var updateReqBody map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&updateReqBody)
	if err != nil {
		shared.EncodeError(w, err, 400, h.Logger)
		return
	}

	var groupedErr []error

	for k, v := range updateReqBody {
		if str, ok := v.(string); ok {
			err := h.CustomerService.UpdateCustomer(custID, k, str)
			if err != nil {
				groupedErr = append(groupedErr, errors.New("update error at field "+k+" "+err.Error()))
			}
		} else {
			groupedErr = append(groupedErr, errors.New("value not valid"))
		}
	}

	if len(groupedErr) > 0 {
		shared.EncodeError(w, err, 400, h.Logger)
	} else {
		shared.EncodeJSON(w, &shared.ResponseTemplate{Message: "success"}, h.Logger)
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
