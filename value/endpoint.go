package value

import (
	"fmt"
	"net"
	neturl "net/url"
	"strconv"
)

// Endpoint holds all parameters for a netx.Endpoint call.
type Endpoint interface {
	Network() string
	Address() string
	Options() []Option
}

type endpoint struct {
	network string
	address string
	options []Option
}

// NewEndpoint returns a new Endpoint with the provided values.
func NewEndpoint(network, address string, options ...Option) Endpoint {
	return &endpoint{
		network: network,
		address: address,
		options: options,
	}
}

// NewEndpointFromAddr returns a new Endpoint with values from the provided net.Addr.
func NewEndpointFromAddr(addr net.Addr, options ...Option) Endpoint {
	return NewEndpoint(addr.Network(), addr.String(), options...)
}

// ParseEndpointURL parses the provided url and returns a Endpoint.
func ParseEndpointURL(url string, options ...Option) (Endpoint, error) {
	u, err := neturl.Parse(url)
	if err != nil {
		return nil, fmt.Errorf("parse url [%s]: %v", url, err)
	}
	return NewEndpoint(u.Scheme, net.JoinHostPort(u.Hostname(), u.Port()), options...), nil
}

// MustParseEndpointURL works like ParseEndpointURL, but panics on error.
func MustParseEndpointURL(url string, options ...Option) Endpoint {
	d, err := ParseEndpointURL(url, options...)
	if err != nil {
		panic(err)
	}
	return d
}

func (d *endpoint) Network() string {
	return d.network
}

func (d *endpoint) Address() string {
	return d.address
}

func (d *endpoint) Options() []Option {
	return d.options
}

// EndpointURL builds an url from the provided Endpoint.
func EndpointURL(ep Endpoint) string {
	return fmt.Sprintf("%s://%s", ep.Network(), ep.Address())
}

// EndpointPort returns the port of the provided endpoint.
func EndpointPort(ep Endpoint) int {
	_, p, err := net.SplitHostPort(ep.Address())
	if err != nil {
		panic(fmt.Sprintf("could not split address [%s]", ep.Address()))
	}
	port, err := strconv.Atoi(p)
	if err != nil {
		panic(fmt.Sprintf("could not convert [%s] to int", p))
	}
	return port
}
