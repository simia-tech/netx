package netx_test

import (
	"net"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/simia-tech/netx"
)

func setUpTestEchoListener(tb testing.TB, addresses ...string) (net.Listener, chan error) {
	network := os.Getenv("LISTEN_NETWORK")
	if network == "" {
		tb.Skip("LISTEN_NETWORK is unset")
	}

	address := netx.RandomAddress("echo-")
	if len(addresses) > 0 {
		address = addresses[0]
	}

	listener, err := netx.Listen(network, address)
	require.NoError(tb, err)

	errChan := make(chan error, 1)
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				errChan <- err
				return
			}

			data := requireRead(tb, conn)
			requireWrite(tb, conn, data)

			if err := conn.Close(); err != nil {
				errChan <- err
				return
			}
		}
	}()

	return listener, errChan
}

func setUpTestEchoClient(tb testing.TB, address string) net.Conn {
	network := os.Getenv("DIAL_NETWORK")
	if network == "" {
		tb.Skip("DIAL_NETWORK is unset")
	}

	conn, err := netx.Dial(network, address)
	require.NoError(tb, err)

	return conn
}
