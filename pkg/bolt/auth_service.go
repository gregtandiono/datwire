package bolt

import (
	"datwire/pkg/apps/auth"
	"time"

	"golang.org/x/crypto/bcrypt"

	jwt "github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
)

// AuthService represents all service related to authentication
type AuthService struct{}

// Authorize returns an auth obj (id and token) after a successful authentication
func (s *AuthService) Authorize(email, password, passwordFromDB string, userID uuid.UUID) (auth *auth.Auth, err error) {
	if err = bcrypt.CompareHashAndPassword([]byte(password), []byte(passwordFromDB)); err != nil {
		return
	} else if token, tknerr := s.generateToken(userID, "d855496646e88b7c12e0a80135bef652"); tknerr != nil {
		err = tknerr
		return
	} else {
		auth.ID = userID.String()
		auth.Token = token
		return
	}
}

func (s *AuthService) generateToken(userID uuid.UUID, hashString string) (string, error) {
	mySigningKey := []byte(hashString)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"nbf": time.Now().Unix(),
		"id":  userID,
	})

	tokenString, err := token.SignedString(mySigningKey)
	return tokenString, err
}
