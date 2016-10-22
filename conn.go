package netx

import (
	"net"
	"strings"

	"github.com/simia-tech/netx/internal/nats"
)

// Dial establishs a connection on the provided network to the provided address.
func Dial(network, address string) (net.Conn, error) {
	switch {
	case strings.HasPrefix(network, "nats:"):
		return nats.Dial(network, address)
	default:
		return net.Dial(network, address)
	}
}
