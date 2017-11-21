package netx

import (
	"fmt"
	"net"

	"github.com/simia-tech/netx/provider"
	"github.com/simia-tech/netx/selector"
)

// MultiDialer implements a multi dialer.
type MultiDialer struct {
	provider provider.Interface
	selector selector.Interface
}

// NewMultiDialer returns a new nulti dialer.
func NewMultiDialer(p provider.Interface, s selector.Interface) (*MultiDialer, error) {
	return &MultiDialer{
		provider: p,
		selector: s,
	}, nil
}

// Dial dials one the endpoints from the provided service.
func (md *MultiDialer) Dial(service string) (net.Conn, error) {
retry:
	endpoints, err := md.provider.Endpoints(service)
	if err != nil {
		return nil, fmt.Errorf("provider: %v", err)
	}

	endpoint, err := md.selector.Select(endpoints)
	if err != nil {
		return nil, fmt.Errorf("selector: %v", err)
	}

	conn, err := Dial(endpoint.Network(), endpoint.Address(), endpoint.Options()...)
	if err != nil {
		if _, ok := err.(*net.OpError); ok {
			goto retry
		}
		return nil, err
	}

	return conn, nil
}
