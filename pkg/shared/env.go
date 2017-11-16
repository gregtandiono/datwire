package shared

import (
	"datwire/pkg/consul"
	"fmt"
	"log"
	"os"
	"strconv"
)

// ServiceConfig represents service configuration variables
type ServiceConfig struct {
	Address string
	Port    string
	Hash    string // hash is very specific to a service that actually needs it (auth, user, jtw middleware)
	Env     string
}

// GetEnvironmentVariables returns a serviceConfig based on the server environment.
// If the environment is set to DEV, then it will look for a config.toml file.
// On the other hand, it can also query config from consul.
// Acceptable env variables: `DEVELOPMENT`, `CLUSTER_MODE`. `PRODUCTION`.
// @TODO:
// each service may have a unique requirement, so we can't simply generalize
// their configuration needs, so, either abstract this method per service name basis
// or enhance it to satisfy all services' needs.
func GetEnvironmentVariables(serviceName string) *ServiceConfig {
	env := os.Getenv("ENV")
	if env == "" {
		env = "DEVELOPMENT"
	}
	switch env {
	// @NOTE
	// I'm gonna set to cluster mode by default for now
	case "DEVELOPMENT":
		fallthrough
	case "PRODUCTION":
		fallthrough
	case "CLUSTER_MODE":
		consuld := consul.NewConsuld(nil)
		services, err := consuld.GetAllServices()
		if err != nil {
			log.Fatal(err)
		}
		hash, err := consuld.GetKV("datwire/config/hashString", nil)
		if err != nil {
			fmt.Println(err)
		}
		if service := services[serviceName]; service != nil {
			return &ServiceConfig{
				Address: service.Address,
				Port:    strconv.Itoa(service.Port),
				Hash:    hash,
			}
		}
	}
	return &ServiceConfig{}
}
