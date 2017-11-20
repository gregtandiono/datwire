package shared

import (
	"encoding/json"
	"net/http"
	"strings"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"
)

// TokenParser parses token from request header
func TokenParser(authHeader string) string {
	split := strings.Split(authHeader, " ")
	token := split[1]
	return token
}

// JWTMiddleware is a pluggable middleware to an existing route
// to check whether there is a valid token in a request header.
func JWTMiddleware(secretHash string) *jwtmiddleware.JWTMiddleware {
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(secretHash), nil
		},
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err string) {
			w.WriteHeader(401)
			json.NewEncoder(w).Encode(&ResponseTemplate{Message: "fail", Error: "required authorization token is invalid"})
		},
		SigningMethod: jwt.SigningMethodHS256,
	})

	return jwtMiddleware
}
