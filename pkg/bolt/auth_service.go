package bolt

import (
	"datwire/pkg/apps/auth"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"

	jwt "github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
)

// AuthService represents all service related to authentication
type AuthService struct {
	Hash string
}

var _ auth.AuthService = &AuthService{}

// Authorize returns an auth obj (id and token) after a successful authentication
func (s *AuthService) Authorize(password, hashedPassword string, userID uuid.UUID) (a *auth.Auth, err error) {
	if err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		err = errors.New("password does not match")
		return
	} else if token, tknerr := s.generateToken(userID, s.Hash); tknerr != nil {
		err = tknerr
		return
	} else {
		a = &auth.Auth{ID: userID.String(), Token: token}
		return
	}
}

func (s *AuthService) generateToken(userID uuid.UUID, hashString string) (string, error) {
	mySigningKey := []byte(hashString)
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"nbf": time.Now().Unix(),
			"id":  userID,
		},
		// @TODO:
		// jwt.StandardClaims{
		// 	ExpiresAt:
		// },
	)

	tokenString, err := token.SignedString(mySigningKey)
	return tokenString, err
}
