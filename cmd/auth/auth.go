package main

import (
	dwhttp "datwire/pkg/http"
	"datwire/pkg/shared"
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

	h := dwhttp.NewAuthHandler()

	n.Use(negroni.HandlerFunc(shared.Logger))
	n.Use(c)
	n.UseHandler(h)

	log.Fatal(http.ListenAndServe(":1338", n))
}
