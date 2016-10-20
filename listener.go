package netx

import (
	"net"
	"strings"

	"github.com/simia-tech/netx/nats"
)

func Listen(network, address string) (net.Listener, error) {
	switch {
	case strings.HasPrefix(network, "nats:"):
		return nats.Listen(network, address)
	default:
		return net.Listen(network, address)
	}
}
