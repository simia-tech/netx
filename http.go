package netx

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/simia-tech/netx/value"
)

// NewHTTPTransport returns a new transport instance for a http client. The instance will use
// netx.Dial to establish a connection.
func NewHTTPTransport(network string, options ...value.DialOption) *http.Transport {
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

func NewHTTPMultiTransport(md *MultiDialer) *http.Transport {
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: func(ctx context.Context, _, address string) (net.Conn, error) {
			return md.Dial(address)
		},
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
}
