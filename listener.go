package netx

import "net"

var listenFuncs = map[string]ListenFunc{}

// ListenFunc defines the signature of the Listen function.
type ListenFunc func(string, *Options) (net.Listener, error)

// RegisterListen registers the provided Listen method under the provided network name.
func RegisterListen(network string, listenFunc ListenFunc) {
	listenFuncs[network] = listenFunc
}

// Listen creates a listener on the provided network at the provided address.
func Listen(network, address string, options ...Option) (net.Listener, error) {
	o := &Options{}
	for _, option := range options {
		if option == nil {
			continue
		}
		if err := option(o); err != nil {
			return nil, err
		}
	}

	listenFunc, ok := listenFuncs[network]
	if ok {
		return listenFunc(address, o)
	}
	return net.Listen(network, address)
}
