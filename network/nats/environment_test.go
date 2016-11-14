package nats_test

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/simia-tech/netx"
	_ "github.com/simia-tech/netx/network/nats"
	"github.com/simia-tech/netx/test"
)

func setUpTestEchoListener(tb testing.TB, addresses ...string) (net.Listener, chan error) {
	address := netx.RandomAddress("echo-")
	if len(addresses) > 0 {
		address = addresses[0]
	}

	listener, err := netx.Listen("nats", address, netx.Nodes("nats://localhost:4222"))
	require.NoError(tb, err)

	return listener, test.EchoServer(listener)
}

func setUpTestEchoClient(tb testing.TB, address string) net.Conn {
	conn, err := netx.Dial("nats", address, netx.Nodes("nats://localhost:4222"))
	require.NoError(tb, err)
	return conn
}
