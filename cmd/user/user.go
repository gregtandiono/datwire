package main

import (
	dwhttp "datwire/pkg/http"
	"datwire/pkg/shared"
	"fmt"
	"log"
	"net/http"

	"github.com/rs/cors"
	"github.com/urfave/negroni"
)

func main() {
	n := negroni.New()
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
	})

	h := dwhttp.NewUserHandler()
	h.UserService.Open()
	defer h.UserService.Close()

	n.Use(negroni.HandlerFunc(shared.Logger))
	n.Use(c)
	n.UseHandler(h)

	serviceConfig := shared.GetEnvironmentVariables("datwire-users")

	fmt.Println("user service is running on port " + serviceConfig.Port)
	log.Fatal(http.ListenAndServe(":"+serviceConfig.Port, n))
}
