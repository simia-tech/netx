package consul

import (
	"net"

	"github.com/simia-tech/netx/provider/consul/client"
	"github.com/simia-tech/netx/value"
)

// Consul implements the consul provider.
type Consul struct {
	Client              *client.Consul
	EndpointDialOptions []value.Option
}

// NewProvider returns a new consul provider.
func NewProvider(urls []string, options ...value.Option) (*Consul, error) {
	c, err := client.NewConsul(urls, options...)
	if err != nil {
		return nil, err
	}
	return &Consul{
		Client: c,
	}, nil
}

// AddEndpointAddr adds the provided address to the consul register.
func (p *Consul) AddEndpointAddr(service string, addr net.Addr) (string, error) {
	return p.Client.Register(service, value.NewEndpoint(addr.Network(), addr.String()))
}

// AddEndpoint adds the provided endpoint to the consul register.
func (p *Consul) AddEndpoint(service string, endpoint value.Endpoint) (string, error) {
	return p.Client.Register(service, endpoint)
}

// RemoveEndpoint removes the endpoint with the provided id from the consul register.
func (p *Consul) RemoveEndpoint(id string) error {
	return p.Client.Deregister(id)
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
