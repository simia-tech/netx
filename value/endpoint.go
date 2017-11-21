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
	Options() []DialOption
}

type endpoint struct {
	network string
	address string
	options []DialOption
}

// NewEndpoint returns a new Dial with the provided values.
func NewEndpoint(network, address string, options ...DialOption) Endpoint {
	return &endpoint{
		network: network,
		address: address,
		options: options,
	}
}

// ParseEndpointURL parses the provided url and returns a Dial.
func ParseEndpointURL(url string, options ...DialOption) (Endpoint, error) {
	u, err := neturl.Parse(url)
	if err != nil {
		return nil, fmt.Errorf("parse url [%s]: %v", url, err)
	}
	return NewEndpoint(u.Scheme, net.JoinHostPort(u.Hostname(), u.Port()), options...), nil
}

// MustParseEndpointURL works like ParseDialURL, but panics on error.
func MustParseEndpointURL(url string, options ...DialOption) Endpoint {
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

func (d *endpoint) Options() []DialOption {
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
