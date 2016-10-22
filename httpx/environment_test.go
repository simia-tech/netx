package httpx_test

import (
	"fmt"
	"net"
	"net/http"
	"testing"

	n "github.com/nats-io/nats"
	"github.com/stretchr/testify/require"

	"github.com/simia-tech/netx"
	"github.com/simia-tech/netx/httpx"
	xxx "github.com/simia-tech/netx/internal/nats"
)

const defaultNatsURL = "nats://localhost:4222"

func setUpTestHTTPServer(tb testing.TB) (net.Addr, func()) {
	listener, err := netx.Listen(defaultNatsURL, n.NewInbox())
	require.NoError(tb, err)

	mux := &http.ServeMux{}
	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "test")
	})
	xxx.Dial("", "")

	server := &http.Server{Handler: mux}
	go func() {
		server.Serve(listener)
	}()

	return listener.Addr(), func() {
		require.NoError(tb, listener.Close())
	}
}

func setUpTestHTTPClient(tb testing.TB) *http.Client {
	transport := httpx.NewTransport(defaultNatsURL)
	return &http.Client{Transport: transport}
}
