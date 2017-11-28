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
func NewHTTPTransport(network string, options ...value.Option) *http.Transport {
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

// NewHTTPMultiTransport returns a new transport instance that uses the provided MultiDialer.
func NewHTTPMultiTransport(md *MultiDialer) *http.Transport {
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: func(ctx context.Context, _, address string) (net.Conn, error) {
			host, _, _ := net.SplitHostPort(address)
			return md.Dial(host)
		},
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
}
