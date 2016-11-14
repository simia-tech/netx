package netx

import (
	"net"
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
