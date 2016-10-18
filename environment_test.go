package netx_test

import (
	"fmt"
	"net"
	"net/http"
	"testing"

	n "github.com/nats-io/nats"
	"github.com/stretchr/testify/require"

	"github.com/simia-tech/netx"
)

const defaultNatsURL = "nats://localhost:4222"

func setUpTestListener(tb testing.TB) net.Listener {
	listener, err := netx.Listen(defaultNatsURL, n.NewInbox())
	require.NoError(tb, err)

	go func() {
		conn, err := listener.Accept()
		require.NoError(tb, err)

		buffer := requireRead(tb, conn, 4)
		requireWrite(tb, conn, buffer)

		require.NoError(tb, conn.Close())
	}()

	return listener
}

func setUpTestHTTPServer(tb testing.TB) (net.Addr, func()) {
	listener, err := netx.Listen(defaultNatsURL, n.NewInbox())
	require.NoError(tb, err)

	mux := &http.ServeMux{}
	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "test")
	})

	server := &http.Server{
		Handler: mux,
	}
	go func() {
		require.NoError(tb, server.Serve(listener))
	}()

	return listener.Addr(), func() {
		require.NoError(tb, listener.Close())
	}
}

func setUpTestHTTPClient(tb testing.TB) *http.Client {
	network, err := netx.NewNetwork(defaultNatsURL)
	require.NoError(tb, err)
	transport := netx.NewTransport(network)
	return &http.Client{Transport: transport}
}
