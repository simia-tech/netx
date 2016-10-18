package netx

import (
	"net"
	"strings"

	"github.com/pkg/errors"
	"github.com/simia-tech/netx/nats"
)

type Network interface {
	Listen(string) (net.Listener, error)
	Dial(string) (net.Conn, error)
}

func NewNetwork(net string) (Network, error) {
	switch {
	case strings.HasPrefix(net, "nats:"):
		return nats.JoinNetwork(net)
	default:
		return nil, errors.Errorf("unknown network [%s]", net)
	}
}

func Listen(net, address string) (net.Listener, error) {
	network, err := NewNetwork(net)
	if err != nil {
		return nil, err
	}
	return network.Listen(address)
}

func Dial(net, address string) (net.Conn, error) {
	network, err := NewNetwork(net)
	if err != nil {
		return nil, err
	}
	return network.Dial(address)
}
