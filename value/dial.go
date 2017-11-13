package value

import (
	"fmt"
	"net"
	neturl "net/url"
)

// Dial holds all parameters for a netx.Dial call.
type Dial interface {
	Network() string
	Address() string
	Options() []DialOption
}

type dial struct {
	network string
	address string
	options []DialOption
}

// NewDial returns a new Dial with the provided values.
func NewDial(network, address string, options ...DialOption) Dial {
	return &dial{
		network: network,
		address: address,
		options: options,
	}
}

// ParseDialURL parses the provided url and returns a Dial.
func ParseDialURL(url string, options ...DialOption) (Dial, error) {
	u, err := neturl.Parse(url)
	if err != nil {
		return nil, fmt.Errorf("parse url [%s]: %v", url, err)
	}
	return NewDial(u.Scheme, net.JoinHostPort(u.Hostname(), u.Port()), options...), nil
}

// MustParseDialURL works like ParseDialURL, but panics on error.
func MustParseDialURL(url string, options ...DialOption) Dial {
	d, err := ParseDialURL(url, options...)
	if err != nil {
		panic(err)
	}
	return d
}

func (d *dial) Network() string {
	return d.network
}

func (d *dial) Address() string {
	return d.address
}

func (d *dial) Options() []DialOption {
	return d.options
}

// DialURL builds an url from the provided Dial.
func DialURL(dial Dial) string {
	return fmt.Sprintf("%s://%s", dial.Network(), dial.Address())
}
