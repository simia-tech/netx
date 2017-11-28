package netx

import (
	"fmt"
	"io"

	"github.com/simia-tech/netx/value"
)

var dialMulticastFuncs = map[string]DialMulticastFunc{}

// DialMulticastFunc defines the signature of the Dial function.
type DialMulticastFunc func(string, *value.Options) (io.WriteCloser, error)

// RegisterDialMulticast registers the provided DialMulticast method under the provided network name.
func RegisterDialMulticast(network string, dialMulticastFunc DialMulticastFunc) {
	dialMulticastFuncs[network] = dialMulticastFunc
}

// RegisteredDialMulticastNetworks returns the available networks for the DialMulticast function.
func RegisteredDialMulticastNetworks() []string {
	networks := []string{}
	for network := range dialMulticastFuncs {
		networks = append(networks, network)
	}
	return networks
}

// DialMulticast opens a multicast connection on the provided network to the provided address.
func DialMulticast(network, address string, options ...value.Option) (io.WriteCloser, error) {
	o := &value.Options{}
	for _, option := range options {
		if option == nil {
			continue
		}
		if err := option(o); err != nil {
			return nil, err
		}
	}

	dialMulticastFunc, ok := dialMulticastFuncs[network]
	if ok {
		return dialMulticastFunc(address, o)
	}
	return nil, fmt.Errorf("no DialMulticast function registered for network [%s]", network)
}
