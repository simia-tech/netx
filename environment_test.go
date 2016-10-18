package netx_test

import (
	"net"
	"testing"

	"github.com/simia-tech/netx"
	"github.com/stretchr/testify/require"
)

const defaultNatsURL = "nats://localhost:4222"

func setUpTestListener(tb testing.TB) net.Listener {
	listener, err := netx.Listen(defaultNatsURL, "test")
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
