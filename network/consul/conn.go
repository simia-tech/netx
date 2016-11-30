package consul

import (
	"net"

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

	return netx.DialOne(addrs, options)
}
