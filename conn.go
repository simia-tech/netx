package netx

import (
	"crypto/tls"
	"log"
	"net"

	"github.com/simia-tech/netx/value"
)

var dialFuncs = map[string]DialFunc{}

// DialFunc defines the signature of the Dial function.
type DialFunc func(string, *value.Options) (net.Conn, error)

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
func Dial(network, address string, options ...value.Option) (net.Conn, error) {
	o := &value.Options{}
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

	var (
		conn net.Conn
		err  error
	)
	if o.Timeout == 0 {
		conn, err = net.Dial(network, address)
	} else {
		log.Printf("net %s / addr %s / to %s", network, address, o.Timeout)
		conn, err = net.DialTimeout(network, address, o.Timeout)
	}
	if err != nil {
		return nil, err
	}

	if o.TLSConfig != nil {
		conn = tls.Client(conn, o.TLSConfig)
	}
	return conn, nil
}
