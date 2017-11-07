package value

import (
	"fmt"
	"net"
	neturl "net/url"
)

// Addr provides a static implementation of the net.Addr interface.
type Addr struct {
	network string
	address string
}

// NewAddr returns a new Addr.
func NewAddr(network, address string) *Addr {
	return &Addr{network: network, address: address}
}

// ParseAddrURL parse the provided url and returns an Addr.
func ParseAddrURL(url string) (*Addr, error) {
	u, err := neturl.Parse(url)
	if err != nil {
		return nil, fmt.Errorf("parse url [%s]: %v", url, err)
	}
	return NewAddr(u.Scheme, net.JoinHostPort(u.Hostname(), u.Port())), nil
}

// Network returns the network.
func (a *Addr) Network() string {
	return a.network
}

// String returns the address.
func (a *Addr) String() string {
	return a.address
}
