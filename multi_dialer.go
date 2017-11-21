package netx

import (
	"fmt"
	"net"

	"github.com/simia-tech/netx/filter"
	"github.com/simia-tech/netx/provider"
	"github.com/simia-tech/netx/selector"
)

// MultiDialer implements a multi dialer.
type MultiDialer struct {
	provider provider.Interface
	filter   filter.Interface
	selector selector.Interface
}

// NewMultiDialer returns a new nulti dialer.
func NewMultiDialer(p provider.Interface, f filter.Interface, s selector.Interface) (*MultiDialer, error) {
	return &MultiDialer{
		provider: p,
		filter:   f,
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

	if md.filter != nil {
		endpoints, err = md.filter.Filter(endpoints)
		if err != nil {
			return nil, fmt.Errorf("filter: %v", err)
		}
	}

	endpoint, err := md.selector.Select(endpoints)
	if err != nil {
		return nil, fmt.Errorf("selector: %v", err)
	}

	conn, err := Dial(endpoint.Network(), endpoint.Address(), endpoint.Options()...)
	if err != nil {
		if _, ok := err.(*net.OpError); ok {
			if md.filter != nil {
				if err = md.filter.Failure(endpoint); err != nil {
					return nil, fmt.Errorf("report failure to filter: %v", err)
				}
			}
			goto retry
		}
		return nil, err
	}
	if md.filter != nil {
		if err = md.filter.Success(endpoint); err != nil {
			return nil, fmt.Errorf("report success to filter: %v", err)
		}
	}

	return conn, nil
}
