package consul

import (
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

	if len(addrs) == 0 {
		return nil, errors.Errorf("could find any instances for service [%s]", address)
	}

	return net.Dial("tcp", addrs[0].String())
}
