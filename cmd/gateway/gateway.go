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

	// @TODO:
	// this is only for testing the consul API, will need to abstract this
	// to its own package.
	// client, err := api.NewClient(api.DefaultConfig())
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// agent := client.Agent()
	// services, err := agent.Services()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(services["datwire-gatewa"])

	fmt.Println("api gateway server is running on port 1338")
	log.Fatal(http.ListenAndServe(":1338", n))
}
