package netx

import (
	"context"
	"crypto/tls"
	"net"

	"github.com/simia-tech/netx/value"
)

var dialFuncs = map[string]DialFunc{}

// DialFunc defines the signature of the Dial function.
type DialFunc func(context.Context, string, *value.Options) (net.Conn, error)

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
func Dial(ctx context.Context, network, address string, options ...value.Option) (net.Conn, error) {
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
		return dialFunc(ctx, address, o)
	}

	var d net.Dialer
	conn, err := d.DialContext(ctx, network, address)
	if err != nil {
		return nil, err
	}

	if o.TLSConfig != nil {
		conn = tls.Client(conn, o.TLSConfig)
	}
	return conn, nil
}
