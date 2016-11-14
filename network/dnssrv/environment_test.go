package dnssrv_test

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/simia-tech/netx"
	_ "github.com/simia-tech/netx/network/consul"
	_ "github.com/simia-tech/netx/network/dnssrv"
	"github.com/simia-tech/netx/test"
)

func setUpTestEchoListener(tb testing.TB, addresses ...string) (net.Listener, chan error) {
	address := netx.RandomAddress("echo-")
	if len(addresses) > 0 {
		address = addresses[0]
	}

	listener, err := netx.Listen("consul", address, netx.Nodes("http://localhost:8500"), netx.LocalAddress("localhost:0"))
	require.NoError(tb, err)

	return listener, test.EchoServer(listener)
}

func setUpTestEchoClient(tb testing.TB, address string) net.Conn {
	conn, err := netx.Dial("dnssrv", address, netx.Nodes("localhost:8600"))
	require.NoError(tb, err)
	return conn
}
