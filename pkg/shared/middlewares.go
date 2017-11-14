package shared

// func tokenParser(authHeader string) string {
// 	split := strings.Split(authHeader, " ")
// 	token := split[1]
// 	return token
// }

// func myJWTMiddleware() *jwtmiddleware.JWTMiddleware {
// 	env := anakkandang.GetEnvironmentVariables()
// 	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
// 		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
// 			return []byte(env.Hash), nil
// 		},
// 		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err string) {
// 			w.WriteHeader(401)
// 			json.NewEncoder(w).Encode(&ResponseTemplate{Message: "fail", Error: "required authorization token is invalid"})
// 		},
// 		SigningMethod: jwt.SigningMethodHS256,
// 	})

// 	return jwtMiddleware
// }
