package netx

import (
	"net"
	"time"

	"github.com/simia-tech/netx/value"
)

// NewGRPCDialer returns a dialer that can be passed to the grpc.Dial function.
func NewGRPCDialer(network string, options ...value.Option) func(string, time.Duration) (net.Conn, error) {
	return func(address string, timeout time.Duration) (net.Conn, error) {
		return Dial(network, address, options...)
	}
}

// NewGRPCMultiDialer returns a dialer that can be passed to the grpc.Dial function.
func NewGRPCMultiDialer(md *MultiDialer) func(string, time.Duration) (net.Conn, error) {
	return func(address string, timeout time.Duration) (net.Conn, error) {
		return md.Dial(address)
	}
}
