package consul

import (
	"math/rand"
	"net"
	"time"

	"github.com/pkg/errors"

	"github.com/simia-tech/netx"
)

func init() {
	netx.RegisterDial("consul", Dial)
}

// Dial establishes a connection to the provided address over the consul network.
func Dial(address string, options *netx.Options) (net.Conn, error) {
	consul, err := newConsulFrom(options.Nodes)
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

	rand.Seed(time.Now().UnixNano())

	switch l := len(addrs); l {
	case 0:
		return nil, errors.Errorf("could find any instances for service [%s]", address)
	case 1:
		return net.Dial("tcp", addrs[0].String())
	default:
		return net.Dial("tcp", addrs[rand.Intn(l)].String())
	}
}
