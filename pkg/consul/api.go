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

// GetKV is a wrapper around consul's client.KV.Get method.
// I'm only registering the Get method because for this application,
// I don't see any reason why it should write/post anything to the consul KV.
func (c *Consuld) GetKV(key string, q *api.QueryOptions) (value string, err error) {
	kv := c.client.KV()
	kvp, _, err := kv.Get(key, q)
	value = string(kvp.Value[:])
	return
}
