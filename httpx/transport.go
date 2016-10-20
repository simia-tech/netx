package httpx

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/simia-tech/netx"
)

func NewTransport(network string) *http.Transport {
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: func(ctx context.Context, _, address string) (net.Conn, error) {
			return netx.Dial(network, address)
		},
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
}
