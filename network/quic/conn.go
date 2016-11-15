package quic

import (
	"net"

	"github.com/simia-tech/netx"
)

type conn struct{}

func init() {
	netx.RegisterDial("quic", Dial)
}

func Dial(address string, options *netx.Options) (net.Conn, error) {
	return net.Dial("tcp", address)
}
