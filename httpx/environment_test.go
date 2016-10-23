package httpx_test

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/simia-tech/netx"
	"github.com/simia-tech/netx/httpx"
)

func setUpTestHTTPServer(tb testing.TB, addresses ...string) (net.Addr, func() int, func()) {
	network := os.Getenv("LISTEN_NETWORK")
	if network == "" {
		tb.Skip("LISTEN_NETWORK is unset")
	}

	address := netx.RandomAddress("http-")
	if len(addresses) > 0 {
		address = addresses[0]
	}

	listener, err := netx.Listen(network, address)
	require.NoError(tb, err)

	counter := new(int)
	mux := &http.ServeMux{}
	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		*counter++
		fmt.Fprintf(w, "test")
	})

	server := &http.Server{Handler: mux}
	go func() {
		server.Serve(listener)
	}()

	return listener.Addr(), func() int {
			return *counter
		}, func() {
			require.NoError(tb, listener.Close())
		}
}

func setUpTestHTTPClient(tb testing.TB) *http.Client {
	network := os.Getenv("DIAL_NETWORK")
	if network == "" {
		tb.Skip("DIAL_NETWORK is unset")
	}

	transport := httpx.NewTransport(network)
	return &http.Client{Transport: transport}
}
