package consul

import (
	"math/rand"
	"net"

	"github.com/pkg/errors"
)

// Dial establishes a connection to the provided address over the provided network.
func Dial(network, address string) (net.Conn, error) {
	consul, _, err := newConsulFrom(network)
	if err != nil {
		return nil, err
	}

	host, _, err := net.SplitHostPort(address)
	if err != nil {
		host = address
	}

	addrs, err := consul.service(host)
	if err != nil {
		return nil, err
	}

	switch l := len(addrs); l {
	case 0:
		return nil, errors.Errorf("could find any instances for service [%s]", address)
	case 1:
		return net.Dial("tcp", addrs[0].String())
	default:
		return net.Dial("tcp", addrs[rand.Intn(l)].String())
	}
}
