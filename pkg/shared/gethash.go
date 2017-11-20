package shared

import (
	"datwire/pkg/consul"
	"log"
	"os"
)

// GetHash interacts with Consul KV http API,
// to fetch the hash string to generate token in other services / handlers.
// If environment is set to `TEST`, it will default to a static hash string.
func GetHash() string {
	var hashString string
	if env := os.Getenv("ENV"); env != "TEST" {
		consuld := consul.NewConsuld(nil)
		hash, err := consuld.GetKV("datwire/config/hashString", nil)
		hashString = hash
		if err != nil {
			log.Fatal(err)
		}
	} else {
		hashString = "869826e158da8666906ec2681b19b96b729665fd2fae1328ace29171a1e8b3e2" // just for testing purposes
	}

	return hashString
}
