package netx

import (
	"net"
	"strings"

	"github.com/simia-tech/netx/internal/consul"
	"github.com/simia-tech/netx/internal/dnssrv"
	"github.com/simia-tech/netx/internal/nats"
)

// Dial establishs a connection on the provided network to the provided address.
func Dial(network, address string) (net.Conn, error) {
	switch {
	case strings.HasPrefix(network, "nats:"):
		return nats.Dial(network, address)
	case strings.HasPrefix(network, "consul:"):
		return consul.Dial(network, address)
	case strings.HasPrefix(network, "dnssrv:"):
		return dnssrv.Dial(network, address)
	default:
		return net.Dial(network, address)
	}
}
