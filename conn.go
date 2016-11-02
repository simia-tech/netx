package netx

import (
	"net"

	"github.com/simia-tech/netx/internal/consul"
	"github.com/simia-tech/netx/internal/dnssrv"
	"github.com/simia-tech/netx/internal/nats"
)

// Dial establishs a connection on the provided network to the provided address.
func Dial(network, address string, options ...Option) (net.Conn, error) {
	o := &Options{}
	for _, option := range options {
		if err := option(o); err != nil {
			return nil, err
		}
	}

	switch network {
	case "nats":
		return nats.Dial(address, o.nodes)
	case "consul":
		return consul.Dial(address, o.nodes)
	case "dnssrv":
		return dnssrv.Dial(address, o.nodes)
	default:
		return net.Dial(network, address)
	}
}
