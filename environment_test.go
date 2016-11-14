package netx_test

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/simia-tech/netx"
	"github.com/simia-tech/netx/test"
)

func setUpTestEchoListener(tb testing.TB, address string) (net.Listener, chan error) {
	listener, err := netx.Listen("tcp", address)
	require.NoError(tb, err)

	return listener, test.EchoServer(listener)
}

func setUpTestEchoClient(tb testing.TB, address string) net.Conn {
	conn, err := netx.Dial("tcp", address)
	require.NoError(tb, err)
	return conn
}
