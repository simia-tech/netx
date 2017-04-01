package netx

import (
	"crypto/tls"
	"net"
)

var listenFuncs = map[string]ListenFunc{}

// ListenFunc defines the signature of the Listen function.
type ListenFunc func(string, *Options) (net.Listener, error)

// RegisterListen registers the provided Listen method under the provided network name.
func RegisterListen(network string, listenFunc ListenFunc) {
	listenFuncs[network] = listenFunc
}

// RegisteredListenNetworks returns the available networks for the Listen function.
func RegisteredListenNetworks() []string {
	networks := []string{}
	for network := range listenFuncs {
		networks = append(networks, network)
	}
	return networks
}

// Listen creates a listener on the provided network at the provided address.
func Listen(network, address string, options ...Option) (net.Listener, error) {
	o := &Options{}
	for _, option := range options {
		if option == nil {
			continue
		}
		if err := option(o); err != nil {
			return nil, err
		}
	}

	listenFunc, ok := listenFuncs[network]
	if ok {
		return listenFunc(address, o)
	}
	if o.TLSConfig == nil {
		return net.Listen(network, address)
	}
	return tls.Listen(network, address, o.TLSConfig)
}
