package netx

import (
	"crypto/tls"
	"net"

	"github.com/simia-tech/netx/value"
)

var dialFuncs = map[string]DialFunc{}

// DialFunc defines the signature of the Dial function.
type DialFunc func(string, *value.DialOptions) (net.Conn, error)

// RegisterDial registers the provided Dial method under the provided network name.
func RegisterDial(network string, dialFunc DialFunc) {
	dialFuncs[network] = dialFunc
}

// RegisteredDialNetworks returns the available networks for the Dial function.
func RegisteredDialNetworks() []string {
	networks := []string{}
	for network := range dialFuncs {
		networks = append(networks, network)
	}
	return networks
}

// Dial establishs a connection on the provided network to the provided address.
func Dial(network, address string, options ...value.DialOption) (net.Conn, error) {
	o := &value.DialOptions{}
	for _, option := range options {
		if option == nil {
			continue
		}
		if err := option(o); err != nil {
			return nil, err
		}
	}

	dialFunc, ok := dialFuncs[network]
	if ok {
		return dialFunc(address, o)
	}
	if o.TLSConfig == nil {
		return net.Dial(network, address)
	}
	return tls.Dial(network, address, o.TLSConfig)
}
