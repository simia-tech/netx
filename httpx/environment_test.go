package httpx_test

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"testing"

	n "github.com/nats-io/nats"
	"github.com/stretchr/testify/require"

	"github.com/simia-tech/netx"
	"github.com/simia-tech/netx/httpx"
)

func setUpTestHTTPServer(tb testing.TB) (net.Addr, func()) {
	network := os.Getenv("NETWORK")
	if network == "" {
		tb.Skip("NETWORK is unset")
	}

	listener, err := netx.Listen(network, n.NewInbox())
	require.NoError(tb, err)

	mux := &http.ServeMux{}
	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "test")
	})

	server := &http.Server{Handler: mux}
	go func() {
		server.Serve(listener)
	}()

	return listener.Addr(), func() {
		require.NoError(tb, listener.Close())
	}
}

func setUpTestHTTPClient(tb testing.TB) *http.Client {
	network := os.Getenv("NETWORK")
	if network == "" {
		tb.Skip("NETWORK is unset")
	}

	transport := httpx.NewTransport(network)
	return &http.Client{Transport: transport}
}
