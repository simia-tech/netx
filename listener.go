package netx

import (
	"net"

	"github.com/simia-tech/netx/internal/consul"
	"github.com/simia-tech/netx/internal/nats"
)

// Listen creates a listener on the provided network at the provided address.
func Listen(network, address string, options ...Option) (net.Listener, error) {
	o := &Options{}
	for _, option := range options {
		if err := option(o); err != nil {
			return nil, err
		}
	}

	switch network {
	case "nats":
		return nats.Listen(address, o.nodes, o.tlsConfig)
	case "consul":
		return consul.Listen(address, o.nodes, o.localAddress)
	default:
		return net.Listen(network, address)
	}
}
