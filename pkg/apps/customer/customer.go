package customer

import (
	"time"

	uuid "github.com/satori/go.uuid"
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
	Industry                 string    `json:"industry"`
	Notes                    string    `json:"notes"`
	CreatedAt                time.Time
	UpdatedAt                time.Time
}

// CustomerService represents customer model CRUD implementation to the BoltDB
type CustomerService interface {
	Customer(id uuid.UUID) (*Customer, error)
	Customers() ([]Customer, error)
	CreateCustomer(c *Customer) error
	UpdateCustomer(customerID uuid.UUID, key, value string) error
	DeleteUser(id uuid.UUID) error
}
