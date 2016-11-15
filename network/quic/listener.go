package quic

import (
	"net"

	"github.com/simia-tech/netx"
)

type listener struct{}

func init() {
	netx.RegisterListen("quic", Listen)
}

func Listen(address string, options *netx.Options) (net.Listener, error) {
	return net.Listen("tcp", address)
}
