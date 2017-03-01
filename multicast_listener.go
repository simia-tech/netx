package netx

import (
	"io"

	"github.com/pkg/errors"
)

var listenMulticastFuncs = map[string]ListenMulticastFunc{}

// ListenMulticastFunc defines the signature of the ListenMulticast function.
type ListenMulticastFunc func(string, *Options) (io.ReadCloser, error)

// RegisterListen registers the provided Listen method under the provided network name.
func RegisterListenMulticast(network string, listenMulticastFunc ListenMulticastFunc) {
	listenMulticastFuncs[network] = listenMulticastFunc
}

// ListenMulticast creates a multicast connection on the provided network at the provided address.
func ListenMulticast(network, address string, options ...Option) (io.ReadCloser, error) {
	o := &Options{}
	for _, option := range options {
		if option == nil {
			continue
		}
		if err := option(o); err != nil {
			return nil, err
		}
	}

	listenMulticastFunc, ok := listenMulticastFuncs[network]
	if ok {
		return listenMulticastFunc(address, o)
	}
	return nil, errors.Errorf("no ListenMulticast function registered for network [%s]", network)
}
