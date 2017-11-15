package consul

import (
	"log"

	"github.com/hashicorp/consul/api"
)

// Consuld is wrapper around the consul interface catered specifically for this application
type Consuld struct {
	client *api.Client
}

// NewConsuld creates a new instance of Consuld interface, and initiates a new consul client instance
func NewConsuld(consulConfig *api.Config) *Consuld {
	if consulConfig == nil {
		consulConfig = api.DefaultConfig()
	}
	client, err := api.NewClient(consulConfig)
	if err != nil {
		log.Fatal(err)
	}
	return &Consuld{
		client: client,
	}
}

// GetAllServices returns all services available in a consul cluster
func (c *Consuld) GetAllServices() (services map[string]*api.AgentService, err error) {
	services, err = c.client.Agent().Services()
	return
}
