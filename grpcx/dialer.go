package grpcx

import (
	"net"
	"time"

	"github.com/simia-tech/netx"
)

func NewDialer(network string, options ...netx.Option) func(string, time.Duration) (net.Conn, error) {
	return func(address string, timeout time.Duration) (net.Conn, error) {
		return netx.Dial(network, address, options...)
	}
}
