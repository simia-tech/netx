package netx

import (
	"net"
	"time"
)

// NewGRPCDialer returns a dialer that can be passed to the grpc.Dial function.
func NewGRPCDialer(network string, options ...Option) func(string, time.Duration) (net.Conn, error) {
	return func(address string, timeout time.Duration) (net.Conn, error) {
		return Dial(network, address, options...)
	}
}
