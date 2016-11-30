package netx

import (
	"log"
	"net"
	"sort"
)

var dialFuncs = map[string]DialFunc{}

// DialFunc defines the signature of the Dial function.
type DialFunc func(string, *Options) (net.Conn, error)

// RegisterDial registers the provided Dial method under the provided network name.
func RegisterDial(network string, dialFunc DialFunc) {
	dialFuncs[network] = dialFunc
}

// Dial establishs a connection on the provided network to the provided address.
func Dial(network, address string, options ...Option) (net.Conn, error) {
	o := &Options{}
	for _, option := range options {
		if err := option(o); err != nil {
			return nil, err
		}
	}

	dialFunc, ok := dialFuncs[network]
	if ok {
		return dialFunc(address, o)
	}
	return net.Dial(network, address)
}

// DialOne dials one of the provided addresses using the provided options.
func DialOne(addrs Addrs, options *Options) (net.Conn, error) {
	sort.Sort(addrs)

	balancer := options.Balancer
	if balancer == nil {
		balancer = DefaultOptions.Balancer
	}

	for {
		addr, err := balancer(addrs)
		if err != nil {
			return nil, err
		}
		if addr == nil {
			return nil, ErrServiceUnavailable
		}

		conn, err := net.Dial(addr.Network(), addr.String())
		if err != nil {
			log.Printf("error connecting to %s: %v", addr, err)
			continue
		}
		return conn, nil
	}
}
