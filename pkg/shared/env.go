package shared

import (
	"datwire/pkg/consul"
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
// Acceptable env variables: `DEVELOPMENT`, `CLUSTER_MODE`. `PRODUCTION`
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
		if service := services[serviceName]; service != nil {
			return &ServiceConfig{
				Address: service.Address,
				Port:    strconv.Itoa(service.Port),
			}
		}
	}
	return &ServiceConfig{}
}
