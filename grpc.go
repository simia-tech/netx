package netx

import (
	"net"
	"time"

	"github.com/simia-tech/netx/value"
)

// NewGRPCDialer returns a dialer that can be passed to the grpc.Dial function.
func NewGRPCDialer(network string, options ...value.DialOption) func(string, time.Duration) (net.Conn, error) {
	return func(address string, timeout time.Duration) (net.Conn, error) {
		return Dial(network, address, options...)
	}
}
