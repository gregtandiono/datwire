package user

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// User represents user data model
type User struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Type      string    `json:"type"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// UserService represents user model CRUD implementation to the BoltDB
type UserService interface {
	User(id uuid.UUID) (*User, error)
	Users() ([]User, error)
	CreateUser(u *User) error
	SetName(userID uuid.UUID, name string) error
	DeleteUser(id uuid.UUID) error
	CheckIfUserExists(username string) (uuid.UUID, error)
}
