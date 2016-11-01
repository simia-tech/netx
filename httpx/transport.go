package httpx

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/simia-tech/netx"
)

// NewTransport returns a new transport instance for a http client. The instance will use
// netx.Dial to establish a connection.
func NewTransport(network string, options ...netx.Option) *http.Transport {
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: func(ctx context.Context, _, address string) (net.Conn, error) {
			return netx.Dial(network, address, options...)
		},
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
}
