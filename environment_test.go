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

func setUpTestEchoListener(tb testing.TB, addresses ...string) (net.Listener, chan error) {
	address := n.NewInbox()
	if len(addresses) > 0 {
		address = addresses[0]
	}

	listener, err := netx.Listen(defaultNatsURL, address)
	require.NoError(tb, err)

	errChan := make(chan error, 3)
	go func() {
		conn, err := listener.Accept()
		require.NoError(tb, err)

		buffer := [4]byte{}
		n, err := conn.Read(buffer[:])
		if err != nil {
			errChan <- err
			return
		}
		require.Equal(tb, len(buffer), n)

		n, err = conn.Write(buffer[:])
		if err != nil {
			errChan <- err
			return
		}
		require.Equal(tb, len(buffer), n)

		errChan <- conn.Close()
	}()

	return listener, errChan
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
