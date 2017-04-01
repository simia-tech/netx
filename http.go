package netx

import (
	"context"
	"net"
	"net/http"
	"time"
)

// NewHTTPTransport returns a new transport instance for a http client. The instance will use
// netx.Dial to establish a connection.
func NewHTTPTransport(network string, options ...Option) *http.Transport {
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: func(ctx context.Context, _, address string) (net.Conn, error) {
			return Dial(network, address, options...)
		},
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
}
