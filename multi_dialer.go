package netx

import (
	"fmt"
	"net"

	"github.com/simia-tech/netx/provider"
	"github.com/simia-tech/netx/selector"
)

type MultiDialer struct {
	provider provider.Interface
	selector selector.Interface
}

func NewMultiDialer(p provider.Interface, s selector.Interface) (*MultiDialer, error) {
	return &MultiDialer{
		provider: p,
		selector: s,
	}, nil
}

func (md *MultiDialer) Dial(service string) (net.Conn, error) {
	dials, err := md.provider.Dials(service)
	if err != nil {
		return nil, fmt.Errorf("provider: %v", err)
	}

	dial, err := md.selector.Select(dials)
	if err != nil {
		return nil, fmt.Errorf("selector: %v", err)
	}

	return Dial(dial.Network(), dial.Address(), dial.Options()...)
}
