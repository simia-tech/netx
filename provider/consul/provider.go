package consul

import (
	"github.com/simia-tech/netx/provider/consul/client"
	"github.com/simia-tech/netx/value"
)

// Consul implements the consul provider.
type Consul struct {
	Client              *client.Consul
	EndpointDialOptions []value.DialOption
}

// NewProvider returns a new consul provider.
func NewProvider(urls []string, options ...value.DialOption) (*Consul, error) {
	c, err := client.NewConsul(urls, options...)
	if err != nil {
		return nil, err
	}
	return &Consul{
		Client: c,
	}, nil
}

// Endpoints returns endpoints for the provided service.
func (p *Consul) Endpoints(service string) (value.Endpoints, error) {
	endpoints, err := p.Client.Service(service)
	if err != nil {
		return nil, err
	}
	for index, endpoint := range endpoints {
		endpoints[index] = value.NewEndpoint(endpoint.Network(), endpoint.Address(), p.EndpointDialOptions...)
	}
	return endpoints, nil
}
