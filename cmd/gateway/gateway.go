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

	h := dwhttp.NewGatewayHandler()

	n.Use(negroni.HandlerFunc(shared.Logger))
	n.Use(c)
	n.UseHandler(h)

	fmt.Println("api gateway server is running on port 1338")
	log.Fatal(http.ListenAndServe(":1338", n))
}
