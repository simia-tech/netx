package netx

import (
	"net"
	"strings"

	"github.com/simia-tech/netx/internal/consul"
	"github.com/simia-tech/netx/internal/nats"
)

// Listen creates a listener on the provided network at the provided address.
func Listen(network, address string) (net.Listener, error) {
	switch {
	case strings.HasPrefix(network, "nats:"):
		return nats.Listen(network, address)
	case strings.HasPrefix(network, "consul:"):
		return consul.Listen(network, address)
	default:
		return net.Listen(network, address)
	}
}
