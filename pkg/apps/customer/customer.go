package customer

import (
	"time"

	"github.com/pborman/uuid"
)

// Customer represents customer data model
type Customer struct {
	ID                       uuid.UUID `json:"id"`
	Name                     string    `json:"name"`
	Address                  string    `json:"address"`
	Telephone                string    `json:"telephone"`
	ProcurementPIC           string    `json:"procurement_pic"`
	ProcurementContactNumber string    `json:"procurement_contact_number"`
	OperationsPIC            string    `json:"operations_pic"`
	OperationsContactNumber  string    `json:"operations_contact_number"`
	CreatedAt                time.Time
	UpdatedAt                time.Time
}

// CustomerService represents customer model CRUD implementation to the BoltDB
type CustomerService interface{}
