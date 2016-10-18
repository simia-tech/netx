package netx_test

import (
	"net"
	"testing"

	n "github.com/nats-io/nats"
	"github.com/stretchr/testify/require"

	"github.com/simia-tech/netx"
)

const defaultNatsURL = "nats://localhost:4222"

func setUpTestListener(tb testing.TB) net.Listener {
	address := n.NewInbox()

	listener, err := netx.Listen(defaultNatsURL, address)
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
