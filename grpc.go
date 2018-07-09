package netx

import (
	"context"
	"net"
	"time"

	"github.com/simia-tech/netx/value"
)

// NewGRPCDialer returns a dialer that can be passed to the grpc.Dial function.
func NewGRPCDialer(network string, options ...value.Option) func(string, time.Duration) (net.Conn, error) {
	return func(address string, timeout time.Duration) (net.Conn, error) {
		if timeout > 0 {
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			conn, err := Dial(ctx, network, address, options...)
			cancel()
			return conn, err
		}
		return Dial(context.Background(), network, address, options...)
	}
}

// NewGRPCMultiDialer returns a dialer that can be passed to the grpc.Dial function.
func NewGRPCMultiDialer(md *MultiDialer) func(string, time.Duration) (net.Conn, error) {
	return func(address string, timeout time.Duration) (net.Conn, error) {
		if timeout > 0 {
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			conn, err := md.Dial(ctx, address)
			cancel()
			return conn, err
		}
		return md.Dial(context.Background(), address)
	}
}
