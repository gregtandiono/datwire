package auth

import (
	jwt "github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
)

// Auth represents the return object once a user has been authorized
type Auth struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}

// AuthClaims struct to parse jwt tokens
type AuthClaims struct {
	ID string `json:"id"`
	jwt.StandardClaims
}

// AuthService represents REST client that interacts with the user service API
type AuthService interface {
	Authorize(password, hashedPassword string, userID uuid.UUID) (*Auth, error)
}
