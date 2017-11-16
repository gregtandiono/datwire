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
	Env     string
}

// GetEnvironmentVariables returns a ServiceConfig based on the server environment.
// @TODO:
// each service may have a unique requirement, so we can't simply generalize
// their configuration needs, so, either abstract this method per service name basis
// or enhance it to satisfy all services' needs.
func GetEnvironmentVariables(serviceName string) *ServiceConfig {
	env := os.Getenv("ENV")

	if env == "" {
		env = "DEVELOPMENT"
	}
	consuld := consul.NewConsuld(nil)
	services, err := consuld.GetAllServices()
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}
	if service := services[serviceName]; service != nil {
		return &ServiceConfig{
			Address: service.Address,
			Port:    strconv.Itoa(service.Port),
			Env:     env,
		}
	}
	return &ServiceConfig{}
}
