package bolt

import (
	"datwire/pkg/apps/user"
	"encoding/json"
	"errors"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/boltdb/bolt"
	uuid "github.com/satori/go.uuid"
)

// UserService represents a client to the underlying BoltDB data store.
type UserService struct {
	db *bolt.DB
}

var _ user.UserService = &UserService{}

// Open creates a connection to bolt db and inits a user bucket
func (s *UserService) Open() error {
	db, err := bolt.Open("user.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	s.db = db
	// initialize bucket if it does not exist
	return s.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("users"))
		if err != nil {
			return err
		}
		return nil
	})
}

// Close closes the underlying bolt db
func (s *UserService) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

// User returns one user record
func (s *UserService) User(id uuid.UUID) (*user.User, error) {
	var user *user.User
	err := s.db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte("users"))
		v := bkt.Get([]byte(id.String()))
		if v == nil {
			return errors.New("user does not exist")
		} else if err := json.Unmarshal(v, &user); err != nil {
			return err
		}
		return nil
	})

	return user, err
}

// Users returns an array of users
func (s *UserService) Users() ([]user.User, error) {
	var person user.User
	var users []user.User
	err := s.db.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte("users"))
		c := bkt.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			err := json.Unmarshal(v, &person)
			if err != nil {
				return err
			}
			users = append(users, person)
		}
		return nil
	})

	return users, err
}

// CreateUser saves a new user record to the DB
func (s *UserService) CreateUser(u *user.User) error {
	if u == nil {
		return errors.New("user cannot be empty")
	} else if u.ID == uuid.FromStringOrNil("") {
		return errors.New("user ID is required")
	}

	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	hp, err := s.hashPassword(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hp)

	return s.db.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte("users"))

		// Check if user already exists
		if v := bkt.Get([]byte(u.ID.String())); v != nil {
			return errors.New("user already exist")
		}

		if buf, err := json.Marshal(u); err != nil {
			return err
		} else if err := bkt.Put([]byte(u.ID.String()), buf); err != nil {
			return err
		}
		return nil
	})

}

// SetName updates user's name
func (s *UserService) SetName(userID uuid.UUID, name string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte("users"))

		var user user.User

		if v := bkt.Get([]byte(userID.String())); v == nil {
			return errors.New("user does not exist")
		} else if err := json.Unmarshal(v, &user); err != nil {
			return err
		}

		user.Name = name
		user.UpdatedAt = time.Now()

		if buf, err := json.Marshal(&user); err != nil {
			return err
		} else if err := bkt.Put([]byte(userID.String()), buf); err != nil {
			return err
		}
		return nil
	})
}

// DeleteUser removes user from bolt db by matching id key
func (s *UserService) DeleteUser(id uuid.UUID) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte("users"))
		if err := bkt.Delete([]byte(id.String())); err != nil {
			return err
		}
		return nil
	})
}

// CheckIfUserExists queries the db for the user by username
func (s *UserService) CheckIfUserExists(username string) (uuid.UUID, error) {
	var userID uuid.UUID
	var person *user.User

	err := s.db.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte("users"))
		c := bkt.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if err := json.Unmarshal(v, &person); err != nil {
				return err
			} else if person.Username == username {
				userID = person.ID
			}
		}
		if userID == uuid.Nil {
			return errors.New("user not found")
		}
		return nil
	})

	return userID, err
}

func (s *UserService) hashPassword(password string) ([]byte, error) {
	p := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(p, bcrypt.DefaultCost)
	return hash, err
}
