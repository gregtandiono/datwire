package bolt

import (
	"datwire/pkg/apps/customer"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/boltdb/bolt"
	uuid "github.com/satori/go.uuid"
)

// CustomerService represents a client to the underlying BoltDB data store.
type CustomerService struct {
	db *bolt.DB
}

// Open creates a connection to bolt db and inits a user bucket
func (s *CustomerService) Open() error {
	db, err := bolt.Open("customer.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	s.db = db
	// initialize bucket if it does not exist
	return s.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("customers"))
		if err != nil {
			return err
		}
		return nil
	})
}

// Close closes the underlying bolt db
func (s *CustomerService) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

// Customer returns one user record
func (s *CustomerService) Customer(id uuid.UUID) (*customer.Customer, error) {
	var customer *customer.Customer
	err := s.db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte("customers"))
		v := bkt.Get([]byte(id.String()))
		if v == nil {
			return errors.New("customer does not exist")
		} else if err := json.Unmarshal(v, &customer); err != nil {
			return err
		}
		return nil
	})

	return customer, err
}

// Customers returns an array of users
func (s *CustomerService) Customers() ([]customer.Customer, error) {
	var cust customer.Customer
	var customers []customer.Customer
	err := s.db.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte("customers"))
		c := bkt.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			err := json.Unmarshal(v, &cust)
			if err != nil {
				return err
			}
			customers = append(customers, cust)
		}
		return nil
	})

	return customers, err
}

// CreateCustomer saves a new customer record to the DB
func (s *CustomerService) CreateCustomer(c *customer.Customer) error {
	if c == nil {
		return errors.New("customer cannot be empty")
	} else if c.ID == uuid.FromStringOrNil("") {
		return errors.New("customer ID is required")
	}

	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()

	return s.db.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte("customers"))

		// Check if user already exists
		if v := bkt.Get([]byte(c.ID.String())); v != nil {
			return errors.New("customer already exist")
		}

		if buf, err := json.Marshal(c); err != nil {
			return err
		} else if err := bkt.Put([]byte(c.ID.String()), buf); err != nil {
			return err
		}
		return nil
	})
}

// UpdateCustomer updates a specific key/field in an existing customer record
func (s *CustomerService) UpdateCustomer(custID uuid.UUID, key, value string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte("customers"))
		var cust map[string]interface{}

		if v := bkt.Get([]byte(custID.String())); v == nil {
			return errors.New("customer does not exist")
		} else if err := json.Unmarshal(v, &cust); err != nil {
			return err
		}

		cust[key] = value
		cust["updated_at"] = time.Now().String()

		if buf, err := json.Marshal(cust); err != nil {
			return err
		} else if err := bkt.Put([]byte(custID.String()), buf); err != nil {
			return err
		}

		return nil
	})
}

// DeleteCustomer removes an existing customer from the bolt db
func (s *CustomerService) DeleteCustomer(id uuid.UUID) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte("customers"))
		if err := bkt.Delete([]byte(id.String())); err != nil {
			return err
		}
		return nil
	})
}
